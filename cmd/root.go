/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/ericdebeij/akamai-review/v3/services"
	"github.com/ericdebeij/akamai-review/v3/util/logutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cfgDefaultFile = ".akamai-review.yaml"

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "akamai-review",
	Short: "Review your account assets",
	Long: `akamai-review is a utility collection to extract information from
your akamai account and perform checks on it that need to be performed
on a regular base.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("root command does not have a function of its own")
	//	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Cleanup() {
	logutil.CloseLogFile()
}

func param(cmd *cobra.Command, flag string, vip string, def interface{}, help string) {
	switch def := def.(type) {
	case string:
		cmd.PersistentFlags().String(flag, def, help)
		if def != "" {
			viper.SetDefault(vip, def)
		}
	case int:
		cmd.PersistentFlags().Int(flag, def, help)
		if def != 0 {
			viper.SetDefault(vip, def)
		}
	case bool:
		cmd.PersistentFlags().Bool(flag, def, help)
	case []string:
		cmd.PersistentFlags().StringArray(flag, def, help)
		viper.SetDefault(vip, def)
	default:
		log.Fatalf("type for default value not yet supported %s", def)
	}
	viper.BindPFlag(vip, cmd.PersistentFlags().Lookup(flag))
}
func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgDefaultFile, "config file with all default parameters")
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	for _, p := range services.Parameters {
		param(RootCmd, p.Flag, p.Viber, p.Default, p.Help)
	}

	cobra.OnInitialize(initConfig)
}

func viperAlias(base, key string) string {
	c := viper.GetString(base + "." + key)
	if c == "" {
		c = viper.GetString("default." + key)
	}
	if c == "" {
		return ""
	}

	alias := viper.GetStringMapString("alias." + key)
	v, f := alias[c]
	if f {
		return v
	}
	return c
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Use config file from the flag, searchpath does not work with . files
	if _, estat := os.Stat(cfgFile); estat != nil {
		hconfig := os.Getenv("HOME") + "/" + cfgFile

		if _, estat = os.Stat(hconfig); estat == nil {
			cfgFile = hconfig
		}
	}
	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix("AK_REVIEW")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()

	logutil.OpenLogFile(viper.GetString("log.file"))
	log.SetLevelFromString(viper.GetString("log.level"))
	if err == nil {
		log.Infof("using config file: %s", viper.ConfigFileUsed())
	} else {
		log.Infof("config file %v", err)
		if errors.Is(err, os.ErrNotExist) && viper.ConfigFileUsed() != cfgDefaultFile {
			fmt.Fprintf(os.Stderr, "config file %s not found\n", viper.ConfigFileUsed())
		}
	}

	services.StartServices()
}
