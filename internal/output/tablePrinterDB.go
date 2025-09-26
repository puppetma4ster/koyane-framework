package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

var defaultColumnConfigs = []table.ColumnConfig{
	// white
	{Name: "ID", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Name", Align: text.AlignCenter, WidthMax: 40, Transformer: utils.TruncateWithTail("…", 40), ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},

	// yellow (stats)
	{Name: "Entities", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiYellow, text.Bold}},
	{Name: "Smallest", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiYellow, text.Bold}},
	{Name: "Biggest", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiYellow, text.Bold}},
	{Name: "AvgLen", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiYellow, text.Bold}},
	{Name: "AvgEntropy", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiYellow, text.Bold}},

	// green (percent columns)
	{Name: "%Digits", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%Upper", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%Special", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%Digit+Upper", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%Digit+Special", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%Upper+Special", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},
	{Name: "%All", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgHiGreen, text.Bold}},

	// white
	{Name: "Encoding", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Language", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Category", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Author", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Size", Align: text.AlignCenter, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
	{Name: "Link", Align: text.AlignCenter, WidthMax: 20, ColorsHeader: text.Colors{text.BgHiBlack, text.FgWhite, text.Bold}},
}

var minimalBoxStyle = table.BoxStyle{
	BottomLeft:       "",
	BottomRight:      "",
	BottomSeparator:  "",
	EmptySeparator:   "",
	Left:             "",
	LeftSeparator:    "",
	MiddleHorizontal: "-",
	MiddleSeparator:  "",
	MiddleVertical:   "",
	PaddingLeft:      " ",
	PaddingRight:     " ",
	PageSeparator:    "",
	Right:            "",
	RightSeparator:   "",
	TopLeft:          "",
	TopRight:         "",
	TopSeparator:     "-",
	UnfinishedRow:    "",
}

// PrintWordlistsTableAll
// Prints a complete table of all data from the saved word lists
// Parameters:
//   - resaults: all word lists and their related information
func PrintWordlistsTableAll(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Header
	t.AppendHeader(table.Row{
		"ID", "Name", "Entities", "Smallest", "Biggest",
		"AvgLen", "AvgEntropy",
		"%Digits", "%Upper", "%Special",
		"%Digit+Upper", "%Digit+Special", "%Upper+Special", "%All",
		"Encoding", "Language", "Category", "Author",
		"Size", "Link",
	})

	// Rows
	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID,
			wl.Name,
			wl.Entities,
			wl.SmallestEntity,
			wl.BiggestEntity,
			fmt.Sprintf("%.2f", wl.AverageLength),
			fmt.Sprintf("%.2f", wl.AverageEntropy),
			fmt.Sprintf("%.2f", wl.DigitsPercent),
			fmt.Sprintf("%.2f", wl.UpperCasePercent),
			fmt.Sprintf("%.2f", wl.SpecialCharPercent),
			fmt.Sprintf("%.2f", wl.DigitAndUpperCase),
			fmt.Sprintf("%.2f", wl.DigitAndSpecialChar),
			fmt.Sprintf("%.2f", wl.UpperCaseAndSpecialChar),
			fmt.Sprintf("%.2f", wl.DigitUpperCaseAndSpecial),
			wl.Encoding,
			wl.Language,
			wl.Category,
			wl.Author,
			utils.HumanReadableBytes(uint64(wl.Size)),
			wl.Link,
		})
	}

	// Styling
	style := table.StyleDefault

	style.Box = minimalBoxStyle

	style.Options.DrawBorder = false
	style.Options.SeparateColumns = false
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Options.SeparateFooter = false

	t.SetStyle(style)

	t.SetAllowedRowLength(100000000000)

	// Column configs → center-align except Name/Link
	t.SetColumnConfigs(defaultColumnConfigs)

	t.Render()
}

func PrintWordlistsTableDefault(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Header
	t.AppendHeader(table.Row{
		"ID", "Name", "Entities", "Smallest", "Biggest",
		"AvgLen", "%Digits", "%Upper", "%Special",
		"%Digit+Upper", "%Digit+Special", "%Upper+Special", "%All",
		"Category", "Size",
	})

	// Rows
	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID,
			wl.Name,
			wl.Entities,
			wl.SmallestEntity,
			wl.BiggestEntity,
			fmt.Sprintf("%.2f", wl.AverageLength),
			fmt.Sprintf("%.2f", wl.DigitsPercent),
			fmt.Sprintf("%.2f", wl.UpperCasePercent),
			fmt.Sprintf("%.2f", wl.SpecialCharPercent),
			fmt.Sprintf("%.2f", wl.DigitAndUpperCase),
			fmt.Sprintf("%.2f", wl.DigitAndSpecialChar),
			fmt.Sprintf("%.2f", wl.UpperCaseAndSpecialChar),
			fmt.Sprintf("%.2f", wl.DigitUpperCaseAndSpecial),
			wl.Category,
			utils.HumanReadableBytes(uint64(wl.Size)),
		})
	}

	style := table.StyleDefault

	style.Box = minimalBoxStyle

	style.Options.DrawBorder = false
	style.Options.SeparateColumns = false
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Options.SeparateFooter = false

	t.SetStyle(style)

	t.SetAllowedRowLength(100000000000)

	// Column configs → center-align except Name/Link
	t.SetColumnConfigs(defaultColumnConfigs)

	t.Render()
}

// PrintWordlistsTableSummary Prints a compact summary view for quick scanning.
func PrintWordlistsTableSummary(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Entities", "Size", "Category", "Language"})

	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID,
			wl.Name,
			wl.Entities,
			utils.HumanReadableBytes(uint64(wl.Size)),
			wl.Category,
			wl.Language,
		})
	}

	style := table.StyleDefault

	style.Box = minimalBoxStyle

	style.Options.DrawBorder = false
	style.Options.SeparateColumns = false
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Options.SeparateFooter = false

	t.SetStyle(style)

	t.SetAllowedRowLength(10_000)
	t.SetColumnConfigs(defaultColumnConfigs)
	t.Render()
}

// PrintWordlistsTableStats Prints core length/entropy stats per list.
func PrintWordlistsTableStats(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Entities", "Smallest", "Biggest", "AvgLen", "AvgEntropy"})

	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID,
			wl.Name,
			wl.Entities,
			wl.SmallestEntity,
			wl.BiggestEntity,
			fmt.Sprintf("%.2f", wl.AverageLength),
			fmt.Sprintf("%.2f", wl.AverageEntropy),
		})
	}

	style := table.StyleDefault

	style.Box = minimalBoxStyle

	style.Options.DrawBorder = false
	style.Options.SeparateColumns = false
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Options.SeparateFooter = false

	t.SetStyle(style)

	t.SetAllowedRowLength(10_000)

	t.SetColumnConfigs(defaultColumnConfigs)
	t.Render()
}

// PrintWordlistsTableChars Prints character composition percentages per list.
func PrintWordlistsTableChars(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"ID", "Name", "%Digits", "%Upper", "%Special",
		"%Digit+Upper", "%Digit+Special", "%Upper+Special", "%All",
	})

	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID, wl.Name,
			fmt.Sprintf("%.2f", wl.DigitsPercent),
			fmt.Sprintf("%.2f", wl.UpperCasePercent),
			fmt.Sprintf("%.2f", wl.SpecialCharPercent),
			fmt.Sprintf("%.2f", wl.DigitAndUpperCase),
			fmt.Sprintf("%.2f", wl.DigitAndSpecialChar),
			fmt.Sprintf("%.2f", wl.UpperCaseAndSpecialChar),
			fmt.Sprintf("%.2f", wl.DigitUpperCaseAndSpecial),
		})
	}

	style := table.StyleDefault

	style.Box = minimalBoxStyle

	style.Options.DrawBorder = false
	style.Options.SeparateColumns = false
	style.Options.SeparateHeader = true
	style.Options.SeparateRows = true
	style.Options.SeparateFooter = false

	t.SetStyle(style)

	t.SetAllowedRowLength(10_000)
	t.SetColumnConfigs(defaultColumnConfigs)
	t.Render()
}

// DetailedWordlistInfo Detailed infos about wordlists
// method uses a strings.stringbuilder for output printing
// Parameters:
//   - resault: input string von cli Argument
//
// Returns:
//   - *utils.Uint64Range: pointer range Type or nil if s empty
func DetailedWordlistInfo(result *utils.Wordlist) {

	var printer strings.Builder

	printer.WriteString("\n")
	printer.WriteString(utils.PrintDotted("ID", result.ID))
	printer.WriteString(utils.PrintDotted("NAME", result.Name))
	printer.WriteString("\n")

	printer.WriteString(utils.PrintDotted("ENTITIES", result.Entities))
	printer.WriteString(utils.PrintDotted("SMALLEST ENTITY", result.SmallestEntity))
	printer.WriteString(utils.PrintDotted("BIGGEST ENTITY", result.BiggestEntity))
	printer.WriteString(utils.PrintDotted("AVG LENGTH", result.AverageLength))
	printer.WriteString(utils.PrintDotted("AVG ENTROPY", result.AverageEntropy))
	printer.WriteString("\n")

	printer.WriteString(utils.PrintDotted("%DIGITS", result.DigitsPercent))
	printer.WriteString(utils.PrintDotted("%UPPER", result.UpperCasePercent))
	printer.WriteString(utils.PrintDotted("%SPECIAL", result.SpecialCharPercent))
	printer.WriteString(utils.PrintDotted("%DIGIT+UPPER", result.DigitAndUpperCase))
	printer.WriteString(utils.PrintDotted("%DIGIT+SPECIAL", result.DigitAndSpecialChar))
	printer.WriteString(utils.PrintDotted("%UPPER+SPECIAL", result.UpperCaseAndSpecialChar))
	printer.WriteString(utils.PrintDotted("%ALL", result.DigitUpperCaseAndSpecial))
	printer.WriteString("\n")

	printer.WriteString(utils.PrintDotted("ENCODING", result.Encoding))
	printer.WriteString(utils.PrintDotted("%LANGUAGE", result.Language))
	printer.WriteString(utils.PrintDotted("CATEGORY", result.Category))
	printer.WriteString(utils.PrintDotted("%AUTHOR", result.Author))
	printer.WriteString(utils.PrintDotted("%SIZE", result.Size))
	printer.WriteString(utils.PrintDotted("LINK", result.Link))
	printer.WriteString("\n")

	// printer.WriteString(utils.PrintDotted("INFO", results.Info))
	// for future wordlist info like license, history, maybe ratings

	fmt.Println(printer.String())
}
