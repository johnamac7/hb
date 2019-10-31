package provision

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Devices - collection of Device
type Devices struct {
	Device []Device `json:"device"`
}

// Device - info need to Register a Device in Healthbot
type Device struct {
	Authentication struct {
		Password struct {
			Password string `json:"password"`
			Username string `json:"username"`
		} `json:"password"`
	} `json:"authentication"`
	DeviceID string `json:"device-id", yaml:"device-id"`
	Host     string `json:"host"`
}

// Parse - tries to parse yaml first, then json into the Devices struct
func (c *Devices) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	}
	return nil
}

func postDevices(devices Devices, resourceName, username, password string) {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetDebug(true)
	resp, err := client.R().
		SetBasicAuth(username, password).
		SetBody(devices).
		Post("https://" + resourceName + "/api/v1/devices/")

	fmt.Printf("%v, %v", resp, err)
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Provision a set of Devices from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		directory := cmd.Flag("directory").Value.String()
		resource := cmd.Flag("resource").Value.String()
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()
		filenames := FilesInDirectory(directory)
		provisionDevices(directory, filenames, resource, username, password)
	},
}

func provisionDevices(directory string, filenames []string, resource, username, password string) {
	for _, filename := range filenames {
		var devices Devices
		if err := LoadConfiguration(directory+"/"+filename, &devices); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}
		fmt.Printf("%+v\n", devices)
		postDevices(devices, resource, username, password)
	}
}

func init() {
	provisionCmd.AddCommand(devicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesCmd.PersistentFlags().String("foo", "", "A help for foo")
	devicesCmd.PersistentFlags().StringP("directory", "d", "devices", "Default file location")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
