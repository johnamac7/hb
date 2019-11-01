package provision

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
)

// Configuration - structures that get loaded from files
type Configuration interface {
	Parse(data []byte) error
	Dump(format string) string
}

// Config - bean for common provisioning configuration info
type Config struct {
	directory string
	resource  string
	username  string
	password  string
	erase     string
}

// NewConfig - construct the bean from viper / cmd
func NewConfig(cmd *cobra.Command) Config {
	return Config{
		directory: cmd.Flag("directory").Value.String(),
		resource:  viper.GetString("resource"),
		username:  viper.GetString("username"),
		password:  viper.GetString("password"),
	}
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

// POST - HTTP POST to a Resource
func POST(body interface{}, resource, path, username, password string) (resp *resty.Response, err error) {
	client := resty.New()
	//client.SetDebug(true)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err = client.R().
		SetBasicAuth(username, password).
		SetBody(body).
		Post("https://" + resource + path)

	return
}

// DELETE - HTTP POST to a Resource
func DELETE(body interface{}, resource, path, username, password string) (resp *resty.Response, err error) {
	client := resty.New()
	//client.SetDebug(true)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err = client.R().
		SetBasicAuth(username, password).
		SetBody(body).
		Delete("https://" + resource + path)

	return
}

// DumpYAMLOrJSON - For a configuration, output json or yaml
func DumpYAMLOrJSON(format string, configuration Configuration) string {
	var data []byte
	switch format {
	case "yaml":
		_ = yaml.Unmarshal(data, configuration)
	default:
		_ = json.Unmarshal(data, configuration)
	}
	return fmt.Sprintf("%v", data)
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
