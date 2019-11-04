package cmd

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	// VERSION passed in as a build variable
	VERSION string
	cfgFile string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hb",
	Short: "Healthbot Command Line Interface",
	Long: `A tool for interacting with Healthbot over the REST API. 
	
The intent with this tool is to provide bulk or aggregate functions, that
simplify interacting with Healthbot. 
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Config - bean for common provisioning configuration info
type Config struct {
	Directory string
	Resource  string
	Username  string
	Password  string
	Erase     string
}

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

// Configuration - structures that get loaded from files
type Configuration interface {
	Parse(data []byte) error
	Dump(format string) string
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

// LoadConfiguration - populates a configuration type with data from file
func LoadConfiguration(filelocation string, configuration Configuration) error {
	data, err := ioutil.ReadFile(filelocation)
	if err != nil {
		return err
	}
	if err := configuration.Parse(data); err != nil {
		return err
	}
	return nil
}

// FilesInDirectory - returns a list of filenames for a given directory
func FilesInDirectory(dirname string) (names []string) {
	fmt.Printf("Using directory: %s \n", dirname)
	f, err := os.Open(dirname)
	if err != nil {
		return
	}
	names, err = f.Readdirnames(-1)
	f.Close()
	fmt.Printf("Using files: %s \n", names)
	return
}

// POST - HTTP POST to a Resource
func POST(body interface{}, resource, path, username, password string) (resp *resty.Response, err error) {
	resp, err = resty.R().
		SetBasicAuth(username, password).
		SetBody(body).
		Post("https://" + resource + path)
	return
}

// DELETE - HTTP POST to a Resource
func DELETE(resource, path, username, password string) (resp *resty.Response, err error) {
	resp, err = resty.R().
		SetBasicAuth(username, password).
		Delete("https://" + resource + path)

	return
}

// DumpYAMLOrJSON - For a configuration, output json or yaml
func DumpYAMLOrJSON(format string, configuration Configuration) string {
	var data []byte
	switch format {
	case "yaml":
		_ = yaml.Unmarshal(data, configuration)
	default:
		_ = json.Unmarshal(data, configuration)
	}
	return fmt.Sprintf("%v", data)
}

// Dump - outputs DeviceGroups struct in either 'yaml' or 'json' format
func (c *DeviceGroups) Dump(format string) string {
	return DumpYAMLOrJSON(format, c)
}

// AskForConfirmation - console y/n
func AskForConfirmation(s string, tries int, in io.Reader) bool {
	r := bufio.NewReader(in)
	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)
		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}
		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}
	return false
}

// NewConfig - construct the bean from viper / cmd
func NewConfig(cmd *cobra.Command) Config {
	return Config{
		Resource: viper.GetString("resource"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
	}
}

// GET - HTTP GET to a Resource
func GET(resource, path, username, password string) (resp *resty.Response, err error) {
	resp, err = resty.R().
		SetBasicAuth(username, password).
		Get("https://" + resource + path)
	return
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hb.yaml)")

	RootCmd.PersistentFlags().StringP("resource", "r", "localhost:8080", "Healthbot Resource Name")
	viper.BindPFlag("resource", RootCmd.PersistentFlags().Lookup("resource"))

	RootCmd.PersistentFlags().StringP("username", "u", "admin", "Healthbot Username")
	viper.BindPFlag("username", RootCmd.PersistentFlags().Lookup("username"))

	RootCmd.PersistentFlags().StringP("password", "p", "****", "Healthbot Password")
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))

	RootCmd.PersistentFlags().Bool("debug", false, "Enable REST debugging")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	//setup resty
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	viper.Set("restclient.RedirectPolicy", "always")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
