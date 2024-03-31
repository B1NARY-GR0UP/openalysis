package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   CMD,
	Short: "OPEN ANALYSIS SERVICE",
	Long: `
 ██████╗ ██████╗ ███████╗███╗   ██╗ █████╗ ██╗  ██╗   ██╗███████╗██╗███████╗
██╔═══██╗██╔══██╗██╔════╝████╗  ██║██╔══██╗██║  ╚██╗ ██╔╝██╔════╝██║██╔════╝
██║   ██║██████╔╝█████╗  ██╔██╗ ██║███████║██║   ╚████╔╝ ███████╗██║███████╗
██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║██╔══██║██║    ╚██╔╝  ╚════██║██║╚════██║
╚██████╔╝██║     ███████╗██║ ╚████║██║  ██║███████╗██║   ███████║██║███████║
 ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝╚═╝   ╚══════╝╚═╝╚══════╝
`,
	Version: Version,
}

func init() {
	rootCmd.SetVersionTemplate("{{ .Version }}")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	TokenF string
	CronF  string
	RetryF int
)

var (
	defaultTokenF = ""
	defaultCronF  = ""
	defaultRetry  = -1
)

func setupCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)

	cmd.Flags().StringVarP(&TokenF, "token", "t", defaultTokenF, "your github token")
	cmd.Flags().StringVarP(&CronF, "cron", "c", defaultCronF, "your cron spec")
	cmd.Flags().IntVarP(&RetryF, "retry", "r", defaultRetry, "retry times")
}
