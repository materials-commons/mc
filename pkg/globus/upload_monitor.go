package globus

import (
	"context"
	"time"
)

type UploadMonitor struct {
	client     *Client
	endpointID string
}

func NewUploadMonitor(client *Client, endpointID string) *UploadMonitor {
	return &UploadMonitor{client: client, endpointID: endpointID}
}

func (m *UploadMonitor) Start(c context.Context) {
	go m.monitorAndProcessUploads(c)
}

func (m *UploadMonitor) monitorAndProcessUploads(c context.Context) {
	for {
		m.retrieveAndProcessUploads()
		select {
		case <-c.Done():
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func (m *UploadMonitor) retrieveAndProcessUploads() {
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
	}
}

func (m *UploadMonitor) processTransfers(transfers *TransferItems) {
	//transferItem := transfers.Transfers[0]
	// 1. Determine upload id from dir path

	// 2. Lookup the upload id
	// 2.a if upload id is already being processed then skip it
	// 2.b other mark it as being processed and delete ACL (acl id in the upload record)

	// 4. Load the files
	// 5. Remove the top level directory

	// Some general thoughts -
	// Most of this logic is already in place in the background loader, so could just
	// add the item into the file uploads table and let the background loader take
	// care of loading the files.
	//
	// If so, then logic is really simple:
	//     Look in globus table for request
	//     if no found then just ignore
	//     else {
	//        Delete ACL
	//        Add into file uploads table
	//        Delete globus request in table
	//     }
	//
	// With this logic all cancellation, restarts, etc... will be handled in the file loading logic
	// which already exists.

	// GlobusResponse({
	//    'DATA_TYPE': 'successful_transfer',
	//    'destination_path': '/__upload_staging/transfer-5ac039c9-6254-4b39-90db-32b89ba6b5a9/hello.titan.txt',
	//    'source_path': None
	// })
}
