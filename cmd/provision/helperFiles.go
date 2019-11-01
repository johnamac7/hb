package provision

import (
	"fmt"
	"os"

	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

// helperFilesCmd represents the helperFiles command
var helperFilesCmd = &cobra.Command{
	Use:   "helper-files",
	Short: "Provision a set of Devices from configuration files",
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("debug") == "true" {
			resty.SetDebug(true)
		}
	},
	Run: func(c *cobra.Command, args []string) {
		config := cmd.NewConfig(c)
		config.Directory = c.Flag("directory").Value.String()
		filenames := FilesInDirectory(config.Directory)
		provisionHelperFiles(config, filenames)
	},
}

func provisionHelperFiles(config cmd.Config, filenames []string) {
	for _, filename := range filenames {
		f, err := os.Open(config.Directory + "/" + filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		resp, err := resty.R().
			SetBasicAuth(config.Username, config.Password).
			SetFileReader("up_file", filename, f).
			Post("https://" + config.Resource + "/api/v1/files/helper-files/" + filename + "/")

		if err != nil {
			fmt.Printf("Problem posting to Helper Files %v", err)
		}

		switch resp.StatusCode() {
		case 200:
			fmt.Printf("Successfully uploaded %v %s", len(filenames), "Files \n")
		default:
			fmt.Printf("Problem uploading Files: %v \n", resp.String())
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
