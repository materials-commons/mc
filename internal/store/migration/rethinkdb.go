package migration

import (
	"fmt"

	r "gopkg.in/gorethink/gorethink.v4"
)

func RethinkDB(database, address string) error {
	session, err := r.Connect(r.ConnectOpts{Database: database, Address: address})
	if err != nil {
		return fmt.Errorf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", database, address, err)
	}

	r.DBCreate(database).Exec(session)

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

		// templates table
		r.TableCreate("templates"),
		r.Table("templates").IndexCreate("owner"),
	}

	for _, table := range tables {
		if err := table.Exec(session); err != nil {
			//fmt.Printf("Failed creating table: %s", err)
		}
	}

	return nil
}
