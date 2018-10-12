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

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/store"
	"github.com/spf13/cobra"
	r "gopkg.in/gorethink/gorethink.v4"
)

// siteCmd represents the site command
var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := setupEcho()
		setupAPIRoutes(e)
	},
}

func init() {
	rootCmd.AddCommand(siteCmd)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func setupAPIRoutes(e *echo.Echo) {
	g := e.Group("/api")

	uc := &api.UsersController{}
	g.POST("/getUserAPIKey", uc.GetUserByAPIKey).Name = "getUserByAPIKey"

	db := os.Getenv("MCDB")
	if db == "" {
		db = "materialscommons"
	}

	address := os.Getenv("MCDB_CONNECTION")
	if address == "" {
		address = "localhost:28015"
	}

	session, err := r.Connect(r.ConnectOpts{Database: db, Address: address})
	if err != nil {
		panic(fmt.Sprintf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", db, address, err))
	}

	fileLoaderController := api.NewFileLoaderController(store.NewDBRethinkdb(session))
	g.POST("/loadFilesFromDirectory", fileLoaderController.LoadFilesFromDirectory).Name = "loadFilesFromDirectory"
}
