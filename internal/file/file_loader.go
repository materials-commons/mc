/**
 ** Loader will load all the files and directories at a root. It assumes that the root already exists so it skips
 ** sending the starting directory to the Skipper and Loader functions.
 */
package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"

	"github.com/pkg/errors"
)

type DirectoryLoader interface {
	LoadFileOrDir(path string, info os.FileInfo) error
}

type Skipper func(path string, finfo os.FileInfo) bool

type Loader struct {
	Skipper         Skipper         // Method to call to see if this entry should be skipped
	DirectoryLoader DirectoryLoader // Interface to call to load the entry if it is not skipped
}

// DefaultSkipper doesn't skip any entries.
func DefaultSkipper(path string, finfo os.FileInfo) bool {
	return false
}

// ExcludeListSkipper contains a list of entries to skip
type ExcludeListSkipper struct {
	ExcludeList map[string]string
}

// NewExcludeListSkipper creates a new ExludeListSkipper from the list of entries to exclude.
func NewExcludeListSkipper(excludeList []string) *ExcludeListSkipper {
	e := &ExcludeListSkipper{ExcludeList: make(map[string]string)}
	for _, entry := range excludeList {
		e.ExcludeList[entry] = ""
	}

	return e
}

func (s *ExcludeListSkipper) Skipper(path string, finfo os.FileInfo) bool {
	_, ok := s.ExcludeList[path]
	return ok
}

func NewFileLoader(skipper Skipper, loader DirectoryLoader) *Loader {
	s := DefaultSkipper

	if skipper != nil {
		s = skipper
	}

	return &Loader{Skipper: s, DirectoryLoader: loader}
}

func (l *Loader) LoadFiles(path string) error {
	err := filepath.Walk(path, func(fpath string, finfo os.FileInfo, err error) error {
		return l.processFile(path, fpath, finfo, err)
	})

	return err
}

func (l *Loader) LoadFilesWithCancel(path string, c context.Context) error {
	err := filepath.Walk(path, func(fpath string, finfo os.FileInfo, err error) error {
		select {
		case <-c.Done():
			log.Infof("LoadFilesWithCancel was cancelled")
			return errors.Errorf("canceled: %s", c.Err())
		default:
			return l.processFile(path, fpath, finfo, err)
		}
	})

	return err
}

func (l *Loader) processFile(root, fpath string, finfo os.FileInfo, err error) error {
	log.Infof("processFile %s, %s: %s", root, fpath, err)
	switch {
	case err != nil && os.IsPermission(err):
		// Permission errors are ignored. Just continue walking the tree
		// without processing the file or directory.
		return nil

	case err != nil:
		// All other errors cause walking to stop.
		return err

	case l.Skipper(fpath, finfo):
		// if Skipper returns true then skip processing this entry.
		if finfo.IsDir() {
			// If entry is a directory, then skip processing that
			// entire sub tree.
			return filepath.SkipDir
		}
		return nil

	case root == fpath:
		// Always skip processing the root. We assume this entry already exists or that loading is
		// starting from this entry. For example if you load /tmp/dir, then all entries under /tmp/dir
		// will be processed, but the /tmp/dir starting dir will be skipped.
		return nil

	case l.invalidFile(finfo):
		os.RemoveAll(fpath)
		return nil

	default:
		return l.DirectoryLoader.LoadFileOrDir(fpath, finfo)
	}
}

func (l *Loader) invalidFile(finfo os.FileInfo) bool {
	switch {
	case !finfo.IsDir() && finfo.Size() == 0:
		return true
	default:
		return l.isDotOrTempFile(finfo)
	}
}

func (l *Loader) isDotOrTempFile(fileInfo os.FileInfo) bool {
	name := fileInfo.Name()
	if strings.HasPrefix(name, ".") {
		return true
	}

	if strings.Contains(name, "~") {
		return true
	}

	return false
}
