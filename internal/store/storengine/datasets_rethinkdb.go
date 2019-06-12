package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"
	r "gopkg.in/gorethink/gorethink.v4"
)

type DatasetsRethinkdb struct {
	Session *r.Session
}

func NewDatasetsRethinkdb(session *r.Session) *DatasetsRethinkdb {
	return &DatasetsRethinkdb{Session: session}
}

func (e *DatasetsRethinkdb) GetDatadirsForDataset(datasetID string) ([]model.DatadirEntryModel, error) {
	datadirs := make([]model.DatadirEntryModel, 0)
	errMsg := fmt.Sprintf("Unable to retrieve directories for dataset %s", datasetID)
	res, err := r.Table("dataset2datafile").GetAllByIndex("dataset_id", datasetID).
		EqJoin([]interface{}{r.Row.Field("datafile_id")}, r.Table("datadir2datafile"), r.EqJoinOpts{Index: "datadir_id"}).Zip().
		Distinct(r.DistinctOpts{Index: "datadir_id"}).
		EqJoin([]interface{}{r.Row.Field("datadir_id")}, r.Table("datadirs")).Zip().Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return datadirs, err
	}
	defer res.Close()
	err = res.All(&datadirs)
	return datadirs, err
}

func (e *DatasetsRethinkdb) GetDataset(datasetID string) (model.DatasetSchema, error) {
	var dataset model.DatasetSchema
	errMsg := fmt.Sprintf("No such dataset %s", datasetID)
	res, err := r.Table("datasets").Get(datasetID).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return dataset, err
	}
	defer res.Close()

	err = res.One(&dataset)

	if err != nil {
		return dataset, err
	}

	if dataset.SelectionID != "" {
		var fileSelection model.FileSelection
		errMsg := fmt.Sprintf("No such Selection %s", dataset.SelectionID)

		res2, err := r.Table("fileselection").Get(dataset.SelectionID).Run(e.Session)
		if err := checkRethinkdbQueryError(res2, err, errMsg); err != nil {
			return dataset, err
		}
		defer res2.Close()
		err = res2.One(&fileSelection)
		dataset.FileSelection = fileSelection
		return dataset, err
	}

	return dataset, err
}
