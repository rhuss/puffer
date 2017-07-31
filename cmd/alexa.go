// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	alexa "github.com/mikeflynn/go-alexa/skillserver"
	"path/filepath"
	"os"
	"log"
)

// watchCmd represents the watch command
var alexaCmd = &cobra.Command{
	Use:   "alexa",
	Short: "Alexa Skill server",
	Long: `Provide Alexa skills for puffer information`,
	Run: alexaRun,
}

var config map[string]string


func alexaRun(cmd *cobra.Command, args []string) {
	config = viper.GetStringMapString("alexa")
	port, found := config["port"]
	if !found {
		port = "8443"
	}
	certPath := filepath.Join(os.Getenv("HOME"), ".puffer", "server.crt")
	keyPath := filepath.Join(os.Getenv("HOME"), ".puffer", "server.key")


	var applications = map[string]interface{}{
		"/echo/puffer": alexa.EchoApplication{ // Route
			AppID:    config["appid"],
			OnIntent: PufferHandler,
			OnLaunch: PufferHandler,
		},
	}

	log.Printf("Alexa Skillserver Listening on port %s", port)
	err := alexa.RunSSL(applications, ":" + port, certPath,keyPath)
	if err != nil {
		log.Fatal(err)
	}
}

func PufferHandler(echoReq *alexa.EchoRequest, echoResp *alexa.EchoResponse) {
	msg, err := getPufferSummaryMessage()
	if err != nil {
		log.Fatal(err)
	}
	echoResp.OutputSpeech(msg).Card("Puffer", msg)
}

func init() {
	RootCmd.AddCommand(alexaCmd)
}

