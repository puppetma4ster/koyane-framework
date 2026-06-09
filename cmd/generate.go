package cmd

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
)

var generateCmd = &cobra.Command{
	Use:   output.GenerateMessages["use"],
	Short: output.GenerateMessages["short"],
	Long:  output.GenerateMessages["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]

		outputPath, err := utils.OverwriteFile(outputPath) // checks if output file is existing
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		stop := make(chan struct{}) // spinner animation
		go output.Spinner("Generate File", stop)

		if extractHcPotfiles != "" { // hashcat password extraction is chosen
			err := generator.ExtractHashCatPotfile(extractHcPotfiles, outputPath)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		} else if mask != "" { // mask generation is chosen
			// when generate is called
			if minLength == 0 {
				entities, bytes, err := generator.CalculateMaskStorage(mask)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				output.PrintStatus("statusGenerator", "calculateWords", utils.FormatUint64WithCommas(entities))
				output.PrintStatus("statusGenerator", "calculateSize", utils.HumanReadableBytes(bytes))
				output.PrintStatus("statusGenerator", "buildingMaskWordlist", mask)

				err = generator.ConcurrentGenerateMaskWordlist(mask, outputPath)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			} else {
				entities, bytes, err := generator.CalculateMaskStorage(mask)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
				output.PrintStatus("statusGenerator", "calculateWords", utils.FormatUint64WithCommas(entities))
				output.PrintStatus("statusGenerator", "calculateSize", utils.HumanReadableBytes(bytes))
				output.PrintStatus("statusGenerator", "buildingMaskWordlist", mask)

				err = generator.ConcurrentGenerateMaskWordlist(mask, outputPath, minLength)
				if err != nil {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				} else {
					output.PrintError("errors", "error", err)
					os.Exit(1)
				}
			}
		} else if permutationFile != "" && maxLength != 0 { // permutation chosen
			if maxLength <= 0 {
				output.PrintError("errors", "wrongMaxLength", "")
				os.Exit(1)
			}
			if minLength > maxLength && maxLength < 0 {
				output.PrintError("errors", "wrongMinLength", "")
			}
			err := generator.ConcurrentPermutation(permutationFile, outputPath, minLength, maxLength, 0)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		}
		oPath, err := utils.ResolvePath(outputPath)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		close(stop) // stopping spinner animation
		output.PrintSuccess("successGenerator", "wordlistCreated", oPath+utils.ListSuffix)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&mask, "mask", "M", "", output.GenerateMessages["mask"])

	generateCmd.Flags().IntVar(&minLength, "min", 0, output.GenerateMessages["minLength"])
	generateCmd.Flags().StringVarP(&permutationFile, "permutation", "p", "", output.GenerateMessages["permutation"])
	generateCmd.Flags().IntVarP(&maxLength, "max", "m", 0, output.GenerateMessages["maxLength"])
	generateCmd.Flags().StringVar(&extractHcPotfiles, "extract-hc-potfile", "", output.GenerateEditHelpTexts["extract-hc-potfile"])

}
