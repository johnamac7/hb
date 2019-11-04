/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/damianoneill/hb/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Summarizes the Healthbot Installation.",
	Long:  `Provides some high level information on the installation version, Provisioned Devices and Device Groups.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if viper.GetString("debug") == "true" {
			resty.SetDebug(true)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		summary(NewConfig(cmd))
	},
}

// SystemDetails - Provides some basic hb info
type SystemDetails struct {
	ServerTime string `json:"server-time"`
	Version    string `json:"version"`
}

// DeviceFacts - Provides Device Facts
type DeviceFacts []struct {
	DeviceID string `json:"device-id"`
	Facts    struct {
		Hostname  string `json:"hostname"`
		JunosInfo []struct {
			LastRebootReason string `json:"last-reboot-reason"`
			MastershipState  string `json:"mastership-state"`
			Model            string `json:"model"`
			Name             string `json:"name"`
			Status           string `json:"status"`
			UpTime           string `json:"up-time"`
		} `json:"junos-info"`
		Platform     string `json:"platform"`
		PlatformInfo []struct {
			Name     string `json:"name"`
			Platform string `json:"platform"`
		} `json:"platform-info"`
		Product      string `json:"product"`
		Release      string `json:"release"`
		SerialNumber string `json:"serial-number"`
	} `json:"facts,omitempty"`
}

// NewTable - provides a blank table for rendering.
func NewTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetColumnSeparator("")
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.Append([]string{"", "", "", ""})
	return table
}

func summary(config Config) {
	resp, err := GET(config.Resource, "/api/v1/system-details/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem retrieving from Healthbot %v", err)
	}

	var systemDetails SystemDetails
	if err := json.Unmarshal(resp.Body(), &systemDetails); err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println("")
	fmt.Printf("Healthbot Resource: %s \n", config.Resource)
	fmt.Printf("Healthbot Version: %s \n", systemDetails.Version)
	fmt.Printf("Healthbot Time: %s \n", systemDetails.ServerTime)

	//

	resp, err = GET(config.Resource, "/api/v1/devices/facts/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem retrieving from Healthbot %v", err)
	}

	var deviceFacts DeviceFacts
	if err := json.Unmarshal(resp.Body(), &deviceFacts); err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println("")
	fmt.Printf("No of Managed Devices: %v \n", len(deviceFacts))

	fmt.Println("")

	table := NewTable()
	table.SetHeader([]string{"Device Id", "Platform", "Release", "Serial Number"})
	for _, fact := range deviceFacts {
		table.Append([]string{fact.DeviceID, fact.Facts.Platform, fact.Facts.Release, fact.Facts.SerialNumber})
	}
	table.Render() // Send output

	//

	resp, err = GET(config.Resource, "/api/v1/device-groups/", config.Username, config.Password)
	if err != nil {
		fmt.Printf("Problem retrieving from Healthbot %v", err)
	}

	var deviceGroups types.DeviceGroups
	if err := json.Unmarshal(resp.Body(), &deviceGroups); err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println("")
	fmt.Printf("No of Device Groups: %v \n", len(deviceGroups.DeviceGroup))

	fmt.Println("")

	table = NewTable()
	table.SetHeader([]string{"Device Group", "No of Devices"})
	for _, deviceGroup := range deviceGroups.DeviceGroup {
		table.Append([]string{deviceGroup.DeviceGroupName, strconv.Itoa(len(*deviceGroup.Devices))})
	}
	table.Render() // Send output

	fmt.Println("")
}

func init() {
	RootCmd.AddCommand(summaryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summaryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summaryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
