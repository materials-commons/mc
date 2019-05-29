package file

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
)

type DownloadDir struct {
	ProjectID  string
	User       model.UserSchema
	ddirsStore *store.DatadirsStore
}

func NewDownloadDir(projectID string, user model.UserSchema, ddirsStore *store.DatadirsStore) *DownloadDir {
	return &DownloadDir{ProjectID: projectID, User: user, ddirsStore: ddirsStore}
}

// CreateDownloadDirectory creates a download directory for the given set of directories. It will
// walk through the directories and create them in this temporary download. It will
// then create a link to in this directory that points at the file entry in the Materials Commons
// store. The reason this needs to be done is that the Materials Commons store is an object store
// (like S3), where as Globus (and users) need to see the imposed directory structure. This is
// reconstructed from the database. Links to files are used so we don't have to create copies
// of the files.
func (d *DownloadDir) CreateDownloadDirectory(basePath string, ddirs []model.DatadirEntryModel) error {
	log.Infof("CreateDownloadDirectory %s", basePath)
	dirs := d.dirsToCreate(basePath, ddirs)

	for _, dir := range dirs {
		log.Infof("os.MkdirAll %s", dir)
		if err := os.MkdirAll(dir, 0700); err != nil {
			log.Infof("Failed to create directory %s", dir)
		}
	}

	for _, dir := range ddirs {
		files, err := d.ddirsStore.GetFilesForDatadir(d.ProjectID, d.User.ID, dir.ID)
		if err != nil && errors.Cause(err) != mc.ErrNotFound {
			log.Infof("GetFilesForDatadir(%s, %s, %s) failed: %s", d.ProjectID, d.User.ID, dir.ID, err)
		}

		for _, file := range files {
			linkToPath := filepath.Join(basePath, dir.Name, file.Name)
			if err := os.Link(file.FirstMCDirPath(), linkToPath); err != nil {
				log.Infof("Failed to create hard link for file %s/%s", file.FirstMCDirPath(), linkToPath, err)
			}
		}
	}

	return nil
}

func (d *DownloadDir) dirsToCreate(basePath string, ddirs []model.DatadirEntryModel) []string {
	var dirs []string

	for _, dir := range ddirs {
		dirs = append(dirs, filepath.Join(basePath, dir.Name))
	}

	return minimumSetOfDirsToCreate(dirs)
}

// minimumSetOfDirsToCreate creates the minimum set of dirs to call MkdirAll on.
// It takes a set of directories and figures out the set of deepest paths that
// encompasses all directories. For example given:
//     dirs = [ "/a", "/a/c", "/a/b", "/a/c/d", "/a/b/c", "/a/d"]
// It will return:
//     ["/a/b/c", "/a/c/d", "/a/d"]
// This list encompasses the minimum  the dirs to be created in a series of MkdirAll
// that will create all the directories.
func minimumSetOfDirsToCreate(dirs []string) []string {
	// The algorithm to get the minimum set of dirs works as follows:
	// 1. Sort the list of dirs
	// 2. Go through list of sorted dirs and check if the previous entry is
	//    contained in the current directory. If it is then set previousDir to
	//    the current entry. If it *is not* then we are looking at a new path
	//    so add previous entry to the list of directories that make up the
	//    minimum set, and set previous dir to current entry.
	// 3. The minimum set is a hash list. We could get a way with a list, but
	//    the last step is to check if the last entry in the sorted list is
	//    in minimum list. If it isn't then add it.
	// 4. Return the keys in the array
	sort.Strings(dirs)
	dirsToKeep := make(map[string]bool)
	previousDir := dirs[0]
	for _, dir := range dirs {
		if strings.Contains(dir, previousDir) {
			previousDir = dir
		} else {
			dirsToKeep[previousDir] = true
			previousDir = dir
		}
	}

	lastDir := dirs[len(dirs)-1]
	if _, ok := dirsToKeep[lastDir]; !ok {
		dirsToKeep[lastDir] = true
	}

	keys := make([]string, 0, len(dirsToKeep))
	for k := range dirsToKeep {
		keys = append(keys, k)
	}

	return keys
}
