package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
	"gopkg.in/gorethink/gorethink.v4/encoding"
)

type ProjectsRethinkdb struct {
	Session *r.Session
}

func NewProjectsRethinkdb(session *r.Session) *ProjectsRethinkdb {
	return &ProjectsRethinkdb{Session: session}
}

func (e *ProjectsRethinkdb) AddProject(project model.ProjectSchema) (model.ProjectSchema, error) {
	errMsg := fmt.Sprintf("Unable to add project %+v", project)
	resp, err := r.Table("projects").Insert(project, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	if err := checkRethinkdbInsertError(resp, err, errMsg); err != nil {
		return project, err
	}

	var proj model.ProjectSchema
	if err := encoding.Decode(&proj, resp.Changes[0].NewValue); err != nil {
		return proj, err
	}

	ddirModel := model.AddDatadirModel{
		Name:      project.Name,
		Owner:     project.Owner,
		ProjectID: proj.ID,
	}

	_, err = AddDatadir(ToDatadirSchema(ddirModel), e.Session)

	return proj, err
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

func (e *ProjectsRethinkdb) DeleteProject(id string) error {
	errMsg := fmt.Sprintf("failed deleting project %s", id)

	resp, err := r.Table("projects").Get(id).
		Update(map[string]interface{}{"owner": "delete@materialscommons.org"}).RunWrite(e.Session)

	if err := checkRethinkdbUpdateError(resp, err, errMsg); err != nil {
		return err
	}

	resp, err = r.Table("access").GetAllByIndex("project_id", id).Delete().RunWrite(e.Session)
	return checkRethinkdbWriteError(resp, err, errMsg)
}
