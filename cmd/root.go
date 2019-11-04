package cmd

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"

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
