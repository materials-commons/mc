package file

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/Jeffail/tunny"
	"github.com/materials-commons/mc/internal/store"
)

type BackgroundLoader struct {
	mcdir           string
	numberOfWorkers int
	db              store.DB
}

func NewBackgroundLoader(mcdir string, numberOfWorkers int, db store.DB) *BackgroundLoader {
	return &BackgroundLoader{
		mcdir:           mcdir,
		numberOfWorkers: numberOfWorkers,
		db:              db,
	}
}

func (l *BackgroundLoader) Start() {
	go l.loadFilesInBackground()
}

func (l *BackgroundLoader) loadFilesInBackground() {
	pool := tunny.NewFunc(l.numberOfWorkers, l.worker)
	fileloadsStore := l.db.FileLoadsStore()

	// There may have been jobs in process when the server was stopped. Mark those jobs
	// at not currently being processed, this will cause them to be re-processed.
	if err := fileloadsStore.MarkAllNotLoading(); err != nil {
		panic(fmt.Sprintf("Unable to mark current jobs as not loading: %s", err))
	}

	// Loop through all file load requests and look for any that are not currently being processed.
	for {
		requests, err := fileloadsStore.GetAllFileLoads()
		fmt.Printf("processing file loads %#v: %s\n", requests, err)
		if err != nil && errors.Cause(err) != store.ErrNotFound {
			fmt.Println("Error retrieving requests:", err)
		}

		for _, req := range requests {
			if !req.Loading {

				// Mark job as loading so it will be ignored in later processing
				err := fileloadsStore.UpdateLoading(req.ID, true)
				if err != nil {
					fmt.Printf("Unable to update file load request %s: %s", req.ID, err)

					// If the job cannot be marked as loading then skip processing it
					continue
				}

				// pool.Process() is synchronous, so run in separate routine and let the pool control
				// how many jobs are running simultaneously.
				go func() {
					pool.Process(req)
				}()
			}
		}

		// Sleep for 10 seconds before getting the next set of loading requests. Ten seconds is an
		// somewhat arbitrary value chosen to balance time to start processing and load.
		time.Sleep(10 * time.Second)
	}
}

func (l *BackgroundLoader) worker(args interface{}) interface{} {
	dfStore := l.db.DatafilesStore()
	ddStore := l.db.DatadirsStore()
	projStore := l.db.ProjectsStore()
	flStore := l.db.FileLoadsStore()

	req := args.(store.FileLoadSchema)

	proj, err := projStore.GetProjectSimple(req.ProjectID)
	if err != nil {
		return err
	}

	loader := NewMCFileLoader(req.Path, req.Owner, l.mcdir, proj, dfStore, ddStore)
	skipper := NewExcludeListSkipper(req.Exclude)
	fl := NewFileLoader(skipper.Skipper, loader)

	if err := fl.LoadFiles(req.Path); err != nil {
		return err
	} else {
		// if loading files was successful then
		//    remove path since all files were processed and
		//    delete this file load request
		_ = os.RemoveAll(req.Path)
		return flStore.DeleteFileLoad(req.ID)
	}
}
