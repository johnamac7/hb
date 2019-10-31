package provision

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helperFilesCmd represents the helperFiles command
var helperFilesCmd = &cobra.Command{
	Use:   "helper-files",
	Short: "Provision a set of Devices from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		provisionHelperFiles(cmd.Flag("directory").Value.String())
	},
}

func provisionHelperFiles(directory string) {
	fmt.Println("Using directory: " + directory)
}

func init() {
	provisionCmd.AddCommand(helperFilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	helperFilesCmd.PersistentFlags().StringP("directory", "d", "helper-files", "Default file location")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helperFilesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
