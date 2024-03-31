/*
Copyright © 2023 Eric Fisher epffisher@gmail.com

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
	"tmplhate/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile           string
	dontNormalizeVars bool
	printVersion      bool
	tmplLocation      string
	varsCase          string
	varsLocation      string
)

var Version = "development"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmplhate",
	Short: ".tmpl generator",
	Long:  `Tool to generate .tmpls from values.`,
	Run: func(cmd *cobra.Command, args []string) {
		h8 := new(core.Tmplhate)
		if printVersion {
			fmt.Printf("Version: %s\n", Version)
			os.Exit(0)
		}

		h8.Init(tmplLocation, varsLocation, dontNormalizeVars, varsCase)
		h8.WriteTemplate(os.Stdout)
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tmplhate.yaml)")
	rootCmd.PersistentFlags().StringVarP(&tmplLocation, "tmpl", "t", "", ".tmpl location")
	rootCmd.PersistentFlags().StringVarP(&varsLocation, "values", "l", "", "values location")
	rootCmd.PersistentFlags().BoolVar(&dontNormalizeVars, "dont-normalize", false, "don't normalize value cases (default is false)")
	rootCmd.PersistentFlags().StringVar(&varsCase, "case", "lower", "case used to reference values ['lower', 'upper', 'title'] (default is 'lower')")
	rootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "v", false, "print version info")
	rootCmd.MarkFlagRequired("values")
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

		// Search config in home directory with name ".tmplhate" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".tmplhate")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
