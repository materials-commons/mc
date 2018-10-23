package storengine

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
)

type AssociationsRethinkdb struct {
	Session *r.Session
}

func NewAssociationsRethinkdb(session *r.Session) *AssociationsRethinkdb {
	return &AssociationsRethinkdb{Session: session}
}

func (e *AssociationsRethinkdb) AssociateSampleWithProject(sampleID, projectID string) error {
	errMsg := fmt.Sprintf("Unable to associate sample %s with project %s", sampleID, projectID)
	resp, err := r.Table("project2sample").
		Insert(map[string]interface{}{"sample_id": sampleID, "project_id": projectID}, r.InsertOpts{ReturnChanges: true}).
		RunWrite(e.Session)
	return checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *AssociationsRethinkdb) AssociateSampleWithExperiment(sampleID, experimentID string) error {
	errMsg := fmt.Sprintf("Unable to associate sample %s with experiment %s", sampleID, experimentID)
	resp, err := r.Table("experiment2sample").
		Insert(map[string]interface{}{"sample_id": sampleID, "experiment_id": experimentID}, r.InsertOpts{ReturnChanges: true}).
		RunWrite(e.Session)
	return checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *AssociationsRethinkdb) AssociateSampleWithFile(sampleID, fileID string) error {
	errMsg := fmt.Sprintf("Unable to associate sample %s with file %s", sampleID, fileID)
	resp, err := r.Table("sample2datafile").
		Insert(map[string]interface{}{"sample_id": sampleID, "datafile_id": fileID}, r.InsertOpts{ReturnChanges: true}).
		RunWrite(e.Session)
	return checkRethinkdbInsertError(resp, err, errMsg)
}
