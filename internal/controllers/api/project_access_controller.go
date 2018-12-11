package api

import "github.com/labstack/echo"

type ProjectAccessController struct {
}

func NewProjectAccessController() *ProjectAccessController {
	return &ProjectAccessController{}
}

func (p *ProjectAccessController) GetProjectUserAccessEntries(c echo.Context) error {
	return nil
}

func (p *ProjectAccessController) AddUserAccessEntryToProject(c echo.Context) error {
	return nil
}

func (p *ProjectAccessController) DeleteUserAccessEntryFromProject(c echo.Context) error {
	return nil
}
