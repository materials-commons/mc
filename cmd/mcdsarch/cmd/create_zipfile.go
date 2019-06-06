// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/materials-commons/mc/internal/store"

	r "gopkg.in/gorethink/gorethink.v4"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createZipfileCmd = &cobra.Command{
	Use:   "create-zipfile",
	Short: "Create a zipfile for the given dataset",
	Long:  `Create a zipfile for the given dataset`,
	Run:   cliCmdCreateZipfile,
}

var (
	session *r.Session
)

func init() {
	rootCmd.AddCommand(createZipfileCmd)
	createZipfileCmd.Flags().StringP("dataset-id", "d", "", "Dataset id to build zipfile for")
	createZipfileCmd.Flags().StringP("project-id", "p", "", "Project id dataset is in")
	createZipfileCmd.Flags().StringP("db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	createZipfileCmd.Flags().StringP("db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")
}

func cliCmdCreateZipfile(cmd *cobra.Command, args []string) {
	fmt.Println("create called")

	projectId, _ := cmd.Flags().GetString("project-id")
	datasetId, _ := cmd.Flags().GetString("dataset-id")
	dbName, _ := cmd.Flags().GetString("db-name")
	dbConnect, _ := cmd.Flags().GetString("db-connection")

	session := connectToDB(dbName, dbConnect)

	createDatasetZipfile(projectId, datasetId, session)
}

func connectToDB(dbName, dbConnection string) *r.Session {
	opts := r.ConnectOpts{
		Database:   dbName,
		Address:    dbConnection,
		InitialCap: 10,
		MaxOpen:    20,
		Timeout:    1 * time.Second,
		NumRetries: 3,
	}
	session, err := r.Connect(opts)
	if err != nil {
		log.Fatalf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", dbName, dbConnection, err)
	}

	r.SetTags("r")

	return session
}

func createDatasetZipfile(projectId, datasetId string, session *r.Session) {
	dbStore := store.NewDBRethinkdb(session)
	GetProjectFilesCursor(session, dbStore)
}

func GetProjectFilesCursor(session *r.Session, db store.DB) {

}
