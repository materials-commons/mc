package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/materials-commons/mc/internal/ds"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/globusapi"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var globusDirectoryCmd = &cobra.Command{
	Use:   "globus-directory",
	Short: "Create the globus directory for the given dataset",
	Long:  "Create the globus directory for the given dataset",
	Run:   cliCmdGlobusDirectory,
}

var (
	// Globus
	globusEndpointID string
	globusCCUser     string
	globusCCToken    string
)

func init() {
	rootCmd.AddCommand(globusDirectoryCmd)
	globusDirectoryCmd.Flags().StringP("dataset-id", "d", "", "Dataset id to build zipfile for")
	globusDirectoryCmd.Flags().StringP("project-id", "p", "", "Project id dataset is in")
	globusDirectoryCmd.Flags().StringP("db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	globusDirectoryCmd.Flags().StringP("db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")
	globusDirectoryCmd.Flags().StringP("mcdir", "d", "/mcfs/data/materialscommons", "Database name to use (MCDB_NAME)")
	globusDirectoryCmd.Flags().BoolP("public", "t", false, "Is dataset public")
}

func cliCmdGlobusDirectory(cmd *cobra.Command, args []string) {
	var (
		datasetGlobusPath string
	)

	dbName, _ := cmd.Flags().GetString("db-name")
	dbConnect, _ := cmd.Flags().GetString("db-connect")
	datasetID, _ := cmd.Flags().GetString("dataset-id")
	projectID, _ := cmd.Flags().GetString("project-id")
	isPublicDataset, _ := cmd.Flags().GetBool("public")
	mcdir, _ := cmd.Flags().GetString("mcdir")

	setGlobusParams()

	globusClient, err := globusapi.CreateConfidentialClient(globusCCUser, globusCCToken)
	if err != nil {
		log.Fatalf("Unable to create globus client: %s", err)
	}

	session := connectToDB(dbName, dbConnect)

	db := store.NewDBRethinkdb(session)
	dataset, err := db.DatasetsStore().GetDataset(datasetID)
	if err != nil {
		fmt.Println("Unable to retrieve dataset:", err)
		os.Exit(1)
	}

	if isPublicDataset {
		datasetGlobusPath = filepath.Join(mcdir, "__published_datasets/%s", datasetID)
	} else {
		datasetGlobusPath = filepath.Join(mcdir, "__datasets/%s", datasetID)
	}

	dsDirLoader := ds.NewDirLoader(datasetGlobusPath, session)
	if err := dsDirLoader.LoadDirFromDataset(dataset, projectID); err != nil {
		// do something
	}

	if !isPublicDataset {
		projectUsers, err := db.ProjectsStore().GetProjectUsers(projectID)
		if err != nil {
			// do something
		}
		setACLAccessOnDataset(globusClient, datasetGlobusPath, projectUsers)
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

func setACLAccessOnDataset(client *globusapi.Client, datasetGlobusPath string, users []model.UserSchema) {
	for _, user := range users {
		if user.GlobusUser != "" {
			_ = setACL(client, datasetGlobusPath, user)
		}
	}
}

func setACL(client *globusapi.Client, path string, user model.UserSchema) error {
	identities, err := client.GetIdentities([]string{user.GlobusUser})
	if err != nil {
		return err
	}

	globusIdentityID := identities.Identities[0].ID

	rule := globusapi.EndpointACLRule{
		EndpointID:  globusEndpointID,
		Path:        path,
		IdentityID:  globusIdentityID,
		Permissions: "r",
	}

	aclRes, err := client.AddEndpointACLRule(rule)
	if err != nil {
		return err
	}

	_ = aclRes

	return nil
}
