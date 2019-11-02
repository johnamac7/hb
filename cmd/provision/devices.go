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

// Devices - collection of Device
type Devices struct {
	Device []Device `json:"device"`
}

// Authentication - Collection type for Auth options
type Authentication struct {
	Password struct {
		Password *string `json:"password"`
		Username *string `json:"username"`
	} `json:"password,omitempty"`
}

// IAgent - configure the NETCONF port
type IAgent struct {
	Port int `json:"port"`
}

// OpenConfig - configure the Open Config port
type OpenConfig struct {
	Port int `json:"port"`
}

// V2 - configure the SNMP community string
type V2 struct {
	Community string `json:"community"`
}

// Snmp - configure the SNMP port or Community String
type Snmp struct {
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
	V2   *V2 `json:"v2,omitempty" yaml:"v2,omitempty"`
}

// Juniper - option to define the Operating system
type Juniper struct {
	OperatingSystem string `json:"operating-system" yaml:"operating-system"`
}

// Cisco - option to define the Operating system
type Cisco struct {
	OperatingSystem string `json:"operating-system" yaml:"operating-system"`
}

// Vendor - Configure the Vendor information
type Vendor struct {
	Juniper *Juniper `json:"juniper,omitempty"`
	Cisco   *Cisco   `json:"cisco,omitempty"`
}

// Device - info needed to Register a Device in Healthbot
type Device struct {
	DeviceID       string          `json:"device-id" yaml:"device-id"`
	Host           string          `json:"host"`
	SystemID       string          `json:"system-id,omitempty" yaml:"system-id,omitempty"`
	Authentication *Authentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	IAgent         *IAgent         `json:"iAgent,omitempty" yaml:"iAgent,omitempty"`
	OpenConfig     *OpenConfig     `json:"open-config,omitempty" yaml:"open-config,omitempty"`
	Snmp           *Snmp           `json:"snmp,omitempty" yaml:"snmp,omitempty"`
	Vendor         *Vendor         `json:"vendor,omitempty" yaml:"vendor,omitempty"`
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

// Dump - outputs Devices struct in either 'yaml' or 'json' format
func (c *Devices) Dump(format string) string {
	return DumpYAMLOrJSON(format, c)
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Provision a set of Devices from configuration files.",
	Long:  `The Devices can be defined in YAML or JSON and conform to the payload definitions for the REST API.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("debug") == "true" {
			resty.SetDebug(true)
		}
	},
	Run: func(c *cobra.Command, args []string) {
		config := cmd.NewConfig(c)
		config.Erase = c.Flag("erase").Value.String()
		config.Directory = c.Flag("directory").Value.String()
		filenames := FilesInDirectory(config.Directory)
		provisionDevices(config, filenames)
	},
}

func deleteDevices(config cmd.Config, devices Devices) {
	noFailures := true
	for _, device := range devices.Device {
		resp, err := DELETE(config.Resource, "/api/v1/device/"+device.DeviceID+"/", config.Username, config.Password)
		if err != nil {
			fmt.Printf("Problem posting to Devices %v", err)
			return
		}
		if resp.StatusCode() != 204 {
			fmt.Printf("Problem updating Device %v: %v \n", device.DeviceID, resp.String())
			noFailures = false
		}
	}
	if noFailures {
		fmt.Printf("Successfully updated %v %s", len(devices.Device), "Devices \n")
	}
}

func createDevices(config cmd.Config, devices Devices) {
	resp, err := POST(devices, config.Resource, "/api/v1/devices/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem posting to Devices %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		fmt.Printf("Successfully updated %v %s", len(devices.Device), "Devices \n")
	default:
		fmt.Printf("Problem updating Devices: %v \n", resp.String())
	}
}

func provisionDevices(config cmd.Config, filenames []string) {
	for _, filename := range filenames {
		var devices Devices
		if err := LoadConfiguration(config.Directory+"/"+filename, &devices); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}
		if config.Erase == "true" {
			deleteDevices(config, devices)
		} else {
			createDevices(config, devices)
		}
	}
}

func init() {
	provisionCmd.AddCommand(devicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesCmd.PersistentFlags().String("foo", "", "A help for foo")
	devicesCmd.PersistentFlags().StringP("directory", "d", "devices", "Default file location")

	devicesCmd.PersistentFlags().BoolP("erase", "e", false, "to erase this configuration")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
