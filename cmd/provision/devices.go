package provision

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Vendor - Configure the Vendor information
type Vendor struct {
	Juniper struct {
		OperatingSystem string `json:"operating-system" yaml:"operating-system"`
	} `json:"juniper"`
}

// Device - info need to Register a Device in Healthbot
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

// Device - info need to Register a Device in Healthbot

// Parse - tries to parse yaml first, then json into the Devices struct
func (c *Devices) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return err
		}
	}
	return nil
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Provision a set of Devices from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		directory := cmd.Flag("directory").Value.String()
		// can be overridden in ~/.hb.yaml so viper used
		resource := viper.GetString("resource")
		username := viper.GetString("username")
		password := viper.GetString("password")
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
		resp, err := POST(devices, resource, "/api/v1/devices/", username, password)
		if err != nil {
			fmt.Printf("Problem posting to Devices %v", err)
		}
		if resp.StatusCode() == 200 {
			fmt.Printf("Successfully provisioned %v %s", len(devices.Device), "Devices \n")
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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
