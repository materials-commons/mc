// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"strings"
	"time"

	"github.com/materials-commons/mc/internal/file"

	"github.com/Jeffail/tunny"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/store"
	"github.com/spf13/cobra"
	r "gopkg.in/gorethink/gorethink.v4"
)

// siteCmd represents the site command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		db := connectToDB()
		e := setupEcho()
		setupAPIRoutes(e, db)
		mcdir := getMCDir()
		startBackgroundFileLoads(10, mcdir, db)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func connectToDB() store.DB {
	dbName := os.Getenv("MCDB")
	if dbName == "" {
		dbName = "materialscommons"
	}

	address := os.Getenv("MCDB_CONNECTION")
	if address == "" {
		address = "localhost:28015"
	}

	session, err := r.Connect(r.ConnectOpts{Database: dbName, Address: address})
	if err != nil {
		panic(fmt.Sprintf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", dbName, address, err))
	}

	r.SetTags("r")

	return store.NewDBRethinkdb(session)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func setupAPIRoutes(e *echo.Echo, db store.DB) {
	g := e.Group("/api")

	uc := &api.UsersController{}
	g.POST("/getUserAPIKey", uc.GetUserByAPIKey).Name = "getUserByAPIKey"

	fileLoaderController := api.NewFileLoaderController(db)
	g.POST("/loadFilesFromDirectory", fileLoaderController.LoadFilesFromDirectory).Name = "loadFilesFromDirectory"
}

func getMCDir() string {
	mcdir := strings.Split(os.Getenv("MCDIR"), ":")[0]
	if mcdir == "" {
		mcdir = "/mcfs/data/materialscommons"
	}
	return mcdir
}

func startBackgroundFileLoads(numberOfWorkers int, mcdir string, db store.DB) {
	pool := createPool(numberOfWorkers, mcdir, db)
	fileloadsStore := db.FileLoadsStore()

	// There may have been jobs in process when the server was stopped. Mark those jobs
	// at not currently being processed, this will cause them to be re-processed.
	if err := fileloadsStore.MarkAllNotLoading(); err != nil {
		panic(fmt.Sprintf("Unable to mark current jobs as not loading: %s", err))
	}

	// Loop through all file load requests and look for any that are not currently being processed.
	// The call to pool.Process(req) will block until the request is completed, so each processing
	// request is started in a new go routine. The pool will limit how many of these are currently
	// being processed.
	for {
		requests, err := fileloadsStore.GetAllFileLoads()
		if err != nil {
			break
		}

		for _, req := range requests {
			if !req.Loading {

				// Mark job as loading so it will be ignored in later processing
				err := fileloadsStore.UpdateLoading(req.ID, true)
				if err != nil {
					fmt.Printf("Unable to update file load request %s: %s", req.ID, err)

					// If the job cannot be marked as loading then skip processing it
					continue
				}

				// pool.Process() is synchronous, so run in separate routine and let the pool control
				// how many jobs are running simultaneously.
				go func() {
					pool.Process(req)
				}()
			}
		}

		// Sleep for 10 seconds before getting the next set of loading requests.
		time.Sleep(10 * time.Second)
	}
}

func createPool(numberOfWorkers int, mcdir string, db store.DB) *tunny.Pool {
	pool := tunny.NewFunc(numberOfWorkers, func(args interface{}) interface{} {
		dfStore := db.DatafilesStore()
		ddStore := db.DatadirsStore()
		projStore := db.ProjectsStore()
		flStore := db.FileLoadsStore()

		req := args.(store.FileLoadSchema)

		proj, err := projStore.GetProjectSimple(req.ProjectID)
		if err != nil {
			return err
		}

		loader := file.NewMCFileLoader(req.Path, req.Owner, mcdir, proj, dfStore, ddStore)
		skipper := file.NewExcludeListSkipper(req.Exclude)
		fl := file.NewFileLoader(skipper.Skipper, loader)

		if err := fl.LoadFiles(req.Path); err != nil {
			return err
		} else {
			// if loading files was successful then
			//    remove path since all files were processed and
			//    delete this file load request
			_ = os.RemoveAll(req.Path)
			return flStore.DeleteFileLoad(req.ID)
		}
	})

	return pool
}
