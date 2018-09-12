package store

import (
	"fmt"

	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v4"
)

func checkRethinkdbQueryError(res *r.Cursor, err error, msg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case res.IsNil():
		return errors.Wrapf(ErrNotFound, msg)
	default:
		return nil
	}
}

func checkRethinkdbWriteError(resp r.WriteResponse, err error, msg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case resp.Errors != 0:
		return fmt.Errorf("%s: %s", msg, resp.FirstError)
	default:
		return nil
	}
}
