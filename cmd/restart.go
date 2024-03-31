// Copyright 2024 BINARY Members
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
	"context"
	"fmt"

	"github.com/B1NARY-GR0UP/openalysis/api"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart openalysis service",
	Long: `restart openalysis service 
e.g. oa restart -c "cron-spec" path2config.yaml`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := ""
		if !util.IsEmptySlice(args) {
			configPath = args[0]
		}
		if err := api.Init(configPath); err != nil {
			cobra.CheckErr(err)
		}
		if TokenF != "" {
			api.SetToken(TokenF)
		}
		if CronF != "" {
			api.SetCron(CronF)
		}
		if RetryF != -1 {
			api.SetRetry(RetryF)
		}
		fmt.Println(config.GlobalConfig)
		api.Restart(context.Background())
	},
}

func init() {
	setupCommand(restartCmd)
}
