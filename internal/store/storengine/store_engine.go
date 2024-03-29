package storengine

import (
	"time"

	"github.com/materials-commons/mc/internal/store/model"
)

type UsersStoreEngine interface {
	AddUser(user model.UserSchema) (model.UserSchema, error)
	GetUserByID(id string) (model.UserSchema, error)
	GetUserByAPIKey(apikey string) (model.UserSchema, error)
	UpdateUserFullname(id, fullname string, updatedAt time.Time) (model.UserSchema, error)
	UpdateUserPassword(id, password string, updatedAt time.Time) (model.UserSchema, error)
	UpdateUserAPIKey(id, apikey string, updatedAt time.Time) (model.UserSchema, error)
	UpdateUserGlobusUser(id string, globusUser string) error
}

type ProjectsStoreEngine interface {
	GetProjectSimple(id string) (model.ProjectSimpleModel, error)
	AddProject(project model.ProjectSchema) (model.ProjectSchema, error)
	DeleteProject(id string) error
	GetProjectUsers(id string) ([]model.UserSchema, error)
}

type AccessStoreEngine interface {
	AddAccessEntry(entry model.ProjectAccessSchema) (model.ProjectAccessSchema, error)
	DeleteAccess(projectID, userID string) error
	DeleteAllAccessForProject(projectID string) error
	GetProjectAccessEntries(projectID string) ([]model.ProjectAccessSchema, error)
	GetUserAccessEntries(userID string) ([]model.ProjectAccessSchema, error)
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

type SamplesStoreEngine interface {
	AddSample(sample model.SampleSchema) (model.SampleSchema, error)
	DeleteSample(sampleID string) error
	GetSample(sampleID string) (model.SampleSchema, error)
	ModifySampleName(sampleID, name string, updatedAt time.Time) error
}

type ProcessesStoreEngine interface {
	AddProcess(process model.ProcessSchema) (model.ProcessSchema, error)
	GetProcess(processID string) (model.ProjectExtendedModel, error)
}

type AssociationsStoreEngine interface {
	AssociateSampleWithProject(sampleID, projectID string) error
	AssociateSampleWithExperiment(sampleID, experimentID string) error
	AssociateSampleWithFile(sampleID, fileID string) error
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

type BackgroundProcessStoreEngine interface {
	AddBackgroundProcess(bgp model.BackgroundProcessSchema) (model.BackgroundProcessSchema, error)
	GetListBackgroundProcess(glbpg model.GetListBackgroundProcessModel) ([]model.BackgroundProcessSchema, error)
	DeleteBackgroundProcess(id string) error
	GetBackgroundProcess(id string) (model.BackgroundProcessSchema, error)
	UpdateStatusBackgroundProcess(id string, status string, message string) error
	SetFinishedBackgroundProcess(id string, done bool) error
	SetOkBackgroundProcess(id string, success bool) error
}

type ExperimentsStoreEngine interface {
	GetExperimentOverviewsForProject(projectID string) ([]model.ExperimentOverviewModel, error)
}

type DatasetsStoreEngine interface {
	GetDatadirsForDataset(datasetID string) ([]model.DatadirEntryModel, error)
	GetDataset(datasetID string) (model.DatasetSchema, error)
	SetDatasetZipfile(datasetID string, size int64, name string) error
}
