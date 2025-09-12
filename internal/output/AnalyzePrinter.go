package output

import (
	"fmt"
	"strconv"
	"strings"

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
	contentText.WriteString("\n")

	generalText.WriteString("\n")
	generalText.WriteString("-------------------- General Information --------------------")
	generalText.WriteString("\n")
	generalText.WriteString("\n")

	statsText.WriteString("\n")
	statsText.WriteString("-------------------- Wordlist Stats --------------------")
	statsText.WriteString("\n")
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
	wordlist.generalText.WriteString(utils.PrintDotted("Name", wordlist.general.FileName))
}

func (wordlist *AnalyzePrinter) PrintFilePath() {
	wordlist.generalText.WriteString(utils.PrintDotted("Path", wordlist.general.FilePath))
}

func (wordlist *AnalyzePrinter) PrintFileSize() {
	size := wordlist.general.FileSize
	formattedSize := fmt.Sprintf("%d / %s", size, utils.HumanReadableBytes(size))

	wordlist.generalText.WriteString(utils.PrintDotted("Size", formattedSize))
}

func (wordlist *AnalyzePrinter) PrintExtension() {
	wordlist.generalText.WriteString(utils.PrintDotted("Extension", wordlist.general.Extension))
}

func (wordlist *AnalyzePrinter) PrintEncoding() {
	wordlist.generalText.WriteString(utils.PrintDotted("Estimated Encoding", wordlist.general.Encoding))
}

func (wordlist *AnalyzePrinter) PrintHashValue() {
	wordlist.generalText.WriteString(utils.PrintDotted("MD5 Hash", wordlist.general.HashVal))
}

func (wordlist *AnalyzePrinter) PrintLastModified() {
	wordlist.generalText.WriteString(utils.PrintDotted("Last Modified", wordlist.general.LastChanges))
}

// Content Words
func (wordlist *AnalyzePrinter) PrintWordLInes() {
	wordlist.contentText.WriteString(utils.PrintDotted("ENTITIES", wordlist.content.WordLines))
}

func (wordlist *AnalyzePrinter) PrintSmallestWordLen() {
	wordlist.contentText.WriteString(utils.PrintDotted("SMALLEST ENTITY", strconv.Itoa(wordlist.content.SmallestWordLen)+"  ->  "+wordlist.content.SmallestWordStr))
}

func (wordlist *AnalyzePrinter) PrintBiggestWordLen() {
	wordlist.contentText.WriteString(utils.PrintDotted("BIGGEST ENTITY", strconv.Itoa(wordlist.content.BiggestWordLen)+"  ->  "+wordlist.content.BiggestWordStr))
}

func (wordlist *AnalyzePrinter) PrintAvWordLen() {
	wordlist.contentText.WriteString(utils.PrintDotted("AVG WORD LENGTH", wordlist.content.AvWordLen))
}

func (wordlist *AnalyzePrinter) PrintAvEntropy() {
	wordlist.contentText.WriteString(utils.PrintDotted("AVG WORD ENTROPY", wordlist.content.AvEntropy))
}

func (wordlist *AnalyzePrinter) PrintHasDuplicates() {
	wordlist.contentText.WriteString(utils.PrintDotted("HAS DUPLICATES", wordlist.content.HasDuplicates))
}

// stats words
func (wordlist *AnalyzePrinter) PrintDuplicateWords() {
	var result interface{}
	if len(wordlist.content.DuplicateWords) == 0 {
		result = "no duplicates in wordlist"
	} else {
		result = strings.Join(wordlist.content.DuplicateWords, ",")
	}
	wordlist.statsText.WriteString(utils.PrintDotted("DUPLICATES", result))

}

func (wordlist *AnalyzePrinter) PrintCharStatistics() {
	wordlist.statsText.WriteString(utils.PrintDotted("CHARACTER STATISTICS", "\n"))
	wordlist.statsText.WriteString(utils.PrintRuneMapColumns(wordlist.content.CharCount, 6))
}

func (wordlist *AnalyzePrinter) PrintWordStats() {
	wordlist.statsText.WriteString(utils.PrintDotted("%DIGITS", wordlist.content.WordsWDigitsPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%UPPER", wordlist.content.WordsWUpperPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%SPECIAL", wordlist.content.WordsWSpecCharPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%DIGIT+UPPER", wordlist.content.WordsWDigitUpperPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%DIGIT+SPECIAL", wordlist.content.WordsWDigitSpecPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%UPPER+SPECIA", wordlist.content.WordsWUpperSpecPercent, "%"))
	wordlist.statsText.WriteString(utils.PrintDotted("%ALL", wordlist.content.WordsWDigitUpperSpecPercent, "%"))
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
