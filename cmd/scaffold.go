/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"
)

// scaffoldCmd represents the scaffold command
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate a config directory from an existing Healthbot installation",
	Long: `This command when pointed at an existing Healthbot installation, will generate
	valid configuration for the provision sub commands e.g. devices, device-groups, etc.
	
	The command requires a single argument, the directory where the configs should be written too, current directory is valid.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("debug") == "true" {
			resty.SetDebug(true)
		}
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("scaffold requires the name of the directory to store the config files")
		}
		return nil
	},
	Run: func(c *cobra.Command, args []string) {
		config := NewConfig(c)
		scaffold(config, args[0])
	},
}

func scaffold(config Config, path string) {
	fmt.Printf("Healthbot scaffold: %v\n", config.Resource)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	} else {
		cont := AskForConfirmation("scaffold directory "+path+" already exists, do you wish to continue?", 3, os.Stdout)
		if !cont {
			os.Exit(0)
		}
	}

	//

	os.Mkdir(path+string(filepath.Separator)+"devices", os.ModePerm)

	resp, err := GET(config.Resource, "/api/v1/devices/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem getting Devices %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		break
	default:
		fmt.Printf("Problem getting Devices: %v \n", resp.String())
		return
	}

	var devices Devices
	if err := json.Unmarshal(resp.Body(), &devices); err != nil {
		fmt.Printf("%v", err)
		return
	}

	// need to blank out the passwords
	blankPassword := "CHANGEME"
	for _, device := range devices.Device {
		if device.Authentication != nil && device.Authentication.Password.Password != nil {
			device.Authentication.Password.Password = &blankPassword
		}
	}

	data, err := yaml.Marshal(&devices)
	if err != nil {
		fmt.Printf("Problem with Marshalling Yaml: %v", err)
		return
	}
	err = ioutil.WriteFile(path+string(filepath.Separator)+"devices"+string(filepath.Separator)+"devices.yml", data, os.ModePerm)
	if err != nil {
		fmt.Printf("Problem writing Devices Config %v", err)
		return
	}

	//

	os.Mkdir(path+string(filepath.Separator)+"device-groups", os.ModePerm)

	resp, err = GET(config.Resource, "/api/v1/device-groups/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem getting Device Groups %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		break
	default:
		fmt.Printf("Problem getting Device Groups: %v \n", resp.String())
		return
	}

	var deviceGroups DeviceGroups
	if err := json.Unmarshal(resp.Body(), &deviceGroups); err != nil {
		fmt.Printf("%v", err)
		return
	}

	data, err = yaml.Marshal(&deviceGroups)
	if err != nil {
		fmt.Printf("Problem with Marshalling Yaml: %v", err)
		return
	}
	err = ioutil.WriteFile(path+string(filepath.Separator)+"device-groups"+string(filepath.Separator)+"device-groups.yml", data, os.ModePerm)
	if err != nil {
		fmt.Printf("Problem writing Device Groups Config %v", err)
		return
	}

}

func init() {
	RootCmd.AddCommand(scaffoldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scaffoldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scaffoldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
