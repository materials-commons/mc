package file_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/materials-commons/mc/internal/file"

	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestLoadFiles(t *testing.T) {
	var tmpFile string
	tmpDir, err := prepareTestDirTree("test/dir")
	assert.Okf(t, err, "Unable to create test dir %s", err)
	defer os.RemoveAll(tmpDir)

	tmpFile, err = createTmpFile(filepath.Join(tmpDir, "test"), "test file contents 1")
	assert.Okf(t, err, "Unable to create tmpfile %s", err)

	results := []string{""}

	loader := func(path string, finfo os.FileInfo) error {
		results = append(results, path)
		return nil
	}

	createSkipperForPath := func(whatToSkip string) file.Skipper {
		return func(path string, finfo os.FileInfo) bool {
			return path == whatToSkip
		}
	}

	skipEverthing := func(path string, finfo os.FileInfo) bool {
		return true
	}

	excludeListSkipper := file.NewExcludeListSkipper([]string{filepath.Join(tmpDir, "test/dir"), tmpFile})

	tests := []struct {
		skipper  file.Skipper
		expected []string
		name     string
	}{
		{
			skipper:  nil,
			expected: []string{filepath.Join(tmpDir, "test"), tmpFile, filepath.Join(tmpDir, "test/dir")},
			name:     "Collect all with default skipper",
		},
		{
			skipper:  createSkipperForPath(filepath.Join(tmpDir, "test/dir")),
			expected: []string{filepath.Join(tmpDir, "test"), tmpFile},
			name:     "Skip test/dir directory",
		},
		{
			skipper:  createSkipperForPath(tmpFile),
			expected: []string{filepath.Join(tmpDir, "test"), filepath.Join(tmpDir, "test/dir")},
			name:     fmt.Sprintf("Skip file %s", tmpFile),
		},
		{
			skipper:  excludeListSkipper.Skipper,
			expected: []string{filepath.Join(tmpDir, "test")},
			name:     fmt.Sprintf("ExcludeListSkipper only accept single dir imported %s", filepath.Join(tmpDir, "test")),
		},
		{
			skipper:  skipEverthing,
			expected: []string{},
			name:     "skip everything",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fileLoader := file.NewFileLoader(test.skipper, loader)
			results = nil
			err = fileLoader.LoadFiles(tmpDir)
			assert.Okf(t, err, "LoadFiles returned err %s", err)
			assert.Truef(t, compareArrays(test.expected, results), "expected != results %+v/%+v", test.expected, results)
		})
	}
}

func compareArrays(expected, what []string) bool {
	expectedMap := make(map[string]string)
	for _, entry := range expected {
		expectedMap[entry] = entry
	}

	for _, whatEntry := range what {
		if _, ok := expectedMap[whatEntry]; !ok {
			return false
		}
	}

	return true
}

func prepareTestDirTree(tree string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v\n", err)
	}

	err = os.MkdirAll(filepath.Join(tmpDir, tree), 0755)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tmpDir, nil
}

func createTmpFile(path, contents string) (string, error) {
	content := []byte(contents)
	tmpfile, err := ioutil.TempFile(path, "tmpfile.txt")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(content); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}
