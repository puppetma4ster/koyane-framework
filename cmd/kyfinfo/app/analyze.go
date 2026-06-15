package app

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/analyzer"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	allArg         bool
	generalArg     bool
	contentArg     bool
	statsArg       bool
	saveFileArg    string
	withDuplicates bool
	quiet          bool
)

var analyzeCmd = &cobra.Command{
	Use:   output.AnalyzeHelpTexts["use"],
	Short: output.AnalyzeHelpTexts["short"],
	Long:  output.AnalyzeHelpTexts["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputPath := args[0] // input word list path
		var printer = output.NewAnalyzePrinter(analyzer.NewGeneralDummy(), analyzer.NewContentDummy())
		stop := make(chan struct{})

		// print banner
		if !quiet {
			output.PrintBanner()
		}
		err := utils.CreateTempDir() // generate tmp dir if no tmp dir is found
		if err != nil {              // creates temp folder to /tmp/koyane_framework_tmp
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		// starting spinner animation...
		go output.Spinner("Analyze File", stop)

		if allArg || generalArg && contentArg && statsArg { // everything is to be printed
			general, err := analyzer.NewGeneralAnalyzer(inputPath)
			if err != nil { // retrieve general data
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
			isSaved, err := analyzer.IsSavedObj(general)
			var content *analyzer.AnalyzerContent
			if isSaved { // checking if wordlist is not changed and already analyzed
				general, content, err = analyzer.LoadFromYamlAll(*general) // loading wordlist analytics
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			} else { // when file is not analyzed yet
				content, err = analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true, true, true, true, true)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			}

			// printing gathered information...
			printer = output.NewAnalyzePrinter(general, content)
			printer.PrintAllGeneralInfo()
			printer.PrintAllContentInfo()
			printer.PrintAllStatsInfo(withDuplicates)

			err = analyzer.SaveToYamlAll(general, content) // saving the values from the wordlist
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}

			// deleting buffers
			printer.FlushGeneral()
			printer.FlushContent()
			printer.FlushStats()

		} else if statsArg && contentArg { // Only stats and content were accessed

			general := analyzer.NewGeneralDummy() // dummy object because “general” doesn't need to be printed
			content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true,
				true, true, true, true)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}

			// printing gathered information...
			printer = output.NewAnalyzePrinter(general, content)

			printer.PrintAllContentInfo()
			printer.PrintAllStatsInfo(withDuplicates)

			printer.FlushContent()
			printer.FlushStats()
		} else {
			if generalArg { // Extract only general information
				content := analyzer.NewContentDummy()
				general, err := analyzer.NewGeneralAnalyzer(inputPath)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}

				printer = output.NewAnalyzePrinter(general, content)

				printer.PrintAllGeneralInfo()

				printer.FlushGeneral()
			}
			if contentArg { // Extract only content information
				general := analyzer.NewGeneralDummy()
				content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, true, true,
					false, true, true, false)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				printer = output.NewAnalyzePrinter(general, content)

				printer.PrintAllContentInfo()

				printer.FlushContent()

			}
			if statsArg { // Extract only statistic information
				general := analyzer.NewGeneralDummy()
				content, err := analyzer.ConcurrentContentAnalyzer(inputPath, true, false, false,
					true, false, true, true)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				printer = output.NewAnalyzePrinter(general, content)
				printer.PrintAllStatsInfo(withDuplicates)
				printer.FlushStats()
			}
		}
		close(stop)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := analyzeCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// output all information
	analyzeCmd.Flags().BoolVarP(&allArg, "all", "a", false, output.AnalyzeHelpTexts["all"])
	// output only general information
	analyzeCmd.Flags().BoolVarP(&generalArg, "general", "g", false, output.AnalyzeHelpTexts["general"])
	// output only content information
	analyzeCmd.Flags().BoolVarP(&contentArg, "content", "c", false, output.AnalyzeHelpTexts["content"])
	// output only stat information
	analyzeCmd.Flags().BoolVarP(&statsArg, "stats", "s", false, output.AnalyzeHelpTexts["stats"])

	// print duplicate  words
	analyzeCmd.Flags().BoolVarP(&withDuplicates, "duplicates", "d", false, output.AnalyzeHelpTexts["duplicates"])

	// save the output in a file - not implemented yet
	analyzeCmd.Flags().StringVarP(&saveFileArg, "save-file", "o", "", output.AnalyzeHelpTexts["saveFile"])
	// No banner will be displayed when this application is launched
	analyzeCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, output.AnalyzeHelpTexts["quiet"])

}
