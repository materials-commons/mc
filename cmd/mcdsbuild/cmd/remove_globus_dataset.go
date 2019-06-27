package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/materials-commons/mc/pkg/globusapi"
	"github.com/spf13/cobra"
)

var removeGlobusDatasetCmd = &cobra.Command{
	Use:   "remove-globus-dataset",
	Short: "Removes the globus dataset",
	Long:  "Removes the globus dataset",
	Run:   cliCmdRemoveGlobusDataset,
}

func init() {
	rootCmd.AddCommand(removeGlobusDatasetCmd)
	removeGlobusDatasetCmd.Flags().StringP("dataset-id", "d", "", "Dataset id to remove")
	removeGlobusDatasetCmd.Flags().StringP("mcdir", "m", "/mcfs/data/materialscommons", "Database name to use (MCDB_NAME)")
	removeGlobusDatasetCmd.Flags().BoolP("private", "t", false, "Is dataset a private dataset")
}

func cliCmdRemoveGlobusDataset(cmd *cobra.Command, args []string) {
	var datasetGlobusPath string

	datasetID, _ := cmd.Flags().GetString("dataset-id")
	mcdir, _ := cmd.Flags().GetString("mcdir")
	isPrivateDataset, _ := cmd.Flags().GetBool("private")

	if mcdir == "" {
		fmt.Println("Invalid mcdir")
		os.Exit(1)
	}

	if datasetID == "" {
		fmt.Println("Invalid dataset-id")
		os.Exit(1)
	}

	// apiParams is initialized in root.go/init()
	globusClient, err := globusapi.CreateConfidentialClient(apiParams.GlobusCCUser, apiParams.GlobusCCToken)
	if err != nil {
		log.Fatalf("Unable to create globus client: %s", err)
	}

	removeAccessRules(globusClient, datasetID, isPrivateDataset)

	if isPrivateDataset {
		datasetGlobusPath = privateDatasetPath(mcdir, datasetID)
	} else {
		datasetGlobusPath = publicDatasetPath(mcdir, datasetID)
	}
	_ = os.RemoveAll(datasetGlobusPath)
}

func removeAccessRules(client *globusapi.Client, datasetID string, isPrivate bool) {
	var datasetGlobusPath string

	if isPrivate {
		datasetGlobusPath = privateGlobusDatasetPath(datasetID)
	} else {
		datasetGlobusPath = publicGlobusDatasetPath(datasetID)
	}

	endpointAccessRuleList, err := client.GetEndpointAccessRules(apiParams.GlobusEndpointID)
	if err != nil {
		// do something
		return
	}

	for _, rule := range endpointAccessRuleList.AccessRules {
		if rule.Path == datasetGlobusPath {
			// Remove this rule for this path, ignore errors as there is nothing we can do about them
			_, _ = client.DeleteEndpointACLRule(apiParams.GlobusEndpointID, rule.AccessID)
		}
	}
}
