package provision

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
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
	Short: "Provision a set of Device Groups from configuration files.",
	Long:  `The Device groups can be defined in YAML or JSON and conform to the payload definitions for the REST API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("debug") == "true" {
			resty.SetDebug(true)
		}
	},
	Run: func(c *cobra.Command, args []string) {
		config := cmd.NewConfig(c)
		config.Directory = c.Flag("directory").Value.String()
		config.Erase = c.Flag("erase").Value.String()
		filenames := FilesInDirectory(config.Directory)
		provisionDeviceGroups(config, filenames)
	},
}

func provisionDeviceGroups(config cmd.Config, filenames []string) {
	for _, filename := range filenames {
		var deviceGroups DeviceGroups
		if err := LoadConfiguration(config.Directory+"/"+filename, &deviceGroups); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}

		var resp *resty.Response
		var err error
		if config.Erase == "true" {
			resp, err = DELETE(deviceGroups, config.Resource, "/api/v1/device-groups/", config.Username, config.Password)
		} else {
			resp, err = POST(deviceGroups, config.Resource, "/api/v1/device-groups/", config.Username, config.Password)
		}
		if err != nil {
			fmt.Printf("Problem posting to DeviceGroups %v", err)
		}

		switch resp.StatusCode() {
		case 200:
			fmt.Printf("Successfully updated %v %s", len(deviceGroups.DeviceGroup), "Device Groups \n")
		default:
			fmt.Printf("Problem updating Device Groups: %v \n", resp.String())
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

	deviceGroupsCmd.PersistentFlags().BoolP("erase", "e", false, "to erase this configuration")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceGroupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
