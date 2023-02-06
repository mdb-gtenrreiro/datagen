/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var templateFileName string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Command to trigger data generation",
	Long: `The create command takes a "--from" parameter that indicates where to find the 
	JSON data template.
	The JSON file template contains the shape of the data to be generated,
	as well as any "fake:{}" entries. For example: 
	
	{
		"id": "fake:{number:1,100}",
		"firstName": "fake:{firstname}",
		"lastName": "fake:{lastname}",
		"status": "active"
	}

	It is not required to use "fake:{}" commands. For example the 
	"status" key, has a value of "active"; that value will remain the same.
	The other values will be auto generated per the template specification.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		fmt.Println(templateFileName)
		fmt.Println(isKafka)
		fmt.Println(isFile)
		GenData(templateFileName, isKafka, isFile, topicName, limit)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().StringVarP(&templateFileName, "from", "f", "./templates/<filename>", "JSON Template file location.")
	createCmd.MarkFlagRequired("from")
}
