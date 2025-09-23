package cmd

import (
	"os"
	"strings"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/core/wordlistDB"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	viewArg string

	nameArg          string
	wordCountArg     string
	smallWordLenArg  string
	bigWordLenArg    string
	sizeArg          string
	avEntitiesArg    string
	avEntropyArg     string
	digitsArg        string
	upperArg         string
	spArg            string
	digitsUpperArg   string
	digitsSpArg      string
	upperSpArg       string
	digitsUpperSpArg string

	encodingArg string
	languageArg string
	categoryArg string
	authorArg   string
	tagsArg     string
	urlArg      string
)

var searchCmd = &cobra.Command{
	Use:   output.SearchHelpTexts["use"],
	Short: output.SearchHelpTexts["short"],
	Long:  output.SearchHelpTexts["long"],
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dbPath, err := utils.GetDatabasePath()
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		db, err := wordlistDB.NewWordlistRepository(dbPath)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		filter := wordlistDB.NewWordlistFilter(
			csvOrNil(nameArg),
			strOrNil(encodingArg),
			strOrNil(languageArg),
			strOrNil(authorArg),
			strOrNil(urlArg),
			strOrNil(categoryArg),
			uintRangeOrNil(wordCountArg),
			uintRangeOrNil(smallWordLenArg),
			uintRangeOrNil(bigWordLenArg),
			uintRangeOrNil(sizeArg),
			floatRangeOrNil(avEntitiesArg),
			floatRangeOrNil(avEntropyArg),
			floatRangeOrNil(digitsArg),
			floatRangeOrNil(upperArg),
			floatRangeOrNil(spArg),
			floatRangeOrNil(digitsUpperArg),
			floatRangeOrNil(digitsSpArg),
			floatRangeOrNil(upperSpArg),
			floatRangeOrNil(digitsUpperSpArg),
			csvOrNil(tagsArg),
		)

		list, err := db.Filter(filter)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		switch view := strings.ToLower(viewArg); view {
		case "summary":
			output.PrintWordlistsTableSummary(list)
		case "stats":
			output.PrintWordlistsTableStats(list)
		case "chars":
			output.PrintWordlistsTableChars(list)
		case "default":
			output.PrintWordlistsTableDefault(list)
		case "full":
			output.PrintWordlistsTableAll(list)
		default:
			output.PrintWarning("warnings", "invView", viewArg)
			output.PrintWordlistsTableDefault(list)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&viewArg, "view", "", output.SearchHelpTexts["view"])
	searchCmd.Flags().StringVarP(&nameArg, "name", "n", "", output.SearchHelpTexts["name"])
	searchCmd.Flags().StringVarP(&wordCountArg, "entities", "w", "", output.SearchHelpTexts["words"])
	searchCmd.Flags().StringVar(&smallWordLenArg, "smallest-entity", "", output.SearchHelpTexts["small_word"])
	searchCmd.Flags().StringVar(&bigWordLenArg, "biggest-entity", "", output.SearchHelpTexts["big_word"])
	searchCmd.Flags().StringVarP(&sizeArg, "size", "s", "", output.SearchHelpTexts["size"])
	searchCmd.Flags().StringVar(&avEntitiesArg, "avg-entity", "", output.SearchHelpTexts["avg-entity"])
	searchCmd.Flags().StringVar(&avEntropyArg, "avg-entropy", "", output.SearchHelpTexts["avg-entropy"])

	searchCmd.Flags().StringVarP(&digitsArg, "digits-perc", "d", "", output.SearchHelpTexts["digits-perc"])
	searchCmd.Flags().StringVarP(&upperArg, "upper-perc", "u", "", output.SearchHelpTexts["upper-perc"])
	searchCmd.Flags().StringVarP(&spArg, "special-perc", "S", "", output.SearchHelpTexts["special-perc"])
	searchCmd.Flags().StringVar(&digitsUpperArg, "digits-upper-perc", "", output.SearchHelpTexts["digits-upper-perc"])
	searchCmd.Flags().StringVar(&digitsSpArg, "digits-special-perc", "", output.SearchHelpTexts["digits-special-perc"])
	searchCmd.Flags().StringVar(&upperSpArg, "upper-special-perc", "", output.SearchHelpTexts["upper-special-perc"])
	searchCmd.Flags().StringVar(&digitsUpperSpArg, "all-perc", "", output.SearchHelpTexts["all-perc"])

	searchCmd.Flags().StringVar(&encodingArg, "encoding", "", output.SearchHelpTexts["encoding"])
	searchCmd.Flags().StringVarP(&languageArg, "language", "l", "", output.SearchHelpTexts["language"])
	searchCmd.Flags().StringVarP(&categoryArg, "category", "c", "", output.SearchHelpTexts["category"])
	searchCmd.Flags().StringVar(&authorArg, "author", "", output.SearchHelpTexts["author"])
	searchCmd.Flags().StringVarP(&tagsArg, "tags", "t", "", output.SearchHelpTexts["tags"])
	searchCmd.Flags().StringVar(&urlArg, "url", "", output.SearchHelpTexts["url"])
}

// strOrNil converting string variables into pointers
// Parameters:
//   - s: input string von cli Argument
//
// Returns:
//   - *[]string: pointer slice or nil if s empty
func strOrNil(s string) *[]string {
	if s == "" {
		return nil
	}
	sl := []string{s}
	return &sl
}

// csvOrNil method for converting string variables into pointers
// converts comma-separated input into a slice
// Parameters:
//   - s: input string von cli Argument
//
// Returns:
//   - *[]string: pointer slice or nil if s empty
func csvOrNil(s string) *[]string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return &out
}

// uintRangeOrNil converting string  to pointer
// Parameters:
//   - s: input string von cli Argument
//
// Returns:
//   - *utils.Uint64Range: pointer range Type or nil if s empty
func uintRangeOrNil(s string) *utils.Uint64Range {
	if s == "" {
		return nil
	}
	uRange, err := utils.NewUint64Range(s)
	if err != nil {
		output.PrintError("errors", "error", err)
		os.Exit(1)
	}
	return uRange
}

// floatRangeOrNil converting string  to pointer
// Parameters:
//   - s: input string von cli Argument
//
// Returns:
//   - *utils.Float64Range: pointer range Type or nil if s empty
func floatRangeOrNil(s string) *utils.Float64Range {
	if s == "" {
		return nil
	}
	fRange, err := utils.NewFloat64Range(s)
	if err != nil {
		output.PrintError("errors", "error", err)
		os.Exit(1)
	}
	return fRange
}
