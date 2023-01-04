/*
Copyright © 2022 Dipak Parmar <hi@dipak.tech>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mta-sts-server",
	Short: "mta-sts server is a simple server to serve mta-sts.txt",
	Long:  `mta-sts server is a simple server to serve mta-sts.txt`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mta-sts-server.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".mta-sts-server" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mta-sts-server")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
func createConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Search config in home directory .config/mta-sts-server/config.yaml
	viper.AddConfigPath(home + "/.config/mta-sts-server")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	// instead of the default config use from environment variables/flags
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("domain", rootCmd.Flags().Lookup("domain"))
	viper.BindPFlag("mode", rootCmd.Flags().Lookup("mode"))
	viper.BindPFlag("max_age", rootCmd.Flags().Lookup("max_age"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))

	validateConfig()

	// write config
	err = viper.SafeWriteConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing config file:", viper.ConfigFileUsed())
	}
}

func validateConfig() {

	var errors []string

	// validate config
	if viper.GetString("domain") == "" {
		// append error to errors
		errors = append(errors, "domain is required")
	}

	if viper.GetString("mode") == "" {
		errors = append(errors, "mode is required")
	}

	if viper.GetInt("max_age") == 0 {
		errors = append(errors, "max_age is required")
	}

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, "Error in config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, "Please fix the following errors:")
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

}

// environment variables take precedence over config file values if set so override config file values with env vars
func overrideConfig() {
	// override config with env var
	viper.AutomaticEnv()

	// override config with flags
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("domain", rootCmd.Flags().Lookup("domain"))
	viper.BindPFlag("mode", rootCmd.Flags().Lookup("mode"))
	viper.BindPFlag("max_age", rootCmd.Flags().Lookup("max_age"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))

	// write config
	err := viper.WriteConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing config file:", viper.ConfigFileUsed())
	}

}
