package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"
	r "gopkg.in/gorethink/gorethink.v4"
)

type ExperimentsRethinkdb struct {
	Session *r.Session
}

func NewExperimentsRethinkdb(session *r.Session) *ExperimentsRethinkdb {
	return &ExperimentsRethinkdb{Session: session}
}

func (e *ExperimentsRethinkdb) GetExperimentOverviewsForProject(projectID string) ([]model.ExperimentOverviewModel, error) {
	var experiments []model.ExperimentOverviewModel
	errMsg := fmt.Sprintf("Can't retrieve experiments for project %s", projectID)
	res, err := r.Table("project2experiment").GetAllByIndex("project_id", projectID).
		EqJoin("experiment_id", r.Table("experiments")).Zip().
		Merge(experimentOverview).CoerceTo("array").Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return experiments, err
	}

	defer res.Close()
	err = res.All(&experiments)
	return experiments, err
}

func experimentOverview(e r.Term) interface{} {
	return map[string]interface{}{
		"owner_details": r.Table("users").Get(e.Field("owner")).Pluck("fullname"),
		"files_count":   r.Table("experiment2datafile").GetAllByIndex("experiment_id", e.Field("id")).Count(),
		"samples_count": r.Table("experiment2sample").GetAllByIndex("experiment_id", e.Field("id")).Count(),
	}
}
