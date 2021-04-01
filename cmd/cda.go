package cmd

import (
	"github.com/linuxing3/cda/cda"
	"github.com/spf13/cobra"
)

// cdaCmd represents the cda command
var cdaCmd = &cobra.Command{
	Use:   "start",
	Short: "s",
	Long: `Start study in e-cda with default settings`,
	Run: func(cmd *cobra.Command, args []string) { 
		cda.StartCrawlCda()
	},
}

func init() {
	rootCmd.AddCommand(cdaCmd)
	cdaCmd.Flags().BoolP("default", "d", false, "Play cda")
}
