package editor

import (
	"bufio"
	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"os"
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

func SortWordlist(wordlist *EditWordlist) (*EditWordlist, error) {
	var listPath string = wordlist.tempPath
	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return nil, err
	}
	err = utils.ExternalSort(listPath, newTempPath)
	if err != nil {
		return nil, err
	}
	err = os.Remove(listPath)
	if err != nil {
		return nil, err
	}
	wordlist.isSorted = true
	wordlist.tempPath = newTempPath
	return wordlist, nil
}

func RemoveWordsWithMask(wordlist *EditWordlist, msk string) (*EditWordlist, error) {

	newTempPath, err := utils.GenerateRandomTempPath() //generate new temp path
	if err != nil {
		return nil, err
	}
	currentFile, err := os.Open(wordlist.tempPath) //open old wordlist
	if err != nil {
		return nil, err
	}
	defer currentFile.Close()
	newFile, err := os.Create(newTempPath) // create new Wordlist
	if err != nil {
		return nil, err
	}
	defer newFile.Close()
	mask, err := generator.NewMaskInterpreter(msk)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	err = os.Remove(wordlist.tempPath)
	if err != nil {
		return nil, err
	}
	wordlist.tempPath = newTempPath
	return wordlist, nil
}

func FlushFinishedWordlist(wordlist *EditWordlist) error {
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
