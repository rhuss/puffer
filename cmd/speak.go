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
	"github.com/spf13/cobra"
)

// speakCmd represents the speak command
var speakCmd = &cobra.Command{
	Use:   "speak",
	Short: "Read the values of the puffer storage and speak it via audio",
	Long: `Get the puffer values and speak it out via audio.

	An external program is used to play the sound, which is:

	- afplay for OSX
	- mpg123 for Linux
	`,
	Run: func(cmd *cobra.Command, args []string) {
		PufferButtonPushed()
	},
}

func init() {
	RootCmd.AddCommand(speakCmd)
}
