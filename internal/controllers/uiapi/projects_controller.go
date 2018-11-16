package uiapi

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
)

type ProjectsController struct {
	projectsStore *store.ProjectsStore
}

func NewProjectsController(db store.DB) *ProjectsController {
	return &ProjectsController{projectsStore: db.ProjectsStore()}
}

func (p *ProjectsController) GetProjectsForUser(c echo.Context) error {
	user := c.Get("User").(model.UserSchema)
	if projects, err := p.projectsStore.GetProjectsForUser(user.ID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, projects)
	}
}

type projectReq struct {
	ProjectID string `json:"project_id"`
}

func (p *ProjectsController) GetProjectOverview(c echo.Context) error {
	var req projectReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)
	if project, err := p.projectsStore.GetProjectOverview(req.ProjectID, user.ID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, project)
	}
}

func (p *ProjectsController) GetProjectAccessEntries(c echo.Context) error {
	var req projectReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	if users, err := p.projectsStore.GetProjectAccessEntries(req.ProjectID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, users)
	}
}

func (p *ProjectsController) AddUserToProject(c echo.Context) error {
	var req struct {
		ProjectID string `json:"project_id"`
		UserID    string `json:"user_id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)

	proj, err := p.projectsStore.GetProjectSimple(req.ProjectID)
	if err != nil {
		return err
	}

	if proj.Owner != user.ID {
		return echo.ErrUnauthorized
	}

	users, err := p.projectsStore.GetProjectAccessEntries(req.ProjectID)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.UserID == req.UserID {
			return echo.ErrForbidden
		}
	}

	if entry, err := p.projectsStore.AddUserToProject(req.ProjectID, req.UserID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusCreated, entry)
	}
}

func (p *ProjectsController) DeleteUserFromProject(c echo.Context) error {
	return nil
}
