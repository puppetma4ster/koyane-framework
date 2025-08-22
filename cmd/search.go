package cmd

import (
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	nameArg      string
	fileSizeArg  string
	wordCountArg string
	wordLenArg   string

	encodingArg string
	languageArg string
	categoryArg string
	author      string
	tags        string
)

var searchCmd = &cobra.Command{
	Use:   output.SearchHelpTexts["use"],
	Short: output.SearchHelpTexts["short"],
	Long:  output.SearchHelpTexts["long"],
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&nameArg, "name", "n", "", output.SearchHelpTexts["name"])
	searchCmd.Flags().StringVarP(&fileSizeArg, "file-size", "f", "", output.SearchHelpTexts["file_size"])
	searchCmd.Flags().StringVarP(&wordCountArg, "words", "w", "", output.SearchHelpTexts["words"])
	searchCmd.Flags().StringVarP(&wordLenArg, "word-length", "l", "", output.SearchHelpTexts["word_length"])
	searchCmd.Flags().StringVarP(&encodingArg, "encoding", "e", "", output.SearchHelpTexts["encoding"])
	searchCmd.Flags().StringVarP(&languageArg, "language", "L", "", output.SearchHelpTexts["language"])
	searchCmd.Flags().StringVarP(&categoryArg, "category", "c", "", output.SearchHelpTexts["category"])
	searchCmd.Flags().StringVarP(&author, "author", "a", "", output.SearchHelpTexts["author"])
	searchCmd.Flags().StringVarP(&tags, "tags", "t", "", output.SearchHelpTexts["tags"])
}
