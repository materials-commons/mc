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
	fmt.Printf("checkRethinkdbWriteError err = %s\n", err)
	fmt.Printf("checkRethinkdbWriteError resp = %+v\n", resp)
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case resp.Errors != 0:
		return errors.Errorf("%s: %s", msg, resp.FirstError)
	default:
		return nil
	}
}

func checkRethinkdbUpdateError(resp r.WriteResponse, err error, msg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, msg)
	case resp.Errors != 0:
		return errors.Errorf("%s: %s", msg, resp.FirstError)
	default:
		if resp.Updated == 0 && resp.Replaced == 0 {
			return errors.Errorf("%s: No documents found to change", msg)
		}
		return nil
	}
}

func checkRethinkdbDeleteError(resp r.WriteResponse, err error, errMsg string) error {
	switch {
	case err != nil:
		return errors.Wrapf(err, errMsg)
	case resp.Errors != 0:
		return errors.Errorf("%s: %s", errMsg, resp.FirstError)
	default:
		if resp.Deleted == 0 {
			return errors.Errorf("%s: No documents found to delete", errMsg)
		}
		return nil
	}
}

func checkRethinkdbInsertError(resp r.WriteResponse, err error, errMsg string) error {
	switch {
	case err != nil:
		return errors.Errorf("%s: %s", errMsg, err)
	case resp.Errors != 0:
		return errors.Errorf("%s: %s", errMsg, resp.FirstError)
	default:
		if len(resp.Changes) == 0 {
			return errors.Errorf("%s: nothing inserted", errMsg)
		}

		return nil
	}
}
