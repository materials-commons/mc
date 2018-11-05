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
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (m *UploadMonitor) retrieveAndProcessUploads(c context.Context) {
	yesterday := ""
	tasks, err := m.client.GetEndpointTaskList(m.endpointID, map[string]string{
		"filter_completion_time": yesterday,
		"filter_status":          "SUCCEEDED",
	})

	if err != nil {
		return
	}

	for _, task := range tasks.Tasks {
		transfers, err := m.client.GetTaskSuccessfulTransfers(task.TaskID, 0)

		switch {
		case err != nil:
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
		}
	}
}

func (m *UploadMonitor) processTransfers(transfers *TransferItems) {
	transferItem := transfers.Transfers[0]

	// Destination path will have the following format:
	// /__globus_uploads/<id of upload request>/...
	// So the second entry in the array is the id in the globus_uploads table we want to look up.
	pieces := strings.Split(transferItem.DestinationPath, "/")
	if len(pieces) < 3 {
		// sanity check, because the destination path should at least be /__globus_uploads/<id>/....
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

	if _, err := m.client.DeleteEndpointACLRule(m.endpointID, globusUpload.GlobusAclID); err != nil {
		log.Infof("Unable to delete ACL: %s", err)
	}

	flAdd := model.AddFileLoadModel{
		ProjectID: globusUpload.ProjectID,
		Owner:     globusUpload.Owner,
		Path:      globusUpload.Path,
	}

	if _, err := m.fileLoads.AddFileLoad(flAdd); err != nil {
		log.Infof("Unable to add upload request: %s", err)
		return
	}

	// Delete the globus upload request as we have not turned it into a file loading request
	// and won't have to process this request again.
	m.globusUploads.DeleteGlobusUpload(id)
}
