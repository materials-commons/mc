package store

import (
	"fmt"
	"time"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type SamplesStoreEngineRethinkdb struct {
	Session *r.Session
}

func NewSamplesStoreEngineRethinkdb(session *r.Session) *SamplesStoreEngineRethinkdb {
	return &SamplesStoreEngineRethinkdb{Session: session}
}

func (e *SamplesStoreEngineRethinkdb) AddSample(sample SampleSchema) (SampleSchema, error) {
	errMsg := fmt.Sprintf("Unable to add sample %+v", sample)

	resp, err := r.Table("samples").Insert(sample, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return sample, err
	}

	var createdSample SampleSchema
	err = encoding.Decode(&createdSample, resp.Changes[0].NewValue)
	return createdSample, err
}

func (e *SamplesStoreEngineRethinkdb) DeleteSample(sampleID string) error {
	errMsg := fmt.Sprintf("Unable to delete sample %s", sampleID)
	resp, err := r.Table("samples").Get(sampleID).Delete().RunWrite(e.Session)
	return checkRethinkdbDeleteError(resp, err, errMsg)
}

func (e *SamplesStoreEngineRethinkdb) GetSample(sampleID string) (SampleSchema, error) {
	errMsg := fmt.Sprintf("No such sample %s", sampleID)
	res, err := r.Table("samples").Get(sampleID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return SampleSchema{}, err
	}

	var sample SampleSchema
	err = res.One(&sample)
	return sample, err
}

func (e *SamplesStoreEngineRethinkdb) ModifySampleName(sampleID, name string, updatedAt time.Time) error {
	errMsg := fmt.Sprintf("Unable to update sample %s name to %s", sampleID, name)
	resp, err := r.Table("samples").Get(sampleID).
		Update(map[string]interface{}{"name": name, "updated_at": updatedAt}).RunWrite(e.Session)
	return checkRethinkdbUpdateError(resp, err, errMsg)
}
