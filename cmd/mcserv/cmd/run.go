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
	Short: "Runs the mcserv file processor and API",
	Long:  `mcserv provides and API and background processor for loading files into Materials Commons that were uploaded into a given directory.`,
	Run:   cliCmdRun,
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
