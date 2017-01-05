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

	"github.com/rhuss/puffer/pkg/puffer"
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

var Texts = map[string]map[string]string{
	"puffer": {
		"de": "Puffer. Oben : %d Grad. Mitte : %d Grad. Unten : %d Grad. Kollektor : %d Grad.",
		"en": "Heat storage. High: %d degrees celsius. Middle : %d degrees celsius. Low: %d degrees celsius. Collector : %d degrees celsius",
	},
	"cal-none" : {
		"de": "Heute keine Termine.",
		"en": "No events today",
	},
	"cal-timed-event" : {
		"de": "%s - %d Uhr : %s",
		"en": "%s - %d o'clock : %s",
	},
	"cal-timed-event-with-minute" : {
		"de": "%s - %d Uhr %d : %s",
		"en": "%s - %d %d : %s",
	},
	"cal-tomorrow" : {
		"de": "Termine morgen :",
		"en": "Events tomorrow :",
	},
	"cal-reminder-tomorrow" : {
		"de": "Erinnerung für morgen :",
		"en": "Reminder for tomorrow :",
	},
	"cal-event-no-time": {
		"de": "%s - %s",
		"en": "%s - %s",
	},
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

func PufferOptions() *puffer.Options {
	influxConfig := viper.GetStringMapString("influxdb")
	return &puffer.Options{
		Url:      influxConfig["url"],
		User:     influxConfig["user"],
		Password: influxConfig["password"],
	}
}

func ButtonMacAddress(what string) string {
	buttonConfig := viper.GetStringMapString("buttons")
	return buttonConfig[what]
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

	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.puffer") // adding home directory as first search path
	viper.AutomaticEnv()                 // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
