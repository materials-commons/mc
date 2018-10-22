package globus

import (
	"context"
	"fmt"
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

	marker := 0
	for _, task := range tasks.Tasks {
		for {
			transfers, err := m.client.GetTaskSuccessfulTransfers(task.TaskID, marker)
			if err != nil {
				continue
			}
			for _, transfer := range transfers.Transfers {
				fmt.Printf("transfer: %#v\n", transfer)
			}

			if transfers.NextMarker == 0 {
				break
			}

			marker = transfers.NextMarker
		}
	}
}
