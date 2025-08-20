package output

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/puppetma4ster/koyane-framework/internal/core/analyzer"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

type AnalyzePrinter struct {
	generalText *strings.Builder
	general     *analyzer.GeneralAnalyzer
	contentText *strings.Builder
	statsText   *strings.Builder
	content     *analyzer.AnalyzerContent
}

func NewAnalyzePrinter(generalAnalyzer *analyzer.GeneralAnalyzer, contentAnalyzer *analyzer.AnalyzerContent) *AnalyzePrinter {
	contentText := &strings.Builder{}
	generalText := &strings.Builder{}
	statsText := &strings.Builder{}

	contentText.WriteString("\n")
	contentText.WriteString("-------------------- Content Information --------------------")
	contentText.WriteString("\n")

	generalText.WriteString("\n")
	generalText.WriteString("-------------------- General Information --------------------")
	generalText.WriteString("\n")

	statsText.WriteString("\n")
	statsText.WriteString("-------------------- Wordlist Stats --------------------")
	statsText.WriteString("\n")

	return &AnalyzePrinter{
		generalText: generalText,
		contentText: contentText,
		statsText:   statsText,
		general:     generalAnalyzer,
		content:     contentAnalyzer,
	}
}

func (wordlist *AnalyzePrinter) PrintAllGeneralInfo() {
	wordlist.PrintFileName()
	wordlist.PrintFilePath()
	wordlist.PrintFileSize()
	wordlist.PrintExtension()
	wordlist.PrintEncoding()
	wordlist.PrintHashValue()
	wordlist.PrintLastModified()
}
func (wordlist *AnalyzePrinter) PrintAllContentInfo() {
	wordlist.PrintWordLInes()
	wordlist.PrintSmallestWordLen()
	wordlist.PrintBiggestWordLen()
	wordlist.PrintAvWordLen()
	wordlist.PrintAvEntropy()
	wordlist.PrintHasDuplicates()
}
func (wordlist *AnalyzePrinter) PrintAllStatsInfo() {
	wordlist.PrintWordStats()
	wordlist.PrintDuplicateWords()
	wordlist.PrintCharStatistics()
}

// general words
func (wordlist *AnalyzePrinter) PrintFileName() {
	wordlist.generalText.WriteString(printDotted("File Name", wordlist.general.FileName))
}

func (wordlist *AnalyzePrinter) PrintFilePath() {
	wordlist.generalText.WriteString(printDotted("File Path", wordlist.general.FilePath))
}

func (wordlist *AnalyzePrinter) PrintFileSize() {
	size := wordlist.general.FileSize
	formattedSize := fmt.Sprintf("%d / %s", size, utils.HumanReadableBytes(size))

	wordlist.generalText.WriteString(printDotted("File Size", formattedSize))
}

func (wordlist *AnalyzePrinter) PrintExtension() {
	wordlist.generalText.WriteString(printDotted("File Extension", wordlist.general.Extension))
}

func (wordlist *AnalyzePrinter) PrintEncoding() {
	wordlist.generalText.WriteString(printDotted("Estimated Encoding", wordlist.general.Encoding))
}

func (wordlist *AnalyzePrinter) PrintHashValue() {
	wordlist.generalText.WriteString(printDotted("File Hash", wordlist.general.HashVal))
}

func (wordlist *AnalyzePrinter) PrintLastModified() {
	wordlist.generalText.WriteString(printDotted("Last Modified", wordlist.general.LastChanges))
}

// Content Words
func (wordlist *AnalyzePrinter) PrintWordLInes() {
	wordlist.contentText.WriteString(printDotted("Total Words", wordlist.content.WordLines))
}

func (wordlist *AnalyzePrinter) PrintSmallestWordLen() {
	wordlist.contentText.WriteString(printDotted("Smallest Word Length", wordlist.content.SmallestWordLen))
}

func (wordlist *AnalyzePrinter) PrintBiggestWordLen() {
	wordlist.contentText.WriteString(printDotted("Biggest Word Length", wordlist.content.BiggestWordLen))
}

func (wordlist *AnalyzePrinter) PrintAvWordLen() {
	wordlist.contentText.WriteString(printDotted("Average Word Length", wordlist.content.AvWordLen))
}

func (wordlist *AnalyzePrinter) PrintAvEntropy() {
	wordlist.contentText.WriteString(printDotted("Average Word Entropy", wordlist.content.AvEntropy))
}

func (wordlist *AnalyzePrinter) PrintHasDuplicates() {
	wordlist.contentText.WriteString(printDotted("Has Duplicates", wordlist.content.HasDuplicates))
}

// stats words
func (wordlist *AnalyzePrinter) PrintDuplicateWords() {
	var result interface{}
	if len(wordlist.content.DuplicateWords) == 0 {
		result = "no duplicates in wordlist"
	} else {
		result = strings.Join(wordlist.content.DuplicateWords, ",")
	}
	wordlist.statsText.WriteString(printDotted("Duplicate words", result))

}

func (wordlist *AnalyzePrinter) PrintCharStatistics() {
	wordlist.statsText.WriteString(printDotted("Character statistics", "\n"))
	wordlist.statsText.WriteString(printRuneMapColumns(wordlist.content.CharCount, 6))
}

func (wordlist *AnalyzePrinter) PrintWordStats() {
	wordlist.statsText.WriteString(printDotted("Contains digits", wordlist.content.WordsWDigitsPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains upper case letters", wordlist.content.WordsWUpperPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains special chars", wordlist.content.WordsWSpecCharPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains upper case & digits", wordlist.content.WordsWDigitUpperPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains digits & special chars", wordlist.content.WordsWDigitSpecPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains upper case & special chars", wordlist.content.WordsWUpperSpecPercent, "%"))
	wordlist.statsText.WriteString(printDotted("Contains Digits Upper, Case & special chars", wordlist.content.WordsWDigitUpperSpecPercent, "%"))
}

func (wordlist *AnalyzePrinter) FlushGeneral() {
	fmt.Println(wordlist.generalText.String())
}

func (wordlist *AnalyzePrinter) FlushContent() {
	fmt.Println(wordlist.contentText.String())
}

func (wordlist *AnalyzePrinter) FlushStats() {
	fmt.Println(wordlist.statsText.String())
}

func printDotted(label string, value interface{}, unit ...string) string {
	totalWidth := 45
	dots := totalWidth - len(label)
	if dots < 0 {
		dots = 0
	}

	var valueStr string
	switch v := value.(type) {
	case float32:
		valueStr = fmt.Sprintf("%.2f", v)
	case float64:
		valueStr = fmt.Sprintf("%.2f", v)
	case int:
		valueStr = fmt.Sprintln(humanize.Comma(int64(v)))
	case uint64:
		valueStr = fmt.Sprintln(humanize.Comma(int64(v)))
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	// Falls eine Einheit übergeben wurde, diese anhängen
	unitStr := ""
	if len(unit) > 0 {
		unitStr = unit[0]
	}

	return fmt.Sprintf("%s%s: %s%s\n", label, strings.Repeat(".", dots), valueStr, unitStr)
}

func printRuneMapColumns(m map[rune]uint64, columns int) string {
	var builder strings.Builder
	totalWidth := 12 // Breite pro Spalte (z.B. 1 Zeichen Key + 8 Punkte + ": " + max 2-3 Ziffern)

	i := 0
	for k, v := range m {
		keyStr := string(k)
		dots := totalWidth - len(keyStr) - len(fmt.Sprintf(": %d", v))
		if dots < 0 {
			dots = 0
		}

		// Format: key + Punkte + ": " + Wert + 4 Leerzeichen als Abstand
		fmtStr := fmt.Sprintf("%s%s: %d    ", keyStr, strings.Repeat(".", dots), v)
		builder.WriteString(fmtStr)

		i++
		if i%columns == 0 {
			builder.WriteString("\n")
		}
	}
	if i%columns != 0 {
		builder.WriteString("\n")
	}

	return builder.String()
}
