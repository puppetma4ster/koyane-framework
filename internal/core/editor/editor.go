package editor

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

type EditWordlist struct {
	isSorted   bool
	outputPath string
	tempPath   string
}

func NewEditWordlist(inputPath, outputPath string) (*EditWordlist, error) {
	absoluteInputPath, err := utils.ResolvePath(inputPath)
	if err != nil {
		return nil, err
	}
	absoluteOutputPath, err := utils.ResolvePath(outputPath)
	if err != nil {
		return nil, err
	}

	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return nil, err
	}
	err = utils.CopyFileToTemp(absoluteInputPath, newTempPath)
	if err != nil {
		return nil, err
	}
	return &EditWordlist{
		isSorted:   false,
		outputPath: absoluteOutputPath,
		tempPath:   newTempPath,
	}, nil
}

func (wordlist *EditWordlist) SortWordlist() error {
	var listPath string = wordlist.tempPath
	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return err
	}
	err = utils.ExternalSort(listPath, newTempPath)
	if err != nil {
		return err
	}
	err = os.Remove(listPath)
	if err != nil {
		return err
	}
	wordlist.isSorted = true
	wordlist.tempPath = newTempPath
	return nil
}

func (wordlist *EditWordlist) RemoveWordsWithMask(msk string) error {

	newTempPath, err := utils.GenerateRandomTempPath() //generate new temp path
	if err != nil {
		return err
	}
	currentFile, err := os.Open(wordlist.tempPath) //open old wordlist
	if err != nil {
		return err
	}
	defer currentFile.Close()
	newFile, err := os.Create(newTempPath) // create new Wordlist
	if err != nil {
		return err
	}
	defer newFile.Close()
	mask, err := generator.NewMaskInterpreter(msk)
	if err != nil {
		return err
	}

	writer := bufio.NewWriterSize(newFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	scanner := bufio.NewScanner(currentFile)
	for scanner.Scan() {
		if generator.MatchesWord(mask, scanner.Text()) {
			continue
		}
		_, err = writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			return err
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	err = os.Remove(wordlist.tempPath)
	if err != nil {
		return err
	}
	wordlist.tempPath = newTempPath
	return nil
}

// Uint64Range holds optional min and max values for filtering.
// type Uint64Range struct { Min, Max *uint64 }

// RemoveWordsByRangeUint removes all words from the current wordlist
// whose length in runes is outside the given Uint64Range.
// - If Min is set: words shorter than Min are skipped
// - If Max is set: words longer than Max are skipped
// - If both are nil: everything is kept
func (w *EditWordlist) RemoveWordsByRangeUint(r *utils.Uint64Range) error {
	if r == nil {
		return fmt.Errorf("range must not be nil")
	}
	// Extra safety check: Min must not be greater than Max
	if r.Min != nil && r.Max != nil && *r.Min > *r.Max {
		return fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %d Max: %d", *r.Min, *r.Max)
	}

	// generate a new temporary file path
	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return err
	}

	// create a new file for filtered content
	newFile, err := os.Create(newTempPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = newFile.Close()
	}()

	// open the current wordlist file
	currentFile, err := os.Open(w.tempPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = currentFile.Close()
	}()

	// buffered writer for performance (1 MB buffer)
	writer := bufio.NewWriterSize(newFile, 1024*1024)
	defer writer.Flush()

	// scanner with increased buffer size (default 64 KB → bumped to 1 MB)
	scanner := bufio.NewScanner(currentFile)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		// rune count for proper Unicode length
		l := uint64(utf8.RuneCountInString(line))

		// enforce minimum length
		if r.Min != nil && l < *r.Min {
			continue
		}
		// enforce maximum length
		if r.Max != nil && l > *r.Max {
			continue
		}

		// write line to new file if within range
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// remove old file and swap paths
	if err := os.Remove(w.tempPath); err != nil {
		return err
	}
	w.tempPath = newTempPath
	return nil
}

func (wordlist *EditWordlist) FlushFinishedWordlist() error {
	absolutePath, err := utils.ResolvePath(wordlist.outputPath)
	if err != nil {
		return err
	}
	absolutePath, err = utils.ListPath(absolutePath)
	if err != nil {
		return err
	}
	err = utils.CopyFileToTemp(wordlist.tempPath, absolutePath)
	if err != nil {
		return err
	}
	return nil
}
