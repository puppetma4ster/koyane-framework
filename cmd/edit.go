package cmd

import (
	"fmt"
	"github.com/puppetma4ster/koyane-framework/internal/core/editor"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	sort bool
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
			wordlist, err = editor.SortWordlist(wordlist)
			if err != nil {
				fmt.Println("Fehler beim Sortieren:", err)
				return
			}
		}
		err = editor.FlushFinishedWordlist(wordlist)
		if err != nil {
			fmt.Println("Fehler beim Sortieren:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().BoolVarP(&sort, "sort", "s", false, output.GenerateEditHelpTexts["sort"])

}
