package migration

import (
	"fmt"

	"github.com/pkg/errors"

	r "gopkg.in/gorethink/gorethink.v4"
)

func RethinkDB(database, address string) error {
	session, err := r.Connect(r.ConnectOpts{Database: database, Address: address})
	if err != nil {
		return fmt.Errorf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", database, address, err)
	}

	_ = r.DBCreate(database).Exec(session)

	tables := []r.Term{
		// users table
		r.TableCreate("users"),
		r.Table("users").IndexCreate("apikey"),
		r.Table("users").IndexCreate("admin"),

		// projects table
		r.TableCreate("projects"),
		r.Table("projects").IndexCreate("name"),
		r.Table("projects").IndexCreate("owner"),
		r.Table("projects").IndexCreateFunc("name_owner", []interface{}{r.Row.Field("name"), r.Row.Field("owner")}),

		r.TableCreate("project2datadir"),
		r.Table("project2datadir").IndexCreate("datadir_id"),
		r.Table("project2datadir").IndexCreate("project_id"),
		r.Table("project2datadir").IndexCreateFunc("project_datadir", []interface{}{r.Row.Field("project_id"), r.Row.Field("datadir_id")}),

		r.TableCreate("project2datafile"),
		r.Table("project2datafile").IndexCreate("project_id"),
		r.Table("project2datafile").IndexCreate("datafile_id"),
		r.Table("project2datafile").IndexCreateFunc("project_datafile", []interface{}{r.Row.Field("project_id"), r.Row.Field("datafile_id")}),

		// access table
		r.TableCreate("access"),
		r.Table("access").IndexCreate("user_id"),
		r.Table("access").IndexCreate("project_id"),
		r.Table("access").IndexCreateFunc("user_project", []interface{}{r.Row.Field("user_id"), r.Row.Field("project_id")}),

		// samples table
		r.TableCreate("samples"),
		r.Table("samples").IndexCreate("project_id"),

		// experiments table

		// datafiles table
		r.TableCreate("datafiles"),
		r.Table("datafiles").IndexCreate("name"),
		r.Table("datafiles").IndexCreate("owner"),
		r.Table("datafiles").IndexCreate("checksum"),
		r.Table("datafiles").IndexCreate("usesid"),

		// datadirs table
		r.TableCreate("datadirs"),
		r.Table("datadirs").IndexCreate("name"),
		r.Table("datadirs").IndexCreate("project"),
		r.Table("datadirs").IndexCreate("parent"),
		r.Table("datadirs").IndexCreateFunc("datadir_project_name", []interface{}{r.Row.Field("project"), r.Row.Field("name")}),
		r.Table("datadirs").IndexCreateFunc("datadir_project_shortcut", []interface{}{r.Row.Field("project"), r.Row.Field("shortcut")}),

		r.TableCreate("datadir2datafile"),
		r.Table("datadir2datafile").IndexCreate("datadir_id"),
		r.Table("datadir2datafile").IndexCreate("datafile_id"),
		r.Table("datadir2datafile").IndexCreateFunc("datadir_datafile", []interface{}{r.Row.Field("datadir_id"), r.Row.Field("datafile_id")}),

		// templates table
		r.TableCreate("templates"),
		r.Table("templates").IndexCreate("owner"),

		// file loads table
		r.TableCreate("file_loads"),
	}

	var errOnExec error
	for _, table := range tables {
		if err := table.Exec(session); err != nil {
			if errOnExec != nil {
				errOnExec = errors.Errorf("Error processing a table: %s", err)
			}
		}
	}

	return errOnExec
}
