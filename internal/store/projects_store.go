package store

type ProjectsStore struct {
	ProjectsStoreEngine
}

func NewProjectsStore(e ProjectsStoreEngine) *ProjectsStore {
	return &ProjectsStore{ProjectsStoreEngine: e}
}

func (s *ProjectsStore) GetProjectSimple(id string) (ProjectSimpleModel, error) {
	return s.ProjectsStoreEngine.GetProjectSimple(id)
}
