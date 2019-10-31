package provision

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deviceGroupsCmd represents the deviceGroups command
var deviceGroupsCmd = &cobra.Command{
	Use:   "device-groups",
	Short: "Provision a set of Device Groups from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		provisionDeviceGroups(cmd.Flag("directory").Value.String())
	},
}

func provisionDeviceGroups(directory string) {
	fmt.Println("Using directory: " + directory)
}

func init() {
	provisionCmd.AddCommand(deviceGroupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceGroupsCmd.PersistentFlags().String("foo", "", "A help for foo")
	deviceGroupsCmd.PersistentFlags().StringP("directory", "d", "device-groups", "Default file location")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceGroupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
