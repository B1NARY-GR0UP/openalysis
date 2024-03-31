package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start openalysis service",
	Long:  `start openalysis service e.g. oa start -t "your-token" path2config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")
	},
}

func init() {
	setupCommand(startCmd)
}
