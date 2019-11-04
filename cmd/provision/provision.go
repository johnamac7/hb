package provision

import (
	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision Healthbot Entities using config files.",
	Long:  `Grouping for a set of commands for provisioning Healthbot Entities.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

func init() {
	cmd.RootCmd.AddCommand(provisionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// provisionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// provisionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
