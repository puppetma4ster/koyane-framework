package cmd

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/editor"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	sort        bool
	removeMask  bool
	removeRange string
)
var editCmd = &cobra.Command{
	Use:   output.GenerateEditHelpTexts["use"],
	Short: output.GenerateEditHelpTexts["short"],
	Long:  output.GenerateEditHelpTexts["long"],
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := args[1]

		wordlist, err := editor.NewEditWordlist(inputPath, outputPath)
		if err != nil {
			return
		}
		if sort {
			err := wordlist.SortWordlist()
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)

			}
		}
		if cmd.Flags().Changed("remove-range") {
			rRange, err := utils.NewUint64Range(removeRange)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
			err = wordlist.RemoveWordsByRangeUint(rRange)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		}
		err = wordlist.FlushFinishedWordlist()
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)

		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolVarP(&sort, "sort", "s", false, output.GenerateEditHelpTexts["sort"])
	editCmd.Flags().BoolVarP(&removeMask, "remove-mask", "m", false, output.GenerateEditHelpTexts["removeMask"])
	editCmd.Flags().StringVarP(&removeRange, "remove-range", "r", "", output.GenerateEditHelpTexts["removeMask"])

}
