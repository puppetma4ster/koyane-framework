package cmd

import (
	"os"
	"strconv"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/core/wordlistDB"
	"github.com/puppetma4ster/koyane-framework/internal/output"
	"github.com/spf13/cobra"
)

var pathArg string

// getCmd Implementation to download word lists stored in the database.
// Word lists are identified by their entry ID.
// The default download path is taken from the YAML configuration file.
// If a specific path is to be used, it can be specified with the “path” flag.
// If this does not exist, you will be asked if you want to use the default path. Answer with (y)es or (n)o.
var getCmd = &cobra.Command{
	Use:   output.GetHelpTexts["use"],
	Short: output.GetHelpTexts["short"],
	Long:  output.GetHelpTexts["long"],
	Args:  cobra.ExactArgs(1), // database id for choosing download
	Run: func(cmd *cobra.Command, args []string) {

		id, err := strconv.ParseUint(args[0], 10, 64) // convert string into uint64
		if err != nil {
			output.PrintError("errors", "invID", args[0])
			os.Exit(1)
		}
		cfg, err := utils.LoadConfig("config.yaml") // loading yaml for db path
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		db, err := wordlistDB.NewWordlistRepository(cfg.General.DatabasePath) //	loading db repo
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		resault, err := db.FindByID(id) // Load entry from database
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		var path string = cfg.General.DefaultWordlistPath
		if pathArg != "" {
			path = pathArg
		}
		absPath, err := utils.ResolvePath(path)
		if err != nil { //Test whether the path exists. If it does not exist, ask whether the YAML path should be used.
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}

		output.PrintStatus("statusGet", "try", resault.Link)
		err = wordlistDB.DownloadWordlist(resault.Link, absPath, cfg.General.UserAgent)
		if err != nil {
			output.PrintError("errors", "error", err)
			os.Exit(1)
		}
		output.PrintSuccess("successGet", "succeeded", absPath)
	},
}

// init gives the command to root and implements all command flags
func init() {
	rootCmd.AddCommand(getCmd)
	// path flag to specify a custom path if the path in the yaml is not to be used
	getCmd.Flags().StringVarP(&pathArg, "path", "p", "", output.GetHelpTexts["path"])
}
