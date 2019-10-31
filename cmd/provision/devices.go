package provision

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
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
	DeviceID string `json:"device-id"`
	Host     string `json:"host"`
}

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Provision a set of Devices from configuration files",
	Run: func(cmd *cobra.Command, args []string) {
		provisionDevices(cmd.Flag("directory").Value.String())
	},
}

func provisionDevices(directory string) {
	fmt.Printf("Using directory: %s \n", directory)
	filenames, err := filesInDirectory(directory)
	if err != nil {
		fmt.Printf("Error Occurred: %s \n", err)
	} else {
		fmt.Printf("Using files: %s \n", filenames)
	}
	for _, filename := range filenames {
		jsonFile, err := os.Open(directory + "/" + filename)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Successfully Opened %s \n", filename)

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var devices Devices
		json.Unmarshal(byteValue, &devices)

		for _, device := range devices.Device {
			fmt.Printf("%+v\n", device)
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
