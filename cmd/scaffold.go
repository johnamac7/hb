package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/damianoneill/hb/types"
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
	valid configuration for the provision sub commands e.g. devices, device-groups, playbook-instances, etc.
	
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

func collectInfo(config Config, path, folder, resource, message string) *resty.Response {
	os.Mkdir(path+string(filepath.Separator)+folder, os.ModePerm)

	resp, err := GET(config.Resource, resource, config.Username, config.Password)
	if err != nil {
		fmt.Printf(message+" %v", err)
		os.Exit(1)
	}
	switch resp.StatusCode() {
	case 200:
		break
	default:
		fmt.Printf(message+": %v \n", resp.String())
		os.Exit(1)
	}
	return resp
}

func writeInfo(config interface{}, path, folder, filename string) {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Printf("Problem with Marshalling Yaml: %v", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(path+string(filepath.Separator)+folder+string(filepath.Separator)+filename, data, os.ModePerm)
	if err != nil {
		fmt.Printf("Problem writing "+folder+" config %v", err)
		os.Exit(1)
	}
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

	resp := collectInfo(config, path, "devices", "/api/v1/devices/", "Problem getting Devices")

	var devices types.Devices
	if err := json.Unmarshal(resp.Body(), &devices); err != nil {
		fmt.Printf("%v", err)
		return
	}

	// need to blank out the passwords
	blankPassword := "****"
	for _, device := range devices.Device {
		if device.Authentication != nil && device.Authentication.Password.Password != nil {
			device.Authentication.Password.Password = &blankPassword
		}
	}

	writeInfo(devices, path, "devices", "devices.yml")

	//

	dgResp := collectInfo(config, path, "device-groups", "/api/v1/device-groups/", "Problem getting Devices Groups")

	var deviceGroups types.DeviceGroups
	if err := json.Unmarshal(dgResp.Body(), &deviceGroups); err != nil {
		fmt.Printf("%v", err)
		return
	}

	writeInfo(deviceGroups, path, "device-groups", "device-groups.yml")

	//

	os.Mkdir(path+string(filepath.Separator)+"playbook-instances", os.ModePerm)

	var playbookInstances types.PlaybookInstances
	if err := json.Unmarshal(dgResp.Body(), &playbookInstances); err != nil {
		fmt.Printf("%v", err)
		return
	}

	writeInfo(playbookInstances, path, "playbook-instances", "playbook-instances.yml")

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
