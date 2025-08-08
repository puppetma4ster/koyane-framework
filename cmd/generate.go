package cmd

import (
	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
	"os"
)

var (
	mask      string
	minLength int
)

var generateCmd = &cobra.Command{
	Use:   output.GenerateMessages["use"],
	Short: output.GenerateMessages["short"],
	Long:  output.GenerateMessages["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath := args[0]

		// when generate is called
		if minLength == 0 {
			entities, bytes, err := generator.CalculateMaskStorage(mask)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
			output.PrintStatus("statusGenerator", "calculateWords", entities)
			output.PrintStatus("statusGenerator", "calculateSize", utils.HumanReadableBytes(bytes))
			output.PrintStatus("statusGenerator", "buildingMaskWordlist", mask)
			err = generator.GenerateMaskWordlist(mask, outputPath)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		} else {
			err := generator.GenerateMaskWordlist(mask, outputPath, minLength)
			if err != nil {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			} else {
				output.PrintError("errors", "error", err)
				os.Exit(1)
			}
		}
		oPath, err := utils.ResolvePath(outputPath)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		output.PrintSuccess("successGenerator", "wordlistCreated", oPath+utils.ListSuffix)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&mask, "mask", "M", "", output.GenerateMessages["mask"])
	err := generateCmd.MarkFlagRequired("mask")
	if err != nil {
		output.PrintError("errors", "error", err)
		os.Exit(1)
	}
	generateCmd.Flags().IntVarP(&minLength, "min-length", "m", 0, output.GenerateMessages["minLength"])

}
