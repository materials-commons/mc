package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/materials-commons/mc/internal/ds"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/globusapi"
	"github.com/spf13/cobra"
)

var createGlobusDatasetCmd = &cobra.Command{
	Use:   "create-globus-dataset",
	Short: "Create dataset file structure in globus",
	Long:  "Create dataset file structure in globus",
	Run:   cliCmdCreateGlobusDataset,
}

func init() {
	rootCmd.AddCommand(createGlobusDatasetCmd)
	createGlobusDatasetCmd.Flags().StringP("dataset-id", "d", "", "Dataset id to build zipfile for")
	createGlobusDatasetCmd.Flags().StringP("project-id", "p", "", "Project id dataset is in")
	createGlobusDatasetCmd.Flags().StringP("db-connection", "c", "localhost:28015", "Database connection string (MCDB_CONNECTION)")
	createGlobusDatasetCmd.Flags().StringP("db-name", "n", "materialscommons", "Database name to use (MCDB_NAME)")
	createGlobusDatasetCmd.Flags().StringP("mcdir", "m", "/mcfs/data/materialscommons", "Database name to use (MCDB_NAME)")
	createGlobusDatasetCmd.Flags().BoolP("private", "t", false, "Is dataset a private dataset")
}

func cliCmdCreateGlobusDataset(cmd *cobra.Command, args []string) {
	var (
		datasetPath string
	)

	dbName, _ := cmd.Flags().GetString("db-name")
	dbConnect, _ := cmd.Flags().GetString("db-connection")
	datasetID, _ := cmd.Flags().GetString("dataset-id")
	projectID, _ := cmd.Flags().GetString("project-id")
	isPrivateDataset, _ := cmd.Flags().GetBool("private")
	mcdir, _ := cmd.Flags().GetString("mcdir")

	apiParams = globusapi.GetAPIParamsFromEnvFatal()

	// apiParams is initialized in root.go/init()
	globusClient, err := globusapi.CreateConfidentialClient(apiParams.GlobusCCUser, apiParams.GlobusCCToken)
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

	if isPrivateDataset {
		datasetPath = privateDatasetPath(mcdir, datasetID)
	} else {
		datasetPath = publicDatasetPath(mcdir, datasetID)
	}

	dsDirLoader := ds.NewDirLoader(datasetPath, session)
	if err := dsDirLoader.LoadDirFromDataset(dataset, projectID); err != nil {
		log.Fatalf("Unable to create data dir: %s", err)
	}

	if isPrivateDataset {
		projectUsers, err := db.ProjectsStore().GetProjectUsers(projectID)
		if err != nil {
			log.Fatalf("Unable to retrieve list of users: %s", err)
		}

		if err := setACLAccessOnPrivateDataset(globusClient, datasetID, projectUsers); err != nil {
			log.Fatalf("Error setting access rules: %s", err)
		}
	} else {
		if err := setACLAccessOnPublicDataset(globusClient, datasetID); err != nil {
			log.Fatalf("Error setting access rules: %s", err)
		}
	}
}

func setACLAccessOnPrivateDataset(client *globusapi.Client, datasetID string, users []model.UserSchema) error {
	path := privateGlobusDatasetPath(datasetID)
	for _, user := range users {
		if user.GlobusUser != "" {
			if err := setACL(client, path, user); err != nil {
				return err
			}
		}
	}

	return nil
}

func setACL(client *globusapi.Client, path string, user model.UserSchema) error {
	identities, err := client.GetIdentities([]string{user.GlobusUser})
	if err != nil {
		return err
	}

	globusIdentityID := identities.Identities[0].ID

	rule := globusapi.EndpointACLRule{
		PrincipalType: globusapi.ACLPrincipalTypeIdentity,
		EndpointID:    apiParams.GlobusEndpointID,
		Path:          path,
		IdentityID:    globusIdentityID,
		Permissions:   "r",
	}

	_, err = client.AddEndpointACLRule(rule)

	return err
}

func setACLAccessOnPublicDataset(client *globusapi.Client, datasetID string) error {
	path := publicGlobusDatasetPath(datasetID)
	rule := globusapi.EndpointACLRule{
		PrincipalType: globusapi.ACLPrincipalTypeAllAuthenticatedUsers,
		EndpointID:    apiParams.GlobusEndpointID,
		Path:          path,
		Permissions:   "r",
	}

	_, err := client.AddEndpointACLRule(rule)

	return err
}
