/*
 ** MCFileLoader loads files and directories into a Materials Commons project. It does various checks to see
 ** if the file or directory already exists.
 */
package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"

	"github.com/pkg/errors"

	"github.com/materials-commons/mc/internal/store"
)

type MCFileLoader struct {
	root               string
	mcdir              string
	dfStore            *store.DatafilesStore
	ddStore            *store.DatadirsStore
	project            model.ProjectSimpleModel // Project we are adding entries to
	currentRootDatadir model.DatadirSchema      // Track which Materials Commons directory we are currently adding entries to
	owner              string                   // Owner here is the MC user we are performing the processing on behalf, and not the project owner
}

func NewMCFileLoader(root, owner, mcdir string, project model.ProjectSimpleModel, dfStore *store.DatafilesStore, ddStore *store.DatadirsStore) *MCFileLoader {
	return &MCFileLoader{
		root:    root,
		owner:   owner,
		dfStore: dfStore,
		ddStore: ddStore,
		project: project,
		mcdir:   mcdir,
	}
}

func (l *MCFileLoader) LoadFileOrDir(path string, finfo os.FileInfo) error {
	var err error
	// Get the current directory with root replaced by the project name. That is if
	//   root = /tmp/dir
	//   project name = My Project
	//   path = /tmp/dir/dir2
	// then dirName will be set to My Project/dir2
	// The reason for this is that directories are stored in Materials Commons with the
	// root equal to the project name.
	dirName := filepath.Dir(filepath.Join(l.project.Name, strings.TrimPrefix(path, l.root+"/")))
	if l.currentRootDatadir.Name != dirName {
		// The root directory changed, so get the directory from Materials Commons and set it as our "root" we are processing
		l.currentRootDatadir, err = l.ddStore.GetDatadirByPathInProject(dirName, l.project.ID)
		if err != nil {
			return errors.Wrapf(err, "Unable to get dir for path (%s) and project ID (%s) as new current dir", dirName, l.project.ID)
		}
	}

	switch {

	case !finfo.Mode().IsRegular() && !finfo.IsDir():
		// Not a regular file, just skip it
		return nil

	case finfo.IsDir():
		return l.loadDirectory(path, finfo)

	default:
		// Is a file
		return l.loadFile(path, finfo)
	}
}

// loadDirectory adds a directory to a project in Materials Commons if the directory doesn't exist. It takes
// care of constructing the name of the directory to include the project name as the root of the path.
func (l *MCFileLoader) loadDirectory(path string, finfo os.FileInfo) error {
	// Construct name as it exists in Materials Commons: project-name/dirpath, see comment above in LoadFileOrDir
	// on how and why this is done.
	dirPath := filepath.Join(l.project.Name, strings.TrimPrefix(path, l.root+"/"))
	if _, err := l.ddStore.GetDatadirByPathInProject(dirPath, l.project.ID); err == nil {
		return nil
	}

	dir := model.AddDatadirModel{
		Name:      dirPath,
		Owner:     l.owner,
		Parent:    l.currentRootDatadir.ID,
		ProjectID: l.project.ID,
	}

	if _, err := l.ddStore.AddDatadir(dir); err != nil {
		return errors.Wrapf(err, "Unable to create Datadir for path (%s) in project (%s)", dirPath, l.project.ID)
	}

	return nil
}

func (l *MCFileLoader) loadFile(path string, finfo os.FileInfo) error {
	var (
		checksum string
		df       model.DatafileSchema
		err      error
	)
	checksum, err = l.computeFileChecksum(path)
	if err != nil {
		return err
	}

	mediatype := GetMediaTypeByExtension(path)

	addFile := model.AddDatafileModel{
		Name:      filepath.Base(path),
		Owner:     l.owner,
		Checksum:  checksum,
		Size:      finfo.Size(),
		ProjectID: l.project.ID,
		DatadirID: l.currentRootDatadir.ID,
		MediaType: model.DatafileMediaType{
			Mime:        mediatype.Mime,
			Description: mediatype.Description,
		},
	}

	// See if a file with matching checksum already exists, and if so
	// set this entries UsesID to point at it.
	df, err = l.dfStore.GetDatafileWithChecksum(checksum)
	switch {
	case err != nil && errors.Cause(err) == mc.ErrNotFound:
		// Nothing to do, here for documentation purposes:
		// We didn't find a file with a matching checksum, so this is a brand
		// new file we are creating.
	case err != nil:
		return err
	default:
		// err == nil which means a matching entry was found
		// Set the UsesID so we refer to the existing file
		addFile.UsesID = df.ID
	}

	// See if a file with the same name exists in the directory, and if so
	// then set this files ParentID to point at it, and make the other file
	// not current.
	df, err = l.dfStore.GetDatafileInDir(addFile.Name, addFile.DatadirID)
	switch {
	case err != nil && errors.Cause(err) == mc.ErrNotFound:
		// Nothing to do, here for documentation purposes:
		// We didn't find a file with the same name in the given directory
		// so no need to create another version
	case err != nil:
		return err
	default:
		// if the files share the same checksum then make sure the file
		// was correctly copied over. If it wasn't then complete the copy.
		// If the file was completed then there is nothing to do, just delete
		// it.
		dirPath := MCFileDir(l.mcdir, df.ID)
		filePathInMCDir := filepath.Join(dirPath, df.ID)
		if checksum == df.Checksum {
			// We only need to check if things were correctly copied over if this
			// file doesn't point to a different file (ie, UsesID == ""). Otherwise
			// a file matching this was already successfully copied at some point in the past
			// and we are just pointing at that other file.
			if df.UsesID == "" {
				finfo2, err := os.Stat(filePathInMCDir)
				switch {
				case err != nil:
					// File doesn't exist so copy it over
					if err := l.moveFile(path, df); err != nil {
						return err
					}
				case finfo2.Size() != df.Size:
					// For some reason the file copy wasn't completed
					if err := l.moveFile(path, df); err != nil {
						return err
					}
				}
			}

			_ = os.Remove(path)
			return nil
		}
		addFile.Parent = df.ID
		if err := l.dfStore.UpdateDatafileCurrentFlag(df.ID, false); err != nil {
			return err
		}
	}

	var f model.DatafileSchema

	f, err = l.dfStore.AddDatafile(addFile)
	if err != nil {
		return err
	}

	return l.moveFile(path, f) // TODO: Delete file in database if moveFile fails...
}

func (l *MCFileLoader) computeFileChecksum(path string) (string, error) {
	hasher := md5.New()
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to load file %s to compute checksum", path)
	}

	defer f.Close()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func (l *MCFileLoader) moveFile(path string, f model.DatafileSchema) error {
	dirPath := MCFileDir(l.mcdir, f.ID)
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return err
	}

	if err := os.Rename(path, filepath.Join(dirPath, f.ID)); err != nil {
		// If rename fails try directly copying the file
		if err := copyFile(path, filepath.Join(dirPath, f.ID)); err != nil {
			return err
		}
		// File successfully copied so remove it
		_ = os.Remove(path)
	}
	return nil
}

func MCFileDir(path, fileID string) string {
	idSegments := strings.Split(fileID, "-")
	return filepath.Join(path, idSegments[1][0:2], idSegments[1][2:4])
}
