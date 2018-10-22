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
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/file"
	"github.com/materials-commons/mc/internal/store"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	r "gopkg.in/gorethink/gorethink.v4"
)

var (
	port            int
	numberOfWorkers int
	cfgFile         string
	mcdir           string
	dbConnection    string
	dbName          string
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcserv.yaml)")
	rootCmd.Flags().IntVarP(&port, "port", "p", 4000, "Port to listen on")
	rootCmd.Flags().IntVarP(&numberOfWorkers, "workers", "w", 10, "Number of workers to use for processing uploads")
	rootCmd.Flags().StringVarP(&mcdir, "mcdir", "m", "/mcfs/data/materialscommons", "Locations of materials commons repo (MCDIR), can be colon separated list")
	rootCmd.Flags().StringVarP(&dbConnection, "db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	rootCmd.Flags().StringVarP(&dbName, "db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")
	setMCDir()
	setDBParams()
}

func cliCmdRoot(cmd *cobra.Command, args []string) {
	db := connectToDB()
	e := setupEcho()
	setupInternalAPIRoutes(e, db)
	loaderDir := strings.Split(mcdir, ":")[0]
	backgroundLoader := file.NewBackgroundLoader(loaderDir, numberOfWorkers, db)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	backgroundLoader.Start(ctx)
	e.Start(fmt.Sprintf(":%d", port))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mcserv" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".mcserv")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func connectToDB() store.DB {
	session, err := r.Connect(r.ConnectOpts{Database: dbName, Address: dbConnection})
	if err != nil {
		panic(fmt.Sprintf("unable to connect to rethinkdb server, database: %s, address: %s, error: %s", dbName, dbConnection, err))
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
