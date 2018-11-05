package globus

type endpointACLRequest struct {
	DataType      string `json:"DATA_TYPE"`
	PrincipalType string `json:"principal_type"`
	Principal     string `json:"principal"`
	Path          string `json:"path"`
	Permissions   string `json:"permissions"`
	NotifyEmail   string `json:"notify_email,omitempty"`
}

type EndpointACLRule struct {
	EndpointID  string
	Path        string
	IdentityID  string
	Permissions string
}

type AddEndpointACLRuleResult struct {
	Code      string `json:"code"`
	Resource  string `json:"resource"`
	DataType  string `json:"DATA_TYPE"`
	RequestID string `json:"request_id"`
	AccessID  string `json:"access_id"`
	Message   string `json:"message"`
}

/*
{
    "message": "Access rule '123' deleted successfully",
    "code": "Deleted",
    "resource": "/endpoint/user#ep1/access/123",
    "DATA_TYPE": "result",
    "request_id": "ABCdef789"
}
*/
type DeleteEndpointACLRuleResult struct {
	Code      string `json:"code"`
	DataType  string `json:"DATA_TYPE"`
	Resource  string `json:"resource"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}
