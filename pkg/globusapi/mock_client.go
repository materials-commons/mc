package globusapi

type MockClient struct {
	err error
}

func NewMockClient(err error) *MockClient {
	return &MockClient{err: err}
}

func (c *MockClient) Authenticate() error { return c.err }

func (c *MockClient) GetEndpointTaskList(endpointID string, filters map[string]string) (TaskList, error) {
	return TaskList{}, c.err
}

func (c *MockClient) GetTaskSuccessfulTransfers(taskID string, marker int) (TransferItems, error) {
	return TransferItems{}, c.err
}

func (c *MockClient) GetIdentities(users []string) (Identities, error) { return Identities{}, c.err }

func (c *MockClient) AddEndpointACLRule(rule EndpointACLRule) (AddEndpointACLRuleResult, error) {
	return AddEndpointACLRuleResult{}, c.err
}

func (c *MockClient) DeleteEndpointACLRule(endpointID string, accessID int) (DeleteEndpointACLRuleResult, error) {
	return DeleteEndpointACLRuleResult{}, c.err
}

func (c *MockClient) Err(err error) *MockClient {
	c.err = err
	return c
}
