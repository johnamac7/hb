package provision

import (
	"fmt"
	"log"

	"github.com/damianoneill/hb/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

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
		filenames := cmd.FilesInDirectory(config.Directory)
		provisionDeviceGroups(config, filenames)
	},
}

func deleteDeviceGroups(config cmd.Config, deviceGroups cmd.DeviceGroups) {
	noFailures := true
	for _, dg := range deviceGroups.DeviceGroup {
		resp, err := cmd.DELETE(config.Resource, "/api/v1/device-group/"+dg.DeviceGroupName+"/", config.Username, config.Password)
		if err != nil {
			fmt.Printf("Problem posting to DeviceGroups %v", err)
			return
		}
		if resp.StatusCode() != 204 {
			fmt.Printf("Problem updating Device Group %v: %v \n", dg.DeviceGroupName, resp.String())
			noFailures = false
		}
	}
	if noFailures {
		fmt.Printf("Successfully updated %v %s", len(deviceGroups.DeviceGroup), "Device Groups \n")
	}
}

func createDeviceGroups(config cmd.Config, deviceGroups cmd.DeviceGroups) {
	resp, err := cmd.POST(deviceGroups, config.Resource, "/api/v1/device-groups/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem posting to DeviceGroups %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		fmt.Printf("Successfully updated %v %s", len(deviceGroups.DeviceGroup), "Device Groups \n")
	default:
		fmt.Printf("Problem updating Device Groups: %v \n", resp.String())
	}
}

func provisionDeviceGroups(config cmd.Config, filenames []string) {
	for _, filename := range filenames {
		var deviceGroups cmd.DeviceGroups
		if err := cmd.LoadConfiguration(config.Directory+"/"+filename, &deviceGroups); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}
		if config.Erase == "true" {
			deleteDeviceGroups(config, deviceGroups)
		} else {
			createDeviceGroups(config, deviceGroups)
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
