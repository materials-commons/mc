/**
 ** FileLoader will load all the files and directories at a root. It assumes that the root already exists so it skips
 ** sending the starting directory to the Skipper and Loader functions.
 */
package file

import (
	"os"
	"path/filepath"
)

type Skipper func(path string, finfo os.FileInfo) bool
type Loader func(path string, info os.FileInfo) error

type FileLoader struct {
	Skipper Skipper // Method to call to see if this entry should be skipped
	Loader  Loader  // Method to call to load the entry if it is not skipped
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

func NewFileLoader(skipper Skipper, loader Loader) *FileLoader {
	s := DefaultSkipper

	if skipper != nil {
		s = skipper
	}

	if loader == nil {
		panic("A loader for NewFileLoader must be specified")
	}

	return &FileLoader{Skipper: s, Loader: loader}
}

func (l *FileLoader) LoadFiles(path string) error {
	err := filepath.Walk(path, func(fpath string, finfo os.FileInfo, err error) error {
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

		case path == fpath:
			// Always skip processing the root. We assume this entry already exists or that loading is
			// starting from this entry. For example if you load /tmp/dir, then all entries under /tmp/dir
			// will be processed, but the /tmp/dir starting dir will be skipped.
			return nil

		default:
			l.Loader(fpath, finfo)
			return nil
		}
	})

	return err
}
