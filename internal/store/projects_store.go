package store

type ProjectsStore struct {
	pStoreEngine ProjectsStoreEngine
}

func NewProjectsStore(e ProjectsStoreEngine) *ProjectsStore {
	return &ProjectsStore{pStoreEngine: e}
}

func (s *ProjectsStore) GetProjectSimple(id string) (ProjectSimpleModel, error) {
	return s.pStoreEngine.GetProjectSimple(id)
}
