package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart openalysis service",
	Long:  `restart openalysis service e.g. oa restart -c "cron-spec" path2config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restart called")
	},
}

func init() {
	setupCommand(restartCmd)
}
