package file

type FileAndDirFilter interface {
	IsIncludedFile(filePath string) bool
	IsIncludedDir(dirPath string) bool
}
