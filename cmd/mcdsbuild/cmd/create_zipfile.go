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
	"os"
	"path/filepath"
	"time"

	"github.com/materials-commons/mc/internal/store"

	"github.com/materials-commons/mc/internal/file"

	"github.com/materials-commons/mc/internal/ds"

	"github.com/materials-commons/mc/internal/store/model"

	"github.com/apex/log"

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

// 65 Gigabytes is the max zip file size, this takes into account legacy zip file builds
const MaxZipfileSize int64 = 1024 * 1024 * 1024 * 65

var (
	zipfilePath string
)

func init() {
	rootCmd.AddCommand(createZipfileCmd)
	createZipfileCmd.Flags().StringP("dataset-id", "d", "", "Dataset id to build zipfile for")
	createZipfileCmd.Flags().StringP("project-id", "p", "", "Project id dataset is in")
	createZipfileCmd.Flags().StringP("db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	createZipfileCmd.Flags().StringP("db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")
	createZipfileCmd.Flags().StringP("zipfile", "z", "", "Full path to zipfile to create")
}

func cliCmdCreateZipfile(cmd *cobra.Command, args []string) {
	projectId, _ := cmd.Flags().GetString("project-id")
	datasetId, _ := cmd.Flags().GetString("dataset-id")
	dbName, _ := cmd.Flags().GetString("db-name")
	dbConnect, _ := cmd.Flags().GetString("db-connection")
	zipfilePath, _ = cmd.Flags().GetString("zipfile")

	if file.Exists(zipfilePath) {
		fmt.Println("Zipfile already exists")
		return
	}

	session := connectToDB(dbName, dbConnect)

	db := store.NewDBRethinkdb(session)
	dataset, err := db.DatasetsStore().GetDataset(datasetId)
	if err != nil {
		fmt.Println("Unable to retrieve dataset:", err)
		os.Exit(1)
	}

	selection := ds.FromFileSelection(&dataset.FileSelection)

	createDatasetZipfile(projectId, session, selection)
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

func createDatasetZipfile(projectId string, session *r.Session, selection *ds.Selection) {
	cursor, err := GetProjectDirsSortedCursor(projectId, session)
	if err != nil {
		log.Fatalf("Unable to retrieve project directories %s", err)
	}

	zipper, err := file.CreateZipper(zipfilePath)
	if err != nil {
		fmt.Printf("Unable to create zipfile %s: %s", zipfilePath, err)
	}
	defer zipper.Close()

	var dir model.DatadirSimpleModel
	for cursor.Next(&dir) {

		// Check if dir exists in selection, if not, then check its parent dir, and if that
		// exists set this dir to the parent dir setting. This reflects recursive selection as
		// parent directories that are included automatically include all descendants, and parent
		// directories that are excluded automatically exclude all descendants. These can be
		// overridden and selection will take that into account.
		if exists, _ := selection.DirExists(dir.Name); !exists {
			if exists, included := selection.DirExists(filepath.Dir(dir.Name)); exists {
				selection.AddDir(dir.Name, included)
			}
		}

		fileCursor, err := GetDirFilesCursor(dir.ID, session)
		if err != nil {
			continue
		}

		var totalSize int64 = 0

		var f model.DatafileSimpleModel
		for fileCursor.Next(&f) {
			if !f.Current {
				continue
			}

			totalSize += f.Size

			if totalSize > MaxZipfileSize {
				zipper.RemoveZip()
				fmt.Println("Zipfile too big, not building...")
				return
			}

			fullMCFilePath := filepath.Join(dir.Name, f.Name)
			if selection.IsIncludedFile(fullMCFilePath) {
				if err := zipper.AddToZipfile(f.FirstMCDirPath(), fullMCFilePath); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func GetProjectDirsSortedCursor(projectID string, session *r.Session) (*r.Cursor, error) {
	return r.Table("project2datadir").GetAllByIndex("project_id", projectID).
		EqJoin("datadir_id", r.Table("datadirs")).Zip().
		OrderBy("name").
		Run(session)
}

func GetDirFilesCursor(dirID string, session *r.Session) (*r.Cursor, error) {
	return r.Table("datadir2datafile").GetAllByIndex("datadir_id", dirID).
		EqJoin("datafile_id", r.Table("datafiles")).Zip().
		Run(session)
}
