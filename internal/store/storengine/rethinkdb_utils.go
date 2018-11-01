package storengine

import (
	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v4"
)

func checkRethinkdbQueryError(res *r.Cursor, err error, msg string) error {
	switch {
	case err != nil:
		if res != nil {
			res.Close()
		}
		return errors.Wrapf(err, msg)
	case res.IsNil():
		res.Close()
		return errors.Wrapf(mc.ErrNotFound, msg)
	default:
		return nil
	}
}

func checkRethinkdbWriteError(resp r.WriteResponse, err error, msg string) error {
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
			return mc.ErrNotFound
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
			return mc.ErrNotFound
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
