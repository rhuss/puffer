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
	"errors"
	"log"
	"net"

	"github.com/rhuss/puffer/pkg/puffer"
	"github.com/rhuss/puffer/pkg/speak"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rhuss/dash"
	"path/filepath"
	"os"
	"encoding/json"
	"golang.org/x/oauth2"
	"fmt"
	"io/ioutil"
	"github.com/rhuss/puffer/pkg/calendar"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for the press of a Amazon Dash button",
	Long: `Watch the press of a Amazon Dash button

	`,
	Run: watch,
}

func watch(cmd *cobra.Command, args []string) {
	netInterface := viper.GetString("interface")
	if netInterface == "" {
		netInterface = "en3"
	}

	iface, err := net.InterfaceByName(netInterface)
	if err != nil {
		panic(err)
	}

	addr, err := extractAddress(iface)
	if err != nil {
		panic(err)
	}

	log.Printf("Using network range %v for interface %v", addr, iface.Name)
	pufferChan := dash.WatchButton(iface, ButtonMacAddress("puffer"))
	calendarChan := dash.WatchButton(iface, ButtonMacAddress("calendar"))
    for {
		select {
		case <- *pufferChan:
			pufferButtonPushed()
		case <- *calendarChan:
			calendarButtonPushed()
		}
	}
}
func calendarButtonPushed() {
	log.Print("Calendar Button pushed")

	jsonKey, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".puffer", "google-client-secret.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v", err)
		return
	}

	tokenCache := filepath.Join(os.Getenv("HOME"), ".puffer", "calendar-token.json")
	token, err := tokenFromFile(tokenCache)
	if err != nil {
		token, err = calendar.FetchToken(jsonKey)
		if err != nil {
			fmt.Printf("Cannot fetch token: %v", err)
			return
		}
		saveToken(tokenCache, token)
	}
	events, err := calendar.GetNextEvents(token, jsonKey, viper.GetStringSlice("calendars"), viper.GetStringSlice("allday"))
	if err != nil {
		fmt.Printf("Cannot fetch events: %v", err)
		return
	}
	if events.TodayEvents != nil {
		fmt.Println("Today")
		for _, event := range *events.TodayEvents {
			fmt.Printf("%s -- [%s,%s] : %s\n", event.Calendar, event.Start, event.End, event.Summary)
		}
	}
	if events.TomorrowEvents != nil {
		fmt.Println("Tomorrow")
		for _, event := range *events.TomorrowEvents {
			fmt.Printf("%s -- [%s,%s] : %s\n", event.Calendar, event.Start, event.End, event.Summary)
		}
	}
	if events.TomorrowAllDayEvents != nil {
		fmt.Println("All Day Events")
		for _, event := range *events.TomorrowAllDayEvents {
			fmt.Printf("%s -- %s\n", event.Calendar, event.Summary)
		}
	}
}


// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
  f, err := os.Open(file)
  if err != nil {
    return nil, err
  }
  t := &oauth2.Token{}
  err = json.NewDecoder(f).Decode(t)
  defer f.Close()
  return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
  fmt.Printf("Saving credential file to: %s\n", file)
  f, err := os.Create(file)
  if err != nil {
    log.Fatalf("Unable to cache oauth token: %v", err)
  }
  defer f.Close()
  json.NewEncoder(f).Encode(token)
}


func pufferButtonPushed() {
	log.Print("Puffer Button pushed")
	pufferData, err := puffer.FetchPufferData(PufferOptions())
	if err != nil {
		fmt.Print(err)
		return
	}
	log.Print("Puffer info fetched")
	speak.Speak(PufferMessage(pufferData), SpeakOptions())
}

func extractAddress(iface *net.Interface) (*net.IPNet, error) {
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return nil, err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask) - 4:],
					}
					break
				}
			}
		}
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return nil, errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return nil, errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return nil, errors.New("mask means network is too large")
	}
	return addr, nil
}

func init() {
	RootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
