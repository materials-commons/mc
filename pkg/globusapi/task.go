package globusapi

type TaskExtras struct {
	SourceEndpoint                     string `json:"source_endpoint"`
	SourceEndpointID                   string `json:"source_endpoint_id"`
	SourceEndpointDisplayName          string `json:"source_endpoint_display_name"`
	DestinationEndpoint                string `json:"destination_endpoint"`
	DestinationEndpointID              string `json:"destination_endpoint_id"`
	DestinationEndpointDisplayName     string `json:"destination_endpoint_display_name"`
	SourceHostEndpoint                 string `json:"source_host_endpoint"`
	SourceHostEndpointID               string `json:"source_host_endpoint_id"`
	SourceHostEndpointDisplayName      string `json:"source_host_endpoint_display_name"`
	DestinationHostEndpoint            string `json:"destination_host_endpoint"`
	DestinationHostEndpointID          string `json:"destination_host_endpoint_id"`
	DestinationHostEndpointDisplayName string `json:"destination_host_endpoint_display_name"`
	SourceHostPath                     string `json:"source_host_path"`
	DestinationHostPath                string `json:"destination_host_path"`
	IsOK                               bool   `json:"is_ok"`
	SourceLocalUser                    string `json:"source_local_user"`
	SourceLocalUserStatus              string `json:"source_local_user_status"`
	DestinationLocalUser               string `json:"destination_local_user"`
	DestinationLocalUserStatus         string `json:"destination_local_user"`
	OwnerString                        string `json:"owner_string"`
}

type Task struct {
	DataType                   string `json:"DATA_TYPE"`
	TaskID                     string `json:"task_id"`
	Type                       string `json:"type"`
	Status                     string `json:"status"`
	Label                      string `json:"label"`
	OwnerID                    string `json:"owner_id"`
	RequestTime                string `json:"request_time"`
	CompletionTime             string `json:"completion_time"`
	Deadline                   string `json:"deadline"`
	SyncLevel                  int    `json:"sync_level"`
	EncryptData                bool   `json:"encrypt_data"`
	VerifyChecksum             bool   `json:"verify_checksum"`
	DeleteDestinationExtra     bool   `json:"delete_destination_extra"`
	RecursiveSymlinks          string `json:"recursive_symlinks"`
	PreserveTimestamp          bool   `json:"preserve_timestamp"`
	Command                    string `json:"command"`
	HistoryDeleted             bool   `json:"history_deleted"`
	Faults                     int    `json:"faults"`
	FilesCount                 int    `json:"files"`
	DirectoriesCount           int    `json:"directories"`
	SymlinksCount              int    `json:"symlinks"`
	FilesSkipped               int    `json:"files_skipped"`
	FilesTransferred           int    `json:"files_transferred"`
	SubtasksTotal              int    `json:"subtasks_total"`
	SubtasksPending            int    `json:"subtasks_pending"`
	SubtasksRetrying           int    `json:"subtasks_retrying"`
	SubtasksSucceeded          int    `json:"subtasks_succeeded"`
	SubtasksFailed             int    `json:"subtasks_failed"`
	BytesTransferred           int    `json:"bytes_transferred"`
	BytesChecksummed           int    `json:"bytes_checksummed"`
	EffectiveBytesPerSecond    int    `json:"effective_bytes_per_second"`
	NiceStatus                 string `json:"nice_status"`
	NiceStatusShortDescription string `json:"nice_status_short_description"`
	NiceStatusExpiresIn        string `json:"nice_status_expires_in"`
	CanceledByAdmin            string `json:"canceled_by_admin"`
	CanceledByAdminMessage     string `json:"canceled_by_admin_message"`
	IsPaused                   bool   `json:"is_paused"`
	TaskExtras
}

type TaskList struct {
	DataType    string `json:"DATA_TYPE"`
	Limit       int    `json:"limit"`
	LastKey     string `json:"last_key"`
	HasNextPage bool   `json:"has_next_page"`
	Tasks       []Task `json:"DATA"`
}

type TransferItems struct {
	DataType   string     `json:"DATA_TYPE"`
	Marker     int        `json:"marker"`
	NextMarker int        `json:"next_marker"`
	Transfers  []Transfer `json:"DATA"`
}

type Transfer struct {
	DataType        string `json:"DATA_TYPE"`
	SourcePath      string `json:"source_path"`
	DestinationPath string `json:"destination_path"`
}
