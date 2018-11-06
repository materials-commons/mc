package file

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"

	"github.com/pkg/errors"

	"github.com/Jeffail/tunny"
	"github.com/materials-commons/mc/internal/store"
)

type BackgroundLoader struct {
	mcdir           string
	numberOfWorkers int
	db              store.DB
	c               context.Context
	activeProjects  sync.Map
}

func NewBackgroundLoader(mcdir string, numberOfWorkers int, db store.DB) *BackgroundLoader {
	return &BackgroundLoader{
		mcdir:           mcdir,
		numberOfWorkers: numberOfWorkers,
		db:              db,
	}
}

func (l *BackgroundLoader) Start(c context.Context) {
	go l.processLoadFileRequests(c)
}

func (l *BackgroundLoader) processLoadFileRequests(c context.Context) {
	l.c = c
	pool := tunny.NewFunc(l.numberOfWorkers, l.worker)
	fileloadsStore := l.db.FileLoadsStore()

	// There may have been jobs in process when the server was stopped. Mark those jobs
	// at not currently being processed, this will cause them to be re-processed.
	if err := fileloadsStore.MarkAllNotLoading(); err != nil && errors.Cause(err) != mc.ErrNotFound {
		log.Infof("Unable to mark current jobs as not loading: %s\n", err)
	}

	// Loop through all file load requests and look for any that are not currently being processed.
	for {
		requests, err := fileloadsStore.GetAllFileLoads()
		//fmt.Printf("processing file loads %#v: %s\n", requests, err)
		if err != nil && errors.Cause(err) != mc.ErrNotFound {
			log.Infof("Error retrieving requests: %s", err)
		}

		for _, req := range requests {

			if req.Loading {
				// This request is already being processed so ignore it
				continue
			}

			// Check if project is already being processed
			if _, ok := l.activeProjects.Load(req.ProjectID); ok {
				// There is already a load active for this project, so skip
				// processing this load request
				continue
			}

			log.Infof("processing request %#v\n", req)

			// If we are here then the current request is not being processed
			// and it is for a project that is *not* currently being processed.

			// Mark job as loading so we won't attempt to load this request a second time
			if err := fileloadsStore.UpdateLoading(req.ID, true); err != nil {
				log.Infof("Unable to update file load request %s: %s", req.ID, err)

				// If the job cannot be marked as loading then skip processing it
				continue
			}

			// Lock the project so no other uploads for this project will be processed. This is
			// done to prevent issues such as checks and writes to the database creating two entries.
			l.activeProjects.Store(req.ProjectID, true)

			// pool.Process() is synchronous, so run in separate routine and let the pool control
			// how many jobs are running simultaneously.
			go func() {
				// pool.Process will call the worker function (below) for processing the request
				pool.Process(req)
			}()
		}

		// Sleep for 10 seconds before getting the next set of loading requests. Ten seconds is an
		// somewhat arbitrary value chosen to balance time to start processing and load on the system.
		select {
		case <-time.After(10 * time.Second):
		case <-c.Done():
			log.Infof("Shutting down file loading...")
			return
		}
	}
}

// worker is the worker function for the pool. A new worker will run for each request that
// is being run up to the pool limit.
func (l *BackgroundLoader) worker(args interface{}) interface{} {
	dfStore := l.db.DatafilesStore()
	ddStore := l.db.DatadirsStore()
	projStore := l.db.ProjectsStore()
	flStore := l.db.FileLoadsStore()

	req := args.(model.FileLoadSchema)

	proj, err := projStore.GetProjectSimple(req.ProjectID)
	if err != nil {
		return err
	}

	loader := NewMCFileLoader(req.Path, req.Owner, l.mcdir, proj, dfStore, ddStore)
	skipper := NewExcludeListSkipper(req.Exclude)
	fl := NewFileLoader(skipper.Skipper, loader)

	if err := fl.LoadFilesWithCancel(req.Path, l.c); err != nil {
		// Load wasn't successful. Release the project so other load
		// requests for the project can proceed. Also, mark this
		// load request as not being loaded so it can be retried.
		//
		// Ignore errors as we are already in an error situation, either the updates
		// to the database will work or they won't in which case nothing else will work.
		flStore.UpdateLoading(req.ID, false)
		l.activeProjects.Delete(req.ProjectID)
		return err
	} else {
		// if loading files was successful then
		//    1. remove path since all files were processed and
		//    2. remove project from activeProjects map so other requests can be processed
		//    3. delete this file load request
		_ = os.RemoveAll(req.Path)
		l.activeProjects.Delete(req.ProjectID)
		return flStore.DeleteFileLoad(req.ID)
	}
}
