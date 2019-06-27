package globusapi

type endpointACLRequest struct {
	DataType      string `json:"DATA_TYPE"`
	PrincipalType string `json:"principal_type"`
	Principal     string `json:"principal"`
	Path          string `json:"path"`
	Permissions   string `json:"permissions"`
	NotifyEmail   string `json:"notify_email,omitempty"`
}

const ACLPrincipalTypeIdentity = "identity"
const ACLPrincipalTypeAllAuthenticatedUsers = "all_authenticated_users"

type EndpointACLRule struct {
	PrincipalType string
	EndpointID    string
	Path          string
	IdentityID    string
	Permissions   string
}

type AddEndpointACLRuleResult struct {
	Code      string `json:"code"`
	Resource  string `json:"resource"`
	DataType  string `json:"DATA_TYPE"`
	RequestID string `json:"request_id"`
	AccessID  string `json:"access_id"`
	Message   string `json:"message"`
}

type DeleteEndpointACLRuleResult struct {
	Code      string `json:"code"`
	DataType  string `json:"DATA_TYPE"`
	Resource  string `json:"resource"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type EndpointAccessRuleList struct {
	Length      int          `json:"length"`
	Endpoint    string       `json:"endpoint"`
	AccessRules []AccessRule `json:"DATA"`
	DataType    string       `json:"DATA_TYPE"`
}

type AccessRule struct {
	DataType      string `json:"DATA_TYPE"`
	PrincipalType string `json:"principal_type"`
	Path          string `json:"path"`
	Principal     string `json:"principal"`
	AccessID      string `json:"id"`
	Permissions   string `json:"permissions"`
	// There are also role_id and role_type fields. Its not clear what
	// their type is and since we don't need them they are ignored.
}
