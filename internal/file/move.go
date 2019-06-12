package file

import (
	"io"
	"os"
)

func MoveFile(src, dst string, removeSrc bool) error {
	if err := os.Rename(src, dst); err != nil {
		// If rename files try directly copying the file
		if err := copyFile(src, dst); err != nil {
			return err
		}

		if removeSrc {
			_ = os.Remove(src)
		}
	}

	return nil
}

// copyFile the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
