package globus

import (
	"context"
	"strings"
	"time"

	"github.com/materials-commons/mc/internal/store/model"

	"github.com/materials-commons/mc/internal/store"

	"github.com/apex/log"
)

type UploadMonitor struct {
	client        *Client
	globusUploads *store.GlobusUploadsStore
	fileLoads     *store.FileLoadsStore
	endpointID    string
}

func NewUploadMonitor(client *Client, endpointID string, db store.DB) *UploadMonitor {
	return &UploadMonitor{
		client:        client,
		endpointID:    endpointID,
		globusUploads: db.GlobusUploadsStore(),
		fileLoads:     db.FileLoadsStore(),
	}
}

func (m *UploadMonitor) Start(c context.Context) {
	go m.monitorAndProcessUploads(c)
}

func (m *UploadMonitor) monitorAndProcessUploads(c context.Context) {
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
	lastWeek := getLastWeek()
	tasks, err := m.client.GetEndpointTaskList(m.endpointID, map[string]string{
		"filter_completion_time": lastWeek,
		"filter_status":          "SUCCEEDED",
	})

	if err != nil {
		log.Infof("globus.GetEndpointTaskList returned the following error: %s - %#v", err, m.client.GetGlobusErrorResponse())
		return
	}

	for _, task := range tasks.Tasks {
		log.Infof("Getting successful transfers for Globus Task %s", task.TaskID)
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
		}
	}
}

func getLastWeek() string {
	now := time.Now()
	now.AddDate(0, 0, -7)
	return now.Format("2006-01-02")
}

func (m *UploadMonitor) processTransfers(transfers *TransferItems) {
	transferItem := transfers.Transfers[0]

	// Destination path will have the following format: /__globus_uploads/<id of upload request>/...rest of path...
	// So the second entry in the array is the id in the globus_uploads table we want to look up.
	pieces := strings.Split(transferItem.DestinationPath, "/")
	if len(pieces) < 3 {
		// sanity check, because the destination path should at least be /__globus_uploads/<id>/...rest of path...
		// thus should at least have 3 entries in it
		log.Infof("Invalid globus DestinationPath: %s", transferItem.DestinationPath)
		return
	}

	id := pieces[1]
	globusUpload, err := m.globusUploads.GetGlobusUpload(id)
	if err != nil {
		// Upload is already being processed
		return
	}

	log.Infof("Process globus upload %s", id)

	// At this point we have a globus upload. What we are going to do is remove the ACL on the directory
	// so no more files can be uploaded to it. Then we are going to add that directory to the list of
	// directories to upload. Then the file loader will eventually get around to loading these files. In
	// the meantime since we've now created a file load from this globus upload we can delete the entry
	// from the globus_uploads table.

	if _, err := m.client.DeleteEndpointACLRule(m.endpointID, globusUpload.GlobusAclID); err != nil {
		log.Infof("Unable to delete ACL: %s", err)
	}

	flAdd := model.AddFileLoadModel{
		ProjectID: globusUpload.ProjectID,
		Owner:     globusUpload.Owner,
		Path:      globusUpload.Path,
	}

	if fl, err := m.fileLoads.AddFileLoad(flAdd); err != nil {
		log.Infof("Unable to add file load request: %s", err)
		return
	} else {
		log.Infof("Created file load (id: %s) for globus upload %s", fl.ID, id)
	}

	// Delete the globus upload request as we have now turned it into a file loading request
	// and won't have to process this request again. If the server stops while loading the
	// request or their is some other failure, the file loader will take care of picking up
	// where it left off.
	m.globusUploads.DeleteGlobusUpload(id)
}
