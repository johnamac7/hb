package provision

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// DeviceGroups - collection of Device Groups
type DeviceGroups struct {
	DeviceGroup []DeviceGroup `json:"device-group" yaml:"device-group"`
}

// DGAuthentication - Option to Override the individual Device Username/Passwords
type DGAuthentication struct {
	Password struct {
		Password *string `json:"password"`
		Username *string `json:"username"`
	} `json:"password,omitempty" yaml:"password,omitempty"`
}

// NativeGpb - Override the default JTI Port(s)
type NativeGpb struct {
	Ports []int `json:"ports"`
}

// DeviceGroup - info needed to Register a DeviceGroup in Healthbot
type DeviceGroup struct {
	DeviceGroupName string            `json:"device-group-name" yaml:"device-group-name"`
	Description     *string           `json:"description,omitempty" yaml:"description,omitempty"`
	Devices         *[]string         `json:"devices,omitempty" yaml:"devices,omitempty"`
	Playbooks       *[]string         `json:"playbooks,omitempty" yaml:"playbooks,omitempty"`
	Authentication  *DGAuthentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	NativeGpb       *NativeGpb        `json:"native-gpb,omitempty" yaml:"native-gpb,omitempty"`
}

// Parse - tries to parse yaml first, then json into the Devices struct
func (c *DeviceGroups) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	}
	return nil
}

// Dump - outputs DeviceGroups struct in either 'yaml' or 'json' format
func (c *DeviceGroups) Dump(format string) string {
	return DumpYAMLOrJSON(format, c)
}

// deviceGroupsCmd represents the deviceGroups command
var deviceGroupsCmd = &cobra.Command{
	Use:   "device-groups",
	Short: "Provision a set of Device Groups from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		directory := cmd.Flag("directory").Value.String()
		// can be overridden in ~/.hb.yaml so viper used
		resource := viper.GetString("resource")
		username := viper.GetString("username")
		password := viper.GetString("password")
		filenames := FilesInDirectory(directory)
		provisionDeviceGroups(directory, filenames, resource, username, password)
	},
}

func provisionDeviceGroups(directory string, filenames []string, resource, username, password string) {
	for _, filename := range filenames {
		var deviceGroups DeviceGroups
		if err := LoadConfiguration(directory+"/"+filename, &deviceGroups); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}
		resp, err := POST(deviceGroups, resource, "/api/v1/device-groups/", username, password)
		if err != nil {
			fmt.Printf("Problem posting to DeviceGroups %v", err)
		}
		if resp.StatusCode() == 200 {
			fmt.Printf("Successfully provisioned %v %s", len(deviceGroups.DeviceGroup), "Device Groups \n")
		}
	}
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
