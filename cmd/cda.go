package cmd

import (
	"github.com/linuxing3/cda/cda"
	"github.com/spf13/cobra"
)

// cdaCmd represents the cda command
var cdaCmd = &cobra.Command{
	Use:   "cda",
	Short: "c",
	Long: `Study in e-cda`,
	Run: func(cmd *cobra.Command, args []string) { 
		cda.StartCrawlCda()
	},
}

func init() {
	rootCmd.AddCommand(cdaCmd)
	cdaCmd.Flags().BoolP("default", "d", false, "Play cda")
}
