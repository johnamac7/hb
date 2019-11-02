/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func generateMarkdown() {
	err := doc.GenMarkdownTree(RootCmd, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate Markdown for the commands",
	Long:  `For hb generate Markdown Documents for each of the commands and write them to a folder named ./docs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Writing command descriptions to ./docs")
		generateMarkdown()
	},
}

func init() {
	RootCmd.AddCommand(docsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
