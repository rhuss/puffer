// Copyright © 2016 Roland Huss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/rhuss/puffer/pkg/api"
	"github.com/rhuss/puffer/pkg/speak"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var gender string
var language string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "puffer",
	Short: "Managing data of a Sonnenkraft Puffer storage",
	Long: `puffer: Managing data of a Sonnenkraft Puffer storage
	
It can be used to query the current temperatur of the puffer storage.
	`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// SpeakOptions create the options for the text to speech service
func SpeakOptions() *speak.Options {
	ivonaConfig := viper.GetStringMapString("ivona")
	if ivonaConfig == nil {
		log.Fatal("No authentication for ivona configured")
	}
	access, found := ivonaConfig["access"]
	if !found {
		log.Fatal("No access for ivona found")
	}
	secret, found := ivonaConfig["secret"]
	if !found {
		log.Fatal("No secret given for accessing ivona")
	}
	return &speak.Options{
		Access:   access,
		Secret:   secret,
		Gender:   gender,
		Language: language,
	}
}

// PufferMessage returns the message to speak, depending on the language
func PufferMessage(info *api.PufferInfo) string {
	var format string
	if language == "de" {
		format = "Puffertemperatur. Oben : %d Grad Celsius. Mitte : %d Grad Celsius. Unten : %d Grad Celsius"
	} else {
		format = "Heat storage temperature. Up: %d degrees celsius. Middle : %d degrees celsius. Low: %d degrees Celsius"
	}
	return fmt.Sprintf(format, int(info.HighTemp+0.5), int(info.MidTemp+0.5), int(info.LowTemp+0.5))
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.puffer.yaml)")
	RootCmd.PersistentFlags().StringVarP(&gender, "gender", "g", "female", "Gender of voice to use (male or female)")
	RootCmd.PersistentFlags().StringVarP(&language, "language", "l", "de", "Language to use ('de' or 'en')")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".puffer") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}