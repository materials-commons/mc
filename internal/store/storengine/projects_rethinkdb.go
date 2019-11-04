package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
)

type ProjectsRethinkdb struct {
	Session *r.Session
}

func NewProjectsRethinkdb(session *r.Session) *ProjectsRethinkdb {
	return &ProjectsRethinkdb{Session: session}
}

func (e *ProjectsRethinkdb) GetProjectSimple(id string) (model.ProjectSimpleModel, error) {
	var project model.ProjectSimpleModel
	errMsg := fmt.Sprintf("No such project %s", id)
	res, err := r.Table("projects").Get(id).Merge(projectTopLevelDir).Run(e.Session)

	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return project, err
	}
	defer res.Close()

	err = res.One(&project)
	return project, err
}

func projectTopLevelDir(p r.Term) interface{} {
	return map[string]interface{}{
		"root_dir": r.Table("datadirs").
			GetAllByIndex("datadir_project_name", []interface{}{p.Field("id"), p.Field("name")}).CoerceTo("array"),
	}
}

func (e *ProjectsRethinkdb) GetProjectUsers(id string) ([]model.UserSchema, error) {
	users := make([]model.UserSchema, 0)
	errMsg := fmt.Sprintf("No such project %s", id)
	res, err := r.Table("access").GetAllByIndex("project_id", id).
		EqJoin("user_id", r.Table("users")).Zip().Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return users, err
	}

	defer res.Close()

	err = res.All(&users)
	return users, err
}
