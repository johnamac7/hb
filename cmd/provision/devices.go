package provision

import (
	"fmt"
	"log"

	"github.com/damianoneill/hb/cmd"
	"github.com/damianoneill/hb/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

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
		filenames := cmd.FilesInDirectory(config.Directory)
		provisionDevices(config, filenames)
	},
}

func deleteDevices(config cmd.Config, devices types.Devices) {
	noFailures := true
	for _, device := range devices.Device {
		resp, err := cmd.DELETE(config.Resource, "/api/v1/device/"+device.DeviceID+"/", config.Username, config.Password)
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

func createDevices(config cmd.Config, devices types.Devices) {
	resp, err := cmd.POST(devices, config.Resource, "/api/v1/devices/", config.Username, config.Password)
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
		var devices types.Devices
		if err := types.LoadConfiguration(config.Directory+"/"+filename, &devices); err != nil {
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
