package ds

import (
	"path/filepath"

	"github.com/materials-commons/mc/internal/store/model"
)

// Selection keeps track of directories and files and whether or not they have been included
type Selection struct {
	// IncludeFiles are files that have been explicitly included
	IncludeFiles map[string]bool

	// ExcludeFiles are files that have been explicitly excluded
	ExcludeFiles map[string]bool

	// IncludeFiles are directories that have been explicitly included
	IncludeDirs map[string]bool

	// ExcludeFiles are directories that have been explicitly excluded
	ExcludeDirs map[string]bool

	// Parents is used to track directories that are not in the included and excluded lists.
	// This is used, for example, to track descendant directories under a directory that has
	// been explicitly included or excluded. The application using selection calls AddDir to
	// add directories into this list.
	Parents map[string]bool
}

// NewSelection creates a new Selection. All map entries have been initialized.
func NewSelection() *Selection {
	return &Selection{
		IncludeFiles: make(map[string]bool),
		ExcludeFiles: make(map[string]bool),
		IncludeDirs:  make(map[string]bool),
		ExcludeDirs:  make(map[string]bool),
		Parents:      make(map[string]bool),
	}
}

// FromFileSelection creates a new Selection by taking a FileSelection and inserting
// its entries into the Selection fields.
func FromFileSelection(s *model.FileSelection) *Selection {
	selection := NewSelection()
	selection.loadIntoMap(selection.IncludeFiles, s.IncludeFiles)
	selection.loadIntoMap(selection.ExcludeFiles, s.ExcludeFiles)
	selection.loadIntoMap(selection.IncludeDirs, s.IncludeDirs)
	selection.loadIntoMap(selection.ExcludeDirs, s.ExcludeDirs)
	return selection
}

// loadIntoMap takes a Selection map and an array of string paths and loads the array
// entries into the map.
func (s *Selection) loadIntoMap(m map[string]bool, paths []string) {
	for _, p := range paths {
		m[p] = true
	}
}

// IsIncludedDir checks if a directory is a part of a selection. It first checks the IncludeDirs field,
// then checks the ExcludeDirs field. If no entry is found in either of these then it checks the Parents
// field and returns its value if it exists. If after all these checks there is no match then it returns
// false signifying that the directory is not one that has been selected.
func (s *Selection) IsIncludedDir(dirPath string) bool {
	if _, ok := s.IncludeDirs[dirPath]; ok {
		return true
	}

	if _, ok := s.ExcludeDirs[dirPath]; ok {
		return false
	}

	dirName := filepath.Dir(dirPath)
	for {
		if dirName == "." {
			break
		}

		if _, ok := s.IncludeDirs[dirName]; ok {
			return true
		}

		if _, ok := s.ExcludeDirs[dirName]; ok {
			return false
		}

		dirName = filepath.Dir(dirName)
	}

	if val, ok := s.Parents[dirPath]; ok {
		return val
	}

	val, ok := s.Parents[filepath.Dir(dirPath)]
	if !ok {
		return false
	}

	return val
}

// IsIncludedFile checks if a file is a part of a selection. It first checks IncludeFiles, then
// checks ExcludeFiles. If no entry has been found then it returns the files directory by calling
// IsIncludedDir on the file directory.
func (s *Selection) IsIncludedFile(filePath string) bool {
	if _, ok := s.IncludeFiles[filePath]; ok {
		return true
	}

	if _, ok := s.ExcludeFiles[filePath]; ok {
		return false
	}

	return s.IsIncludedDir(filepath.Dir(filePath))
}

// AddDir will add a directory to the Parents field. Callers should add directories to this list
// when that entry is not already in the selection. Directories are not added if they are in
// either IncludeDirs or ExcludeDirs. It checks ExcludeDirs when included is false, otherwise
// it checks IncludeDirs.
func (s *Selection) AddDir(dirPath string, included bool) {
	if included {
		if _, ok := s.IncludeDirs[dirPath]; ok {
			// Entry is already in IncludeDirs
			return
		}
	} else {
		// included is false so check ExcludeDirs to see if dirPath is already in there.
		if _, ok := s.ExcludeDirs[dirPath]; ok {
			// Entry is already in ExcludeDirs
			return
		}
	}

	s.Parents[dirPath] = included
}

// DirExists checks if the given dir path exists in the Selection. It returns
// 2 booleans. The first is true if dir is in IncludeDirs, and false if it is
// is ExcludeDirs. The second boolean denotes whether or not the dirPath was
// found. For example:
//   DirExists("path") => true, true --> in include dirs (first boolean) and found
//   DirExists("path") => false, true --> in exclude dirs (first boolean) and found
//   DirExists("path") => false, false --> first value is not meaningful, entry not found
func (s *Selection) DirExists(dirPath string) (isIncludedDir bool, exists bool) {
	if _, ok := s.IncludeDirs[dirPath]; ok {
		return true, true
	}

	if _, ok := s.ExcludeDirs[dirPath]; ok {
		return false, true
	}

	return false, false
}
