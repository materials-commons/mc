package store

import "time"

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
	ModifyUserFullname(id, fullname string, updatedAt time.Time) (UserSchema, error)
	ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error)
	ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error)
}

type ProjectsStoreEngine interface {
	AddProject(project ProjectSchema) (ProjectSchema, error)
	GetProject(id string) (ProjectExtendedModel, error)
	GetProjectSimple(id string) (ProjectSimpleModel, error)
	GetAllProjectsForUser(user string) ([]ProjectExtendedModel, error)
	DeleteProject(id string) error
	UpdateProjectName(id string, name string, updatedAt time.Time) error
	UpdateProjectDescription(id string, description string, updatedAt time.Time) error
}

type AccessStoreEngine interface {
	AddAccessEntry(entry AccessSchema) (AccessSchema, error)
	DeleteAccess(projectID, userID string) error
	DeleteAllAccessForProject(projectID string) error
	GetProjectAccessEntries(projectID string) ([]AccessSchema, error)
	GetUserAccessEntries(userID string) ([]AccessSchema, error)
}

type DatafilesStoreEngine interface {
	AddFile(file DatafileSchema, projectID, datadirID string) (DatafileSchema, error)
	GetFile(id string) (DatafileSchema, error)
	GetFileWithChecksum(checksum string) (DatafileSchema, error)
	GetFileInDir(name string, dirID string) (DatafileSchema, error)
	UpdateFileCurrentFlag(fileID string, current bool) error
}

type DatadirsStoreEngine interface {
	AddDir(dir DatadirSchema) (DatadirSchema, error)
	GetDatadirByPathInProject(path, projectID string) (DatadirSchema, error)
	GetDatadir(id string) (DatadirSchema, error)
}

type SamplesStoreEngine interface {
	AddSample(sample SampleSchema) (SampleSchema, error)
	DeleteSample(sampleID string) error
	GetSample(sampleID string) (SampleSchema, error)
	ModifySampleName(sampleID, name string, updatedAt time.Time) error
}

type ProcessesStoreEngine interface {
	AddProcess(process ProcessSchema) (ProcessSchema, error)
	GetProcess(processID string) (ProjectExtendedModel, error)
}

type AssociationsStoreEngine interface {
	AssociateSampleWithProject(sampleID, projectID string) error
	AssociateSampleWithExperiment(sampleID, experimentID string) error
	AssociateSampleWithFile(sampleID, fileID string) error
}

type FileLoadsStoreEngine interface {
	AddFileLoad(fileLoad FileLoadSchema) (FileLoadSchema, error)
	DeleteFileLoad(id string) error
	GetFileLoad(id string) (FileLoadSchema, error)
	GetAllFileLoads() ([]FileLoadSchema, error)
	MarkAllNotLoading() error
	UpdateLoading(id string, loading bool) error
}
