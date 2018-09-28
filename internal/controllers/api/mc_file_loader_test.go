package api_test

import (
	"os"
	"testing"

	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestMCFileLoaderLoadDirectory(t *testing.T) {
	finfo, err := os.Stat(".")
	assert.Okf(t, err, "Unable to stat current dir %s", err)

	var project store.ProjectSimpleModel
	project.Name = "My Project"
	mcFileLoader := api.NewMCFileLoader("/tmp", project, store.NewDatafilesStore(nil))
	mcFileLoader.LoadFileOrDir("/tmp/dir", finfo)
}
