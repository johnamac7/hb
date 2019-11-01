package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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

func generateMarkdown() {
	err := doc.GenMarkdownTree(RootCmd, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version
	//generateMarkdown()
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
