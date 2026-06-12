package app

import (
	"os"
	"strconv"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/core/wordlistDB"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   output.ShowHelpTexts["use"],
	Short: output.ShowHelpTexts["short"],
	Long:  output.ShowHelpTexts["long"],
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id, err := strconv.ParseUint(args[0], 10, 64) // convert string into uint64
		if err != nil {
			output.PrintError("errors", "invID", args[0])
			os.Exit(1)
		}
		dbPath, err := utils.GetDatabasePath()
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		db, err := wordlistDB.NewWordlistRepository(dbPath) //	loading db repo
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		resault, err := db.FindByID(id)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		output.DetailedWordlistInfo(resault)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
