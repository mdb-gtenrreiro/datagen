/*
Copyright © 2023 Gustavo Tenrreiro gus.tenrreiro@mongodb.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var isFile bool
var isKafka bool
var topicName string
var limit uint64 = 0

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "datagen",
	Short: "Data generator utility",
	Long: `Datagen generates data from a predefined JSON template. The user must provide the template as a 
	parameter to the "create --from" command. By default the data will be outputed to STDOUT, however if the 
	"--filesystem" flag, or the "--kafka" flag is specified; or both; the data will be redirected there.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.datagen.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVar(&isFile, "filesystem", false, "If the data should be stored on the local file system, under ./data/")
	rootCmd.PersistentFlags().BoolVar(&isKafka, "kafka", false, "If the data should be sent to Kafka. Connection properties for kafka must be stored in ./conf/kafka.properties")
	rootCmd.PersistentFlags().StringVarP(&topicName, "topic", "t", "", "Kafka topic name to send data to")
	rootCmd.MarkFlagsRequiredTogether("kafka", "topic")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".datagen" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".datagen")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
