/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ericdebeij/akamai-review/v2/internal/aksv"
)

const cfgDefaultFile = ".akamai-review.yaml"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "akamai-review",
	Short: "Review your account assets",
	Long: `akamai-review is a utility collection to extract information from
your akamai account and perform checks on it that need to be performed
on a regular base.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) { fmt.Println("hi") },
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

	akamaiConfig = &aksv.EdgeConfig{}

	viper.SetDefault("resolver", "8.8.8.8:53")
	viper.SetDefault("export", "export.csv")
	viper.SetDefault("log.level", "INFO")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgDefaultFile, "config file with all default parameters")
	rootCmd.PersistentFlags().StringVar(&akamaiConfig.Edgerc, "edgerc", "", "location of the credentials file")
	rootCmd.PersistentFlags().StringVar(&akamaiConfig.Section, "section", "", "section of the credentials file")
	rootCmd.PersistentFlags().StringVar(&akamaiConfig.AccountID, "account", "", "account switch key")
	cobra.OnInitialize(initConfig)
}

var akamaiConfig *aksv.EdgeConfig
var akamaiSession session.Session

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use config file from the flag.
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.SetLevelFromString(viper.GetString("log.level"))
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	} else if !errors.Is(err, os.ErrNotExist) && viper.ConfigFileUsed() != cfgDefaultFile {
		log.Infof("error reading config file: %s error %w", viper.ConfigFileUsed(), err)
	}

	openSession()
}
func openSession() {
	if akamaiConfig.Edgerc == "" {
		akamaiConfig.Edgerc = viper.GetString("akamai.edgerc")
	}

	if akamaiConfig.Section == "" {
		akamaiConfig.Section = viper.GetString("akamai.section")
	}

	if akamaiConfig.AccountID == "" {
		akamaiConfig.AccountID = viper.GetString("akamai.accountkey")
	}

	sess, err := aksv.NewSession(akamaiConfig)
	if err != nil {
		log.Errorf("session error %v", err)
		os.Exit(1)
	}

	akamaiSession = sess
}
