package globusapi

type APIClient interface {
	Authenticate() error
	GetEndpointTaskList(endpointID string, filters map[string]string) (TaskList, error)
	GetTaskSuccessfulTransfers(taskID string, marker int) (TransferItems, error)
	GetIdentities(users []string) (Identities, error)
	AddEndpointACLRule(rule EndpointACLRule) (AddEndpointACLRuleResult, error)
	DeleteEndpointACLRule(endpointID string, accessID int) (DeleteEndpointACLRuleResult, error)
}
