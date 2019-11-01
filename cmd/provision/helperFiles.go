package provision

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

// helperFilesCmd represents the helperFiles command
var helperFilesCmd = &cobra.Command{
	Use:   "helper-files",
	Short: "Provision a set of Devices from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		config := NewConfig(cmd)
		filenames := FilesInDirectory(config.directory)
		provisionHelperFiles(config, filenames)
	},
}

func provisionHelperFiles(config Config, filenames []string) {
	for _, filename := range filenames {
		f, err := os.Open(config.directory + "/" + filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		client := resty.New().
			//SetDebug(true).
			SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

		resp, err := client.R().
			SetBasicAuth(config.username, config.password).
			SetFileReader("up_file", filename, f).
			Post("https://" + config.resource + "/api/v1/files/helper-files/" + filename + "/")

		if err != nil {
			fmt.Printf("Problem posting to Helper Files %v", err)
		}
		if resp.StatusCode() == 200 {
			fmt.Printf("Successfully uploaded %v %s", len(filenames), "Files \n")
		}
	}
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
