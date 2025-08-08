package output

import (
	"fmt"
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
	contentText.WriteString("---------- Content Information ----------")
	contentText.WriteString("\n")

	generalText.WriteString("\n")
	generalText.WriteString("---------- General Information ----------")
	generalText.WriteString("\n")

	statsText.WriteString("\n")
	statsText.WriteString("---------- Wordlist Stats ----------")
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
	printDotted("Average Word Length", wordlist.content.AvWordLen)
}

func (wordlist *AnalyzePrinter) PrintCharStatistics() {
	printDotted("Average Word Length", wordlist.content.AvWordLen)
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

func printDotted(label string, value interface{}) string {
	totalWidth := 25
	dots := totalWidth - len(label)
	if dots < 0 {
		dots = 0
	}
	return fmt.Sprintf("%s%s: %v\n", label, strings.Repeat(".", dots), value)
}
