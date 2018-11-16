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
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/materials-commons/mc/internal/controllers/uiapi"

	"github.com/apex/log"

	// "github.com/materials-commons/mc/internal/globus"
	"github.com/materials-commons/mc/pkg/globusapi"

	"github.com/materials-commons/mc/internal/store/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/file"
	"github.com/materials-commons/mc/internal/store"

	m "github.com/materials-commons/mc/internal/middleware"
	"github.com/spf13/cobra"

	r "gopkg.in/gorethink/gorethink.v4"
)

var (
	port            int
	numberOfWorkers int
	cfgFile         string

	// Materials Commons Object Store path
	mcdir string

	// Database
	dbConnection string
	dbName       string

	// Globus
	globusEndpointID string
	globusCCUser     string
	globusCCToken    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcserv",
	Short: "Runs the mcserv file processor and API",
	Long: `mcserv provides and API and background processor for loading files into Materials Commons 
that were uploaded into a given directory.`,
	Run: cliCmdRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcserv.yaml)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port to listen on")
	rootCmd.Flags().IntVarP(&numberOfWorkers, "workers", "w", 10, "Number of workers to use for processing uploads")
	rootCmd.Flags().StringVarP(&mcdir, "mcdir", "m", "/mcfs/data/materialscommons", "Locations of materials commons repo (MCDIR), can be colon separated list")
	rootCmd.Flags().StringVarP(&dbConnection, "db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	rootCmd.Flags().StringVarP(&dbName, "db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")

	setMCDir()
	setDBParams()
	setGlobusParams()
}

func cliCmdRoot(cmd *cobra.Command, args []string) {
	log.Infof("Starting mcserv...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := connectToDB()

	e := setupEcho()

	mcdirFirstEntry := strings.Split(mcdir, ":")[0]

	globusClient, err := globusapi.CreateConfidentialClient(globusCCUser, globusCCToken)
	if err != nil {
		log.Fatalf("Unable to create globus client: %s", err)
	}

	setupInternalAPIRoutes(e, db)
	setupAPIRoutes(e, db, mcdirFirstEntry, globusClient)

	backgroundLoader := file.NewBackgroundLoader(mcdirFirstEntry, numberOfWorkers, db)
	backgroundLoader.Start(ctx)

	// globusMonitor := globus.NewUploadMonitor(globusClient, globusEndpointID, db)
	// globusMonitor.Start(ctx)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
			log.Infof("Shutting down mcserv: %s", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Received signal, shutting down mcserv")
	cancel() // Have all the monitors start their shutdown process

	// Wait 2 seconds for everything to shutdown
	select {
	case <-time.After(2 * time.Second):
	}

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %s", err)
	}
}

func connectToDB() store.DB {
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

func setupInternalAPIRoutes(e *echo.Echo, db store.DB) {
	g := e.Group("/intapi")

	//uc := &api.UsersController{}
	//g.POST("/getUserAPIKey", uc.GetUserByAPIKey).Name = "getUserByAPIKey"

	fileLoaderController := api.NewFileLoaderController(db)
	g.POST("/loadFilesFromDirectory", fileLoaderController.LoadFilesFromDirectory).Name = "loadFilesFromDirectory"
	g.POST("/getFilesLoadRequest", fileLoaderController.GetFilesLoadRequest).Name = "getFilesLoadRequest"

	statusController := api.NewStatusController()
	g.POST("/getServerStatus", statusController.GetServerStatus).Name = "getServerStatus"
}

func setupAPIRoutes(e *echo.Echo, db store.DB, mcdir string, client *globusapi.Client) {
	apikey := createAPIKeyMiddleware(db)

	g := e.Group("/api")
	g.Use(apikey)

	globusController := api.NewGlobusController(db, client, mcdir, globusEndpointID)
	g.POST("/createGlobusUploadRequest", globusController.CreateGlobusUploadRequest).Name = "createGlobusUploadRequest"
	g.POST("/getGlobusUploadRequest", globusController.GetGlobusUploadRequest).Name = "getGlobusUploadRequest"
	g.POST("/listGlobusUploadRequests", globusController.ListGlobusUploadRequests).Name = "listGlobusUploadRequests"

	setupUIAPIRoutes(g, db)
}

func setupUIAPIRoutes(parent *echo.Group, db store.DB) {
	g := parent.Group("/ui")

	projectsController := uiapi.NewProjectsController(db)
	g.POST("/getProjectsForUser", projectsController.GetProjectsForUser).Name = "getProjectsForUser"
	g.POST("/getProjectOverview", projectsController.GetProjectOverview).Name = "getProjectOverview"
	g.POST("/getProjectAccessEntries", projectsController.GetProjectAccessEntries).Name = "getProjectAccessEntries"
}

func createAPIKeyMiddleware(db store.DB) echo.MiddlewareFunc {
	usersStore := db.UsersStore()

	apikeyConfig := m.APIKeyConfig{
		Skipper: middleware.DefaultSkipper,
		Keyname: "apikey",
		Retriever: func(apikey string, c echo.Context) (*model.UserSchema, error) {
			user, err := usersStore.GetUserByAPIKey(apikey)
			return &user, err
		},
	}

	return m.APIKeyAuth(apikeyConfig)
}

func setMCDir() {
	mcdirEnv := os.Getenv("MCDIR")
	if mcdirEnv != "" {
		mcdir = mcdirEnv
	}
}

func setDBParams() {
	mcdb := os.Getenv("MCDB_NAME")
	if mcdb != "" {
		dbName = mcdb
	}

	mcdbConnection := os.Getenv("MCDB_CONNECTION")
	if mcdbConnection != "" {
		dbConnection = mcdbConnection
	}
}

func setGlobusParams() {
	globusEndpointID = os.Getenv("MC_CONFIDENTIAL_CLIENT_ENDPOINT")
	globusCCUser = os.Getenv("MC_CONFIDENTIAL_CLIENT_USER")
	globusCCToken = os.Getenv("MC_CONFIDENTIAL_CLIENT_PW")

	if globusEndpointID == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_ENDPOINT env var is unset")
	}

	if globusCCUser == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_USER env var is unset")
	}

	if globusCCToken == "" {
		log.Fatalf("MC_CONFIDENTIAL_CLIENT_PW env var is unset")
	}
}
