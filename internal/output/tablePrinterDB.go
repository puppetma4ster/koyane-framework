package output

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

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
			wl.AverageLength,
			wl.AverageEntropy,
			wl.DigitsPercent,
			wl.UpperCasePercent,
			wl.SpecialCharPercent,
			wl.DigitAndUpperCase,
			wl.DigitAndSpecialChar,
			wl.UpperCaseAndSpecialChar,
			wl.DigitUpperCaseAndSpecial,
			wl.Encoding,
			wl.Language,
			wl.Category,
			wl.Author,
			utils.HumanReadableBytes(uint64(wl.Size)),
			wl.Link,
		})
	}

	// Styling
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Style().Format.Header = text.FormatUpper
	t.Style().Options.SeparateRows = true

	t.SetAllowedRowLength(100000000000)

	// Column configs → center-align except Name/Link
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Align: text.AlignCenter},
		{Name: "Entities", Align: text.AlignCenter},
		{Name: "Smallest", Align: text.AlignCenter},
		{Name: "Biggest", Align: text.AlignCenter},
		{Name: "AvgLen", Align: text.AlignCenter},
		{Name: "AvgEntropy", Align: text.AlignCenter},
		{Name: "%Digits", Align: text.AlignCenter},
		{Name: "%Upper", Align: text.AlignCenter},
		{Name: "%Special", Align: text.AlignCenter},
		{Name: "%Digit+Upper", Align: text.AlignCenter},
		{Name: "%Digit+Special", Align: text.AlignCenter},
		{Name: "%Upper+Special", Align: text.AlignCenter},
		{Name: "%All", Align: text.AlignCenter},
		{Name: "Encoding", Align: text.AlignCenter},
		{Name: "Language", Align: text.AlignCenter},
		{Name: "Category", Align: text.AlignCenter},
		{Name: "Author", Align: text.AlignCenter},
		{Name: "Size", Align: text.AlignCenter},
		{Name: "Link", WidthMax: 20},
	})

	t.Render()
}

func PrintWordlistsTableDefault(results []utils.Wordlist) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Header
	t.AppendHeader(table.Row{
		"ID", "Name", "Entities", "Smallest", "Biggest",
		"AvgLen", "AvgEntropy",
		"%Digits", "%Upper", "%Special",
		"%Digit+Upper", "%Digit+Special", "%Upper+Special", "%All",
		"Encoding", "Language", "Category",
		"Size",
	})

	// Rows
	for _, wl := range results {
		t.AppendRow(table.Row{
			wl.ID,
			wl.Name,
			wl.Entities,
			wl.SmallestEntity,
			wl.BiggestEntity,
			wl.AverageLength,
			wl.AverageEntropy,
			wl.DigitsPercent,
			wl.UpperCasePercent,
			wl.SpecialCharPercent,
			wl.DigitAndUpperCase,
			wl.DigitAndSpecialChar,
			wl.UpperCaseAndSpecialChar,
			wl.DigitUpperCaseAndSpecial,
			wl.Encoding,
			wl.Language,
			wl.Category,
			utils.HumanReadableBytes(uint64(wl.Size)),
		})
	}

	// Styling
	t.SetStyle(table.StyleColoredYellowWhiteOnBlack)
	t.Style().Format.Header = text.FormatUpper
	t.Style().Options.SeparateRows = true

	t.SetAllowedRowLength(100000000000)

	// Column configs → center-align except Name/Link
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Align: text.AlignCenter},
		{Name: "Entities", Align: text.AlignCenter},
		{Name: "Smallest", Align: text.AlignCenter},
		{Name: "Biggest", Align: text.AlignCenter},
		{Name: "AvgLen", Align: text.AlignCenter},
		{Name: "AvgEntropy", Align: text.AlignCenter},
		{Name: "%Digits", Align: text.AlignCenter},
		{Name: "%Upper", Align: text.AlignCenter},
		{Name: "%Special", Align: text.AlignCenter},
		{Name: "%Digit+Upper", Align: text.AlignCenter},
		{Name: "%Digit+Special", Align: text.AlignCenter},
		{Name: "%Upper+Special", Align: text.AlignCenter},
		{Name: "%All", Align: text.AlignCenter},
		{Name: "Encoding", Align: text.AlignCenter},
		{Name: "Language", Align: text.AlignCenter},
		{Name: "Category", Align: text.AlignCenter},
		{Name: "Size", Align: text.AlignCenter},
	})

	t.Render()
}
