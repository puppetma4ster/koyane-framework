package editor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func (wordlist *EditWordlist) RemoveWordsByRange(min, max string) error {
	var minEndless bool = false
	var maxEndless bool = false
	const endlessSymbol string = "*"

	var minNumber int
	var maxNumber int
	if min == "*" {
		minEndless = true
	} else if imin, err := strconv.Atoi(min); err == nil {
		minNumber = imin
	} else {
		return fmt.Errorf("the Minimal Range is not a number or \"*\": %s", min)
	}

	if max == "*" {
		maxEndless = true
	} else if imax, err := strconv.Atoi(max); err == nil {
		maxNumber = imax
	} else {
		return fmt.Errorf("the Minimal Range is not a number or \"*\": %s", max)
	}

	if !minEndless && !maxEndless {
		if minNumber > maxNumber {
			return fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %d Max: %d", minNumber, maxNumber)
		}
	}
	newTempPath, err := utils.GenerateRandomTempPath() //generate new temp path
	if err != nil {
		return err
	}
	newFile, err := os.Create(newTempPath) // create new Wordlist
	if err != nil {
		return err
	}
	defer newFile.Close()

	currentFile, err := os.Open(wordlist.tempPath)
	if err != nil {
		return err
	}
	defer currentFile.Close()

	writer := bufio.NewWriterSize(newFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	scanner := bufio.NewScanner(currentFile)
	for scanner.Scan() {
		if !minEndless {
			if minNumber > utf8.RuneCountInString(scanner.Text()) {
				continue
			}
		}
		if !maxEndless {
			if maxNumber < utf8.RuneCountInString(scanner.Text()) {
				continue
			}
		}
		_, err2 := writer.WriteString(scanner.Text() + "\n")
		if err2 != nil {
			return err2
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
