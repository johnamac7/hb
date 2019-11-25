/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
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

// playbooksCmd represents the playbook command
var playbooksCmd = &cobra.Command{
	Use:   "playbook",
	Short: "Provision Playbook from configuration files.",
	Long:  `The Playbook can be defined in YAML or JSON and conform to the payload definitions for the REST API.`,
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
		provisionPlaybooks(config, filenames)
	},
}

func deletePlaybooks(config cmd.Config, playbooks types.Playbooks) {
}

func createPlaybooks(config cmd.Config, playbooks types.Playbooks) {
	resp, err := cmd.POST(playbooks, config.Resource, "/api/v1/playbooks/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem posting to DeviceGroups %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		fmt.Printf("Successfully updated %v %s", len(playbooks.Playbooks), "Playbooks \n")
	default:
		fmt.Printf("Problem updating Device Groups: %v \n", resp.String())
		return
	}

	resp, err = cmd.POST(nil, config.Resource, "/api/v1/configuration/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem posting to Configuration %v", err)
		return
	}
	switch resp.StatusCode() {
	case 200:
		fmt.Printf("Successfully committed Playbooks configuration \n")
	default:
		fmt.Printf("Problem committing Playbooks configuration: %v \n", resp.String())
	}

}

func provisionPlaybooks(config cmd.Config, filenames []string) {
	for _, filename := range filenames {
		var playbooks types.Playbooks
		if err := types.LoadConfiguration(config.Directory+"/"+filename, &playbooks); err != nil {
			log.Fatal("Problem with "+filename+" ", err)
		}
		if config.Erase == "true" {
			deletePlaybooks(config, playbooks)
		} else {
			createPlaybooks(config, playbooks)
		}
	}
}

func init() {
	provisionCmd.AddCommand(playbooksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbookInstancesCmd.PersistentFlags().String("foo", "", "A help for foo")
	provisionCmd.PersistentFlags().StringP("directory", "d", "playbooks", "Default file location")

	provisionCmd.PersistentFlags().BoolP("erase", "e", false, "to erase this configuration")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbookInstancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
