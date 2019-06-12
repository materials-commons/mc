package file

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
)

type Zipper struct {
	w           *zip.Writer
	zipfilePath string
}

func CreateZipper(zipfilePath string) (*Zipper, error) {
	zipfile, err := os.Create(zipfilePath)
	if err != nil {
		return nil, err
	}

	return &Zipper{w: zip.NewWriter(zipfile), zipfilePath: zipfilePath}, nil
}

func (z *Zipper) AddToZipfile(filePath string, pathInZip string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file (%s/%s): %s", filePath, pathInZip, err)
	}

	zf, err := z.w.Create(pathInZip)
	if err != nil {
		return fmt.Errorf("unable to create file (%s) in zip: %s\n", pathInZip, err)
	}

	if _, err := zf.Write(data); err != nil {
		return fmt.Errorf("failed to write file (%s) to zip: %s\n", pathInZip, err)
	}

	return nil
}

func (z *Zipper) RemoveZip() {
	_ = z.w.Close()
	_ = os.Remove(z.zipfilePath)
}

func (z *Zipper) Close() error {
	return z.w.Close()
}
