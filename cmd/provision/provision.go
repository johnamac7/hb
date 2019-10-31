package provision

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
)

// Configuration - structures that get loaded from files
type Configuration interface {
	Parse(data []byte) error
}

// LoadConfiguration - populates a configuration type with data from file
func LoadConfiguration(filelocation string, configuration Configuration) error {
	data, err := ioutil.ReadFile(filelocation)
	if err != nil {
		return err
	}
	if err := configuration.Parse(data); err != nil {
		return err
	}
	return nil
}

// FilesInDirectory - returns a list of filenames for a given directory
func FilesInDirectory(dirname string) (names []string) {
	fmt.Printf("Using directory: %s \n", dirname)
	f, err := os.Open(dirname)
	if err != nil {
		return
	}
	names, err = f.Readdirnames(-1)
	f.Close()
	fmt.Printf("Using files: %s \n", names)
	return
}

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision Healthbot Entities using config files",
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
