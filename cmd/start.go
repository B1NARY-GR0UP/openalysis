package cmd

import (
	"context"
	"fmt"
	"github.com/B1NARY-GR0UP/openalysis/api"
	"github.com/B1NARY-GR0UP/openalysis/config"
	"github.com/B1NARY-GR0UP/openalysis/util"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start openalysis service",
	Long: `start openalysis service 
e.g. oa start -t "your-token" path2config.yaml`,
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
		api.Start(context.Background())
	},
}

func init() {
	setupCommand(startCmd)
}
