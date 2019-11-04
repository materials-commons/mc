package storengine

import (
	"github.com/materials-commons/mc/internal/store/model"
)

type UsersStoreEngine interface {
	AddUser(user model.UserSchema) (model.UserSchema, error)
	GetUserByID(id string) (model.UserSchema, error)
	GetUserByAPIKey(apikey string) (model.UserSchema, error)
}

type ProjectsStoreEngine interface {
	GetProjectSimple(id string) (model.ProjectSimpleModel, error)
	GetProjectUsers(id string) ([]model.UserSchema, error)
}

type DatafilesStoreEngine interface {
	AddFile(file model.DatafileSchema, projectID, datadirID string) (model.DatafileSchema, error)
	GetFile(id string) (model.DatafileSchema, error)
	GetFileWithChecksum(checksum string) (model.DatafileSchema, error)
	GetFileInDir(name string, dirID string) (model.DatafileSchema, error)
	UpdateFileCurrentFlag(fileID string, current bool) error
}

type DatadirsStoreEngine interface {
	GetFilesForDatadir(projectID, userID, dirID string) ([]model.DatafileSimpleModel, error)
	GetDatadirForProject(projectID, userID, dirID string) (model.DatadirEntryModel, error)
	GetDatadirsForProject(projectID, userID string) ([]model.DatadirEntryModel, error)
	//
	AddDir(dir model.DatadirSchema) (model.DatadirSchema, error)
	GetDatadirByPathInProject(path, projectID string) (model.DatadirSchema, error)
	GetDatadir(id string) (model.DatadirSchema, error)
}

type FileLoadsStoreEngine interface {
	AddFileLoad(fileLoad model.FileLoadSchema) (model.FileLoadSchema, error)
	DeleteFileLoad(id string) error
	GetFileLoad(id string) (model.FileLoadSchema, error)
	GetAllFileLoads() ([]model.FileLoadSchema, error)
	MarkAllNotLoading() error
	UpdateLoading(id string, loading bool) error
}

type GlobusUploadsStoreEngine interface {
	AddGlobusUpload(upload model.GlobusUploadSchema) (model.GlobusUploadSchema, error)
	DeleteGlobusUpload(id string) error
	GetGlobusUpload(id string) (model.GlobusUploadSchema, error)
	GetAllGlobusUploads() ([]model.GlobusUploadSchema, error)
	GetAllGlobusUploadsForUser(user string) ([]model.GlobusUploadSchema, error)
}

type DatasetsStoreEngine interface {
	GetDatadirsForDataset(datasetID string) ([]model.DatadirEntryModel, error)
	GetDataset(datasetID string) (model.DatasetSchema, error)
	SetDatasetZipfile(datasetID string, size int64, name string) error
}
