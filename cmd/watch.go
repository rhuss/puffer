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
	"bytes"
	"errors"
	"log"
	"time"

	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/rhuss/puffer/pkg/api"
	"github.com/rhuss/puffer/pkg/speak"
	"github.com/spf13/cobra"
)

var netInterface = "en3"
var buttonMac = "ac:63:be:fb:13:9d"

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for the press of a Amazon Dash button",
	Long: `Watch the press of a Amazon Dash button
	
	`,
	Run: watch,
}

func watch(cmd *cobra.Command, args []string) {
	iface, err := net.InterfaceByName(netInterface)
	if err != nil {
		panic(err)
	}

	addr, err := extractAddress(iface)
	if err != nil {
		panic(err)
	}

	log.Printf("Using network range %v for interface %v", addr, iface.Name)

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()

	// Start up a goroutine to read in packet data.
	stop := make(chan struct{})
	go watchForButton(handle, iface, stop, buttonPushed)
	defer close(stop)
	for {
		time.Sleep(10 * time.Second)
	}
}

func buttonPushed() {
	log.Printf("Button pushed !")
	pufferData := api.FetchPufferData()
	speak.Speak(PufferMessage(pufferData), SpeakOptions())
}

// watchForButton watches a handle for incoming ARP responses we might care about, and prints them.
//
// watchForButton loops until 'stop' is closed.
func watchForButton(handle *pcap.Handle, iface *net.Interface, stop chan struct{}, callback func()) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	var lastPushed = time.Time{}
	for {
		var packet gopacket.Packet
		select {
		case <-stop:
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation == layers.ARPRequest {
				if addressEquals("00:00:00:00:00:00", arp.DstHwAddress) && addressEquals(buttonMac, arp.SourceHwAddress) {
					var now = time.Now()
					if now.Sub(lastPushed).Seconds() > 5 {
						lastPushed = now
						go callback()
					}
				}
			}
		}
	}
}

func addressEquals(mac string, addr []byte) bool {
	macParsed, error := net.ParseMAC(mac)
	if error != nil {
		return false
	}
	return bytes.Equal(macParsed, addr)
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
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
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
