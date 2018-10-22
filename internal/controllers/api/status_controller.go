package api

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo"
)

type StatusController struct {
}

func NewStatusController() *StatusController {
	return &StatusController{}
}

type ServerStatus struct {
	NumCPUs       int    `json:"num_cpus"`
	NumGoRoutines int    `json:"num_go_routines"`
	GoVersion     string `json:"go_version"`
	Memory        uint64 `json:"memory"`
}

func (s *StatusController) GetServerStatus(c echo.Context) error {
	status := ServerStatus{
		NumCPUs:       runtime.NumCPU(),
		NumGoRoutines: runtime.NumGoroutine(),
		GoVersion:     runtime.Version(),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status.Memory = m.Alloc

	return c.JSON(http.StatusOK, status)
}
