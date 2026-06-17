package app

import (
	"os"

	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var (
	mask              string
	minLength         int
	extractHcPotfiles string
	permutationFile   string
	maxLength         int
	outputPath        string
	stdout            bool
	quiet             bool
)

// generateCmd has implemented Kyfgen's CLI logic.
// It offers three different generation options:
//   - mask,
//   - Hashcat potfile extraction,
//   - permutation.
var generateCmd = &cobra.Command{
	Use:   output.GenerateMessages["use"],
	Short: output.GenerateMessages["short"],
	Long:  output.GenerateMessages["long"],
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// print banner
		if !quiet && !stdout {
			output.PrintBanner()
		}
		err := utils.CreateTempDir() // generate tmp dir if no tmp dir is found
		if err != nil {              // creates temp folder to /tmp/koyane_framework_tmp
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		if !stdout && outputPath == "" { // printing error if no output variant is chosen
			output.PrintError("errors", "noOutputPath", "")
			os.Exit(1)
		}

		outputPath, err := utils.OverwriteFile(outputPath) // checks if output file is existing
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		stop := make(chan struct{}) // spinner animation
		if !stdout {
			go output.Spinner("Generate File", stop)
		}

		if extractHcPotfiles != "" { // hashcat password extraction is chosen
			err := generator.ExtractHashCatPotfile(extractHcPotfiles, outputPath)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		} else if mask != "" { // mask generation is chosen
			// mask generation without minlen called
			if minLength == 0 {
				// Estimate memory usage
				entities, bytes, err := generator.CalculateMaskStorage(mask)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}

				if !stdout { // output some storage information's...
					output.PrintStatus("statusGenerator", "calculateWords", utils.FormatUint64WithCommas(entities))
					output.PrintStatus("statusGenerator", "calculateSize", utils.HumanReadableBytes(bytes))
					output.PrintStatus("statusGenerator", "buildingMaskWordlist", mask)
				}

				// start writing into file
				err = generator.ConcurrentGenerateMaskWordlist(mask, outputPath)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}

			} else { // mask generation with minlen called
				// Estimate memory usage
				entities, bytes, err := generator.CalculateMaskStorage(mask)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}

				if !stdout { // output some storage information's...
					output.PrintStatus("statusGenerator", "calculateWords", utils.FormatUint64WithCommas(entities))
					output.PrintStatus("statusGenerator", "calculateSize", utils.HumanReadableBytes(bytes))
					output.PrintStatus("statusGenerator", "buildingMaskWordlist", mask)
				}

				// start writing into file
				err = generator.ConcurrentGenerateMaskWordlist(mask, outputPath, minLength)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			}
		} else if permutationFile != "" && maxLength != 0 { // permutation chosen

			// error printing
			if maxLength <= 0 {
				output.PrintError("errors", "wrongMaxLength", "")
				os.Exit(1)
			}
			if minLength > maxLength && maxLength < 0 {
				output.PrintError("errors", "wrongMinLength", "")
			}

			// start permutation
			err := generator.ConcurrentPermutation(permutationFile, outputPath, minLength, maxLength, 0)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		}
		if !stdout {
			oPath, err := utils.ResolvePath(outputPath)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}

			close(stop) // stopping spinner animation

			output.PrintSuccess("successGenerator", "wordlistCreated", oPath+utils.ListSuffix)

		}
	},
}

// Execute gets called by main.main() to start the CLI. It only needs to happen once to the rootCmd.
func Execute() {
	err := generateCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// takes the mask to be generated
	generateCmd.Flags().StringVarP(&mask, "mask", "M", "", output.GenerateMessages["mask"])
	// takes a file as input and permutes each line
	generateCmd.Flags().StringVarP(&permutationFile, "permutation", "p", "", output.GenerateMessages["permutation"])
	// accepts a Hashcat potfile for password extraction
	generateCmd.Flags().StringVar(&extractHcPotfiles, "extract-hc-potfile", "", output.GenerateEditHelpTexts["extract-hc-potfile"])

	// accepts the minimum number of characters for mask and permutation generation
	generateCmd.Flags().IntVar(&minLength, "min", 0, output.GenerateMessages["minLength"])
	// accepts the maximum number of characters for mask and permutation generation
	generateCmd.Flags().IntVarP(&maxLength, "max", "m", 0, output.GenerateMessages["maxLength"])

	// outputs the generated strings directly to the shell
	generateCmd.Flags().BoolVar(&stdout, "stdout", false, output.GenerateEditHelpTexts["stdout"])
	// No banner will be displayed when this application is launched
	generateCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, output.GenerateEditHelpTexts["quiet"])
	// accepts your output path where the strings should be saved
	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", output.GenerateEditHelpTexts["output"])

}
