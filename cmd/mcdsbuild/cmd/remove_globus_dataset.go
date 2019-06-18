package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
}

func cliCmdRemoveGlobusDataset(cmd *cobra.Command, args []string) {
	datasetID, _ := cmd.Flags().GetString("dataset-id")
	mcdir, _ := cmd.Flags().GetString("mcdir")

	if mcdir == "" {
		fmt.Println("Invalid mcdir")
		os.Exit(1)
	}

	if datasetID == "" {
		fmt.Println("Invalid dataset-id")
		os.Exit(1)
	}

	datasetGlobusPath := filepath.Join(mcdir, "__published_datasets", datasetID)
	_ = os.RemoveAll(datasetGlobusPath)
}
