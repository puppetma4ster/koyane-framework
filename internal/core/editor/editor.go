package editor

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

type EditWordlist struct {
	isSorted   bool
	outputPath string
	tempPaths  []string
}

// NewEditWordlist Creates a new struct for wordlist processing.
// The method converts each path to an absolute path.
// The list is then copied to the project's temp path.
// There, the list is split into several parts to enable parallel processing.
//
// Parameters:
//   - inputPath: path to the list to be edited
//   - outputPath: where the edited list should be saved
//
// - delOriginal: Deletes the list that was to be edited so that only the new edited one remains.
//
// Returns:
//   - EditWordlist: struct with wordlist information
//   - error: if  path resolve problems, temp path generateing problems, copy problems,
//     wordlist splitting problems, deleting problems
func NewEditWordlist(inputPath, outputPath string, delOriginal bool) (*EditWordlist, error) {
	absoluteInputPath, err := utils.ResolvePath(inputPath)
	if err != nil {
		return nil, err
	}
	absoluteOutputPath, err := utils.ResolvePath(outputPath) // saves absolute path
	if err != nil {
		return nil, err
	}

	tempList, err := utils.SplitWordlist(absoluteInputPath) // splits wordlist and saves to temp dir
	if err != nil {
		return nil, err
	}
	if delOriginal { // if flag is used for deleting the original file
		err = os.Remove(absoluteInputPath)
		if err != nil {
			return nil, err
		}
	}
	return &EditWordlist{ // creating struct
		isSorted:   false,
		outputPath: absoluteOutputPath,
		tempPaths:  tempList,
	}, nil
}

// ConcurrentSortWordlist parallelizes the sorting process
// All temp lists are sorted and the new string is captured in channelNewPaths.
// channelNewPaths is then passed to the struct as a new array and the isSorted switch is set to true.
func (wordlist *EditWordlist) ConcurrentSortWordlist() {
	var threadPool sync.WaitGroup
	channelNewPaths := make(chan string) // channel to catch new paths
	for _, unSortPath := range wordlist.tempPaths {
		threadPool.Add(1)
		go func(p string) {
			defer threadPool.Done()
			result, err := sortWordlist(p)
			if err != nil {
				panic(err)
			}
			channelNewPaths <- result
		}(unSortPath)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewPaths)
	}()

	var newPaths []string
	for path := range channelNewPaths {
		newPaths = append(newPaths, path)
	}
	wordlist.isSorted = true
	wordlist.tempPaths = newPaths
}

// SortWordlist sorts a wordlist with an external sort algorithm
//
// Parameters:
//   - inputPath: path to the list to be edited
//   - outputPath: where the edited list should be saved
//
// - delOriginal: Deletes the list that was to be edited so that only the new edited one remains.
//
// Returns:
//   - EditWordlist: struct with wordlist information
//   - error: if  path resolve problems, temp path generateing problems, copy problems,
//     wordlist splitting problems, deleting problems
func sortWordlist(inputPath string) (string, error) {
	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return "", err
	}
	err = utils.ExternalSort(inputPath, newTempPath)
	if err != nil {
		return "", err
	}
	err = os.Remove(inputPath)
	if err != nil {
		return "", err
	}
	return newTempPath, nil
}

func (wordlist *EditWordlist) ConcurrentRemoveWordsWithMask(msk string) error {
	var threadPool sync.WaitGroup
	channelNewPaths := make(chan string) // channel to catch new paths

	mask, err := generator.NewMaskInterpreter(msk)
	if err != nil {
		return err
	}

	for _, unRemoved := range wordlist.tempPaths {
		threadPool.Add(1)
		go func(m *generator.MaskInterpreter, p string) {
			defer threadPool.Done()
			newPath, err := removeWordsWithMask(m, p)
			if err != nil {
				panic(err)
			}
			channelNewPaths <- newPath
		}(mask, unRemoved)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewPaths)
	}()

	var newPaths []string
	for path := range channelNewPaths {
		newPaths = append(newPaths, path)
	}

	wordlist.tempPaths = newPaths
	return nil
}
func removeWordsWithMask(mask *generator.MaskInterpreter, inputPath string) (string, error) {

	newTempPath, err := utils.GenerateRandomTempPath() //generate new temp path
	if err != nil {
		return "", err
	}
	currentFile, err := os.Open(inputPath) //open old wordlist
	if err != nil {
		return "", err
	}
	defer currentFile.Close()
	newFile, err := os.Create(newTempPath) // create new Wordlist
	if err != nil {
		return "", err
	}
	defer newFile.Close()

	writer := bufio.NewWriterSize(newFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	scanner := bufio.NewScanner(currentFile)
	for scanner.Scan() {
		if generator.MatchesWord(mask, scanner.Text()) {
			continue
		}
		_, err = writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			return "", err
		}
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}

	err = os.Remove(inputPath)
	if err != nil {
		return "", err
	}
	return newTempPath, nil
}

func (wordlist *EditWordlist) ConcurrentRemoveWordsByRange(r string) error {
	var threadPool sync.WaitGroup
	channelNewPaths := make(chan string) // channel to catch new paths

	uintRange, err := utils.NewUint64Range(r)
	if err != nil {
		return err
	}
	for _, unRemoved := range wordlist.tempPaths {
		threadPool.Add(1)
		go func(r *utils.Uint64Range, p string) {
			defer threadPool.Done()
			newPath, err := removeWordsByRangeUint(r, p)
			if err != nil {
				panic(err)
			}
			channelNewPaths <- newPath
		}(uintRange, unRemoved)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewPaths)
	}()

	var newPaths []string
	for path := range channelNewPaths {
		newPaths = append(newPaths, path)
	}

	wordlist.tempPaths = newPaths
	return nil
}

// RemoveWordsByRangeUint removes all words from the current wordlist
// whose length in runes is outside the given Uint64Range.
// - If Min is set: words shorter than Min are skipped
// - If Max is set: words longer than Max are skipped
// - If both are nil: everything is kept
func removeWordsByRangeUint(r *utils.Uint64Range, inputPath string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("range must not be nil")
	}
	// Min must not be greater than Max
	if r.Min != nil && r.Max != nil && *r.Min > *r.Max {
		return "", fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %d Max: %d", *r.Min, *r.Max)
	}

	// generate a new temporary file path
	newTempPath, err := utils.GenerateRandomTempPath()
	if err != nil {
		return "", err
	}

	// create a new file for filtered content
	newFile, err := os.Create(newTempPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = newFile.Close()
	}()

	// open the current wordlist file
	currentFile, err := os.Open(inputPath)
	if err != nil {
		return "", err
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
			return "", err
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// remove old file and swap paths
	if err := os.Remove(inputPath); err != nil {
		return "", err
	}
	return newTempPath, nil
}

func (wordlist *EditWordlist) FlushFinishedWordlist() error {
	err := utils.MergeWordlists(wordlist.tempPaths, wordlist.outputPath)
	if err != nil {
		return err
	}
	err = utils.RemoveSplitWordlist(wordlist.tempPaths)
	if err != nil {
		return err
	}
	return nil
}
