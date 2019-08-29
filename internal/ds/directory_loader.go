package ds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"

	"github.com/materials-commons/mc/internal/store/model"
	r "gopkg.in/gorethink/gorethink.v4"
)

type DirLoader struct {
	basePath    string
	session     *r.Session
	createdDirs map[string]bool
}

func NewDirLoader(basePath string, session *r.Session) *DirLoader {
	return &DirLoader{basePath: basePath, session: session, createdDirs: make(map[string]bool)}
}

func (d *DirLoader) LoadDirFromDataset(dataset model.DatasetSchema, projectID string) error {
	selection := FromFileSelection(&dataset.FileSelection)
	fmt.Printf("selection = %#v\n", selection)

	return d.loadDatasetDir(projectID, dataset.ID, selection)
}

func (d *DirLoader) loadDatasetDir(projectID, datasetID string, selection *Selection) error {
	cursor, err := GetProjectDirsSortedCursor(projectID, d.session)
	if err != nil {
		return err
	}
	var dir model.DatadirSimpleModel
	for cursor.Next(&dir) {
		// Check if dir exists in selection, if not, then check its parent dir, and if that
		// exists set this dir to the parent dir setting. This reflects recursive selection as
		// parent directories that are included automatically include all descendants, and parent
		// directories that are excluded automatically exclude all descendants. These can be
		// overridden and selection will take that into account.
		if exists, _ := selection.DirExists(dir.Name); !exists {
			dirName := filepath.Dir(dir.Name)
			for {
				if dirName == "." {
					break
				}
				if exists, included := selection.DirExists(filepath.Dir(dirName)); exists {
					selection.AddDir(dir.Name, included)
				}
				dirName = filepath.Dir(dirName)
			}
		}

		fileCursor, err := GetDirFilesCursor(dir.ID, d.session)
		if err != nil {
			continue
		}

		var f model.DatafileSimpleModel
		for fileCursor.Next(&f) {
			fullMCFilePath := filepath.Join(dir.Name, f.Name)
			if selection.IsIncludedFile(fullMCFilePath) {
				dstDir := filepath.Join(d.basePath, dir.Name)
				if err := d.linkFile(f.FirstMCDirPath(), dstDir, f.Name); err != nil {
					log.Infof("Failed to create hard link from %s to %s/%s: %s", f.FirstMCDirPath(), dstDir, f.Name, err)
				}
			}
		}
	}
	return nil
}

// linkFile will create a hard link to a file and create any directories in the path.
func (d *DirLoader) linkFile(src, dstDir, fileName string) error {
	// Check if we need to create the directory
	if _, ok := d.createdDirs[dstDir]; !ok {
		d.createdDirs[dstDir] = true
		_ = os.MkdirAll(dstDir, 0700)
	}

	return os.Link(src, filepath.Join(dstDir, fileName))
}

func GetProjectDirsSortedCursor(projectID string, session *r.Session) (*r.Cursor, error) {
	return r.Table("project2datadir").GetAllByIndex("project_id", projectID).
		EqJoin("datadir_id", r.Table("datadirs")).Zip().
		OrderBy("name").
		Run(session)
}

func GetDirFilesCursor(dirID string, session *r.Session) (*r.Cursor, error) {
	return r.Table("datadir2datafile").GetAllByIndex("datadir_id", dirID).
		EqJoin("datafile_id", r.Table("datafiles")).Zip().
		Run(session)
}
