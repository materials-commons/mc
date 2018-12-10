package globus

import (
	"context"
	"strings"
	"time"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/globusapi"

	"github.com/materials-commons/mc/internal/store"

	"github.com/apex/log"
)

type UploadMonitor struct {
	client              *globusapi.Client
	globusUploads       *store.GlobusUploadsStore
	fileLoads           *store.FileLoadsStore
	endpointID          string
	finishedGlobusTasks map[string]bool
}

func NewUploadMonitor(client *globusapi.Client, endpointID string, db store.DB) *UploadMonitor {
	return &UploadMonitor{
		client:              client,
		endpointID:          endpointID,
		globusUploads:       db.GlobusUploadsStore(),
		fileLoads:           db.FileLoadsStore(),
		finishedGlobusTasks: make(map[string]bool),
	}
}

func (m *UploadMonitor) Start(c context.Context) {
	go m.monitorAndProcessUploads(c)
}

func (m *UploadMonitor) monitorAndProcessUploads(c context.Context) {
	log.Infof("Starting globus monitoring...")
	for {
		m.retrieveAndProcessUploads(c)
		select {
		case <-c.Done():
			log.Infof("Shutting down globus monitoring...")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (m *UploadMonitor) retrieveAndProcessUploads(c context.Context) {
	// Build a filter to get all successful tasks that completed in the last week
	lastWeek := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	taskFilter := map[string]string{
		"filter_completion_time": lastWeek,
		"filter_status":          "SUCCEEDED",
	}
	tasks, err := m.client.GetEndpointTaskList(m.endpointID, taskFilter)

	if err != nil {
		log.Infof("globus.GetEndpointTaskList returned the following error: %s - %#v", err, m.client.GetGlobusErrorResponse())
		return
	}

	for _, task := range tasks.Tasks {
		//log.Infof("Getting successful transfers for Globus Task %s", task.TaskID)
		transfers, err := m.client.GetTaskSuccessfulTransfers(task.TaskID, 0)

		switch {
		case err != nil:
			log.Infof("globus.GetTaskSuccessfulTransfers(%d) returned error %s - %#v", task.TaskID, err, m.client.GetGlobusErrorResponse())
			continue
		case len(transfers.Transfers) == 0:
			// No files transferred in this request
			continue
		default:
			// Files were transferred for this request
			m.processTransfers(&transfers)
		}

		// Check if we should stop processing requests
		select {
		case <-c.Done():
			break
		default:
		}
	}
}

func (m *UploadMonitor) processTransfers(transfers *globusapi.TransferItems) {
	transferItem := transfers.Transfers[0]

	// Transfer items with a blank DestinationPath are downloads not uploads.
	if transferItem.DestinationPath == "" {
		return
	}

	// Destination path will have the following format: /__globus_uploads/<id of upload request>/...rest of path...
	// Split will return ["", "__globus_uploads", "<id of upload request", ....]
	// So the 3rd entry in the array is the id in the globus_uploads table we want to look up.
	pieces := strings.Split(transferItem.DestinationPath, "/")
	if len(pieces) < 4 {
		// sanity check, because the destination path should at least be /__globus_uploads/<id>/...rest of path...
		// it should at least have 4 entries in it (See Split return description above)
		log.Infof("Invalid globus DestinationPath: %s", transferItem.DestinationPath)
		return
	}

	id := pieces[2] // id is the 3rd entry in the path
	if _, ok := m.finishedGlobusTasks[id]; ok {
		// We've seen this globus task before and already processed it
		return
	}
	globusUpload, err := m.globusUploads.GetGlobusUpload(id)
	if err != nil {
		// If we find a Globus task, but no corresponding entry in our database that means at some
		// earlier point in time we processed the task by turning it into a file load request and
		// deleting globus upload from our database. So this is an old reference we can just ignore.
		// Add the entry to our hash table of completed requests.
		m.finishedGlobusTasks[id] = true
		return
	}

	// At this point we have a globus upload. What we are going to do is remove the ACL on the directory
	// so no more files can be uploaded to it. Then we are going to add that directory to the list of
	// directories to upload. Then the file loader will eventually get around to loading these files. In
	// the meantime since we've now created a file load from this globus upload we can delete the entry
	// from the globus_uploads table. Finally we are going to update the status for this background process.

	log.Infof("Processing globus upload %s", id)

	if _, err := m.client.DeleteEndpointACLRule(m.endpointID, globusUpload.GlobusAclID); err != nil {
		log.Infof("Unable to delete ACL: %s", err)
	}

	flAdd := model.AddFileLoadModel{
		ProjectID:      globusUpload.ProjectID,
		Owner:          globusUpload.Owner,
		Path:           globusUpload.Path,
		GlobusUploadID: globusUpload.ID,
	}

	if fl, err := m.fileLoads.AddFileLoad(flAdd); err != nil {
		log.Infof("Unable to add file load request: %s", err)
		return
	} else {
		log.Infof("Created file load (id: %s) for globus upload %s", fl.ID, id)
	}

	// Delete the globus upload request as we have now turned it into a file loading request
	// and won't have to process this request again. If the server stops while loading the
	// request or there is some other failure, the file loader will take care of picking up
	// where it left off.
	m.globusUploads.DeleteGlobusUpload(id)
}
