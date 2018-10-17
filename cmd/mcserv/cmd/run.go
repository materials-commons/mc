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

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/file"
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
	Run: cliCmdRun,
}

var (
	port            int
	numberOfWorkers int
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port to listen on")
	runCmd.Flags().IntVarP(&numberOfWorkers, "workers", "w", 10, "Number of workers to use for processing uploads")
}

func cliCmdRun(cmd *cobra.Command, args []string) {
	db := connectToDB()
	e := setupEcho()
	setupAPIRoutes(e, db)
	mcdir := getMCDir()
	backgroundLoader := file.NewBackgroundLoader(mcdir, numberOfWorkers, db)
	backgroundLoader.Start()
	e.Start(fmt.Sprintf(":%d", port))
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

	//uc := &api.UsersController{}
	//g.POST("/getUserAPIKey", uc.GetUserByAPIKey).Name = "getUserByAPIKey"

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
