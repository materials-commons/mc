package file

import "os"

// IsDir returns true if path is a directory.
func IsDir(path string) bool {
	finfo, err := os.Stat(path)
	switch {
	case err != nil:
		return false
	case finfo.IsDir():
		return true
	default:
		return false
	}
}

// Exists returns true if path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
