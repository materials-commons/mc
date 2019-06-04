package ds

import (
	"path/filepath"

	"github.com/materials-commons/mc/internal/store/model"
)

// Selection keeps track of directories and files and whether or not they have been included
type Selection struct {
	// IncludedFiles are files that have been explicitly included
	IncludedFiles map[string]bool

	// ExcludedFiles are files that have been explicitly excluded
	ExcludedFiles map[string]bool

	// IncludedFiles are directories that have been explicitly included
	IncludedDirs map[string]bool

	// ExcludedFiles are directories that have been explicitly excluded
	ExcludedDirs map[string]bool

	// Parents is used to track directories that are not in the included and excluded lists.
	// This is used, for example, to track descendant directories under a directory that has
	// been explicitly included or excluded. The application using selection calls AddDir to
	// add directories into this list.
	Parents map[string]bool
}

// NewSelection creates a new Selection. All map entries have been initialized.
func NewSelection() *Selection {
	return &Selection{
		IncludedFiles: make(map[string]bool),
		ExcludedFiles: make(map[string]bool),
		IncludedDirs:  make(map[string]bool),
		ExcludedDirs:  make(map[string]bool),
		Parents:       make(map[string]bool),
	}
}

// FromFileSelection creates a new Selection by taking a FileSelection and inserting
// its entries into the Selection fields.
func FromFileSelection(s *model.FileSelection) *Selection {
	selection := NewSelection()
	selection.loadIntoMap(selection.IncludedFiles, s.IncludedFiles)
	selection.loadIntoMap(selection.ExcludedFiles, s.ExcludedFiles)
	selection.loadIntoMap(selection.IncludedDirs, s.IncludedDirs)
	selection.loadIntoMap(selection.ExcludedDirs, s.ExcludedDirs)
	return selection
}

// loadIntoMap takes a Selection map and an array of string paths and loads the array
// entries into the map.
func (s *Selection) loadIntoMap(m map[string]bool, paths []string) {
	for _, p := range paths {
		m[p] = true
	}
}

// IsIncludedDir checks if a directory is a part of a selection. It first checks the IncludedDirs field,
// then checks the ExcludedDirs field. If no entry is found in either of these then it checks the Parents
// field and returns its value if it exists. If after all these checks there is no match then it returns
// false signifying that the directory is not one that has been selected.
func (s *Selection) IsIncludedDir(dirPath string) bool {
	if _, ok := s.IncludedDirs[dirPath]; ok {
		return true
	}

	if _, ok := s.ExcludedDirs[dirPath]; ok {
		return false
	}

	val, ok := s.Parents[filepath.Dir(dirPath)]
	if !ok {
		return false
	}

	return val
}

// IsIncludedFile checks if a file is a part of a selection. It first checks IncludedFiles, then
// checks ExcludedFiles. If no entry has been found then it returns the files directory by calling
// IsIncludedDir on the file directory.
func (s *Selection) IsIncludedFile(filePath string) bool {
	if _, ok := s.IncludedFiles[filePath]; ok {
		return true
	}

	if _, ok := s.ExcludedFiles[filePath]; ok {
		return false
	}

	return s.IsIncludedDir(filepath.Dir(filePath))
}

// AddDir will add a directory to the Parents field. Callers should add directories to this list
// when that entry is not already in the selection. Directories are not added if they are in
// either IncludedDirs or ExcludedDirs. It checks ExcludedDirs when included is false, otherwise
// it checks IncludedDirs.
func (s *Selection) AddDir(dirPath string, included bool) {
	if included {
		if _, ok := s.IncludedDirs[dirPath]; ok {
			// Entry is already in IncludedDirs
			return
		}
	} else {
		// included is false so check ExcludedDirs to see if dirPath is already in there.
		if _, ok := s.ExcludedDirs[dirPath]; ok {
			// Entry is already in ExcludedDirs
			return
		}
	}

	s.Parents[dirPath] = included
}
