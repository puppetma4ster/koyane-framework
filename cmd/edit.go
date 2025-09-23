package cmd

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/editor"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	sortArg           bool
	removeMaskArg     string
	removeRangeArg    string
	deleteOriginalArg bool
)
var editCmd = &cobra.Command{
	Use:   output.GenerateEditHelpTexts["use"],
	Short: output.GenerateEditHelpTexts["short"],
	Long:  output.GenerateEditHelpTexts["long"],
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0]
		outputPath := args[1]

		wordlist, err := editor.NewEditWordlist(inputPath, outputPath, deleteOriginalArg)
		if err != nil {
			return
		}
		if sortArg {
			wordlist.ConcurrentSortWordlist()
		}
		if cmd.Flags().Changed("remove-range") {
			err = wordlist.ConcurrentRemoveWordsByRange(removeRangeArg)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		}
		if cmd.Flags().Changed("remove-mask") {
			err = wordlist.ConcurrentRemoveWordsWithMask(removeMaskArg)
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

	editCmd.Flags().BoolVarP(&sortArg, "sort", "s", false, output.GenerateEditHelpTexts["sort"])
	editCmd.Flags().StringVarP(&removeMaskArg, "remove-mask", "m", "", output.GenerateEditHelpTexts["removeMask"])
	editCmd.Flags().StringVarP(&removeRangeArg, "remove-range", "r", "", output.GenerateEditHelpTexts["removeMask"])
	editCmd.Flags().BoolVarP(&deleteOriginalArg, "delete", "d", false, output.GenerateEditHelpTexts["delete"])

}
