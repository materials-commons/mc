package model

type GlobusUploadSchema struct {
	ModelSimple
	Path             string `db:"path" json:"path" r:"path"`
	ProjectID        string `db:"project_id" json:"project_id" r:"project_id"`
	GlobusAclID      int    `db:"globus_acl_id" json:"globus_acl_id" r:"globus_acl_id"`
	GlobusEndpointID string `db:"globus_endpoint_id" json:"globus_endpoint_id" r:"globus_endpoint_id"`
	GlobusIdentityID string `db:"globus_identity_id" json:"globus_identity_id" r:"globus_identity_id"`
}

type AddGlobusUploadModel struct {
	Owner            string
	Path             string
	ProjectID        string
	GlobusAclID      int
	GlobusEndpointID string
	GlobusIdentityID string
}
