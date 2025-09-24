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
	tempFiles  []*os.File
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
		tempFiles:  tempList,
	}, nil
}

// ConcurrentSortWordlist parallelizes the sorting process
// All temp lists are sorted and the new string is captured in channelNewPaths.
// channelNewPaths is then passed to the struct as a new array and the isSorted switch is set to true.
func (wordlist *EditWordlist) ConcurrentSortWordlist() {
	var threadPool sync.WaitGroup
	channelNewFiles := make(chan *os.File) // channel to catch new paths
	for _, unSortPath := range wordlist.tempFiles {
		threadPool.Add(1)
		go func(p *os.File) {
			defer threadPool.Done()
			result, err := sortWordlist(p)
			if err != nil {
				panic(err)
			}
			channelNewFiles <- result
		}(unSortPath)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewFiles)
	}()

	var newFiles []*os.File
	for path := range channelNewFiles {
		newFiles = append(newFiles, path)
	}
	wordlist.isSorted = true
	wordlist.tempFiles = newFiles
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
func sortWordlist(inputPath *os.File) (*os.File, error) {
	newTempPath, err := utils.GenerateNewTempFile("Edit_Sort*")
	if err != nil {
		return nil, err
	}
	err = utils.ExternalSort(inputPath, newTempPath)
	if err != nil {
		return nil, err
	}
	if err := inputPath.Close(); err != nil {
		return nil, err
	}
	err = inputPath.Close() // close and delete old path
	if err != nil {
		return nil, err
	}
	err = os.Remove(inputPath.Name())
	if err != nil {
		return nil, err
	}
	return newTempPath, nil
}

func (wordlist *EditWordlist) ConcurrentRemoveWordsWithMask(msk string) error {
	var threadPool sync.WaitGroup
	channelNewFiles := make(chan *os.File) // channel to catch new paths

	mask, err := generator.NewMaskInterpreter(msk)
	if err != nil {
		return err
	}

	for _, unRemoved := range wordlist.tempFiles {
		threadPool.Add(1)
		go func(m *generator.MaskInterpreter, f *os.File) {
			defer threadPool.Done()
			newFile, err := removeWordsWithMask(m, f)
			if err != nil {
				panic(err)
			}
			channelNewFiles <- newFile
		}(mask, unRemoved)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewFiles)
	}()

	var newPaths []*os.File
	for path := range channelNewFiles {
		newPaths = append(newPaths, path)
	}

	wordlist.tempFiles = newPaths
	return nil
}
func removeWordsWithMask(mask *generator.MaskInterpreter, inputPath *os.File) (*os.File, error) {
	newFile, err := utils.GenerateNewTempFile("Remove_Mask*") // create new Wordlist
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriterSize(newFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	scanner := bufio.NewScanner(inputPath)
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

	err = inputPath.Close()
	if err != nil {
		return nil, err
	}
	err = os.Remove(inputPath.Name())
	if err != nil {
		return nil, err
	}
	return newFile, nil
}

func (wordlist *EditWordlist) ConcurrentRemoveWordsByRange(r string) error {
	var threadPool sync.WaitGroup
	channelNewFiles := make(chan *os.File) // channel to catch new paths

	uintRange, err := utils.NewUint64Range(r)
	if err != nil {
		return err
	}
	for _, unRemoved := range wordlist.tempFiles {
		threadPool.Add(1)
		go func(r *utils.Uint64Range, f *os.File) {
			defer threadPool.Done()
			newFile, err := removeWordsByRangeUint(r, f)
			if err != nil {
				panic(err)
			}
			channelNewFiles <- newFile
		}(uintRange, unRemoved)
	}
	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelNewFiles)
	}()

	var newFiles []*os.File
	for file := range channelNewFiles {
		newFiles = append(newFiles, file)
	}

	wordlist.tempFiles = newFiles
	return nil
}

// RemoveWordsByRangeUint removes all words from the current wordlist
// whose length in runes is outside the given Uint64Range.
// - If Min is set: words shorter than Min are skipped
// - If Max is set: words longer than Max are skipped
// - If both are nil: everything is kept
func removeWordsByRangeUint(r *utils.Uint64Range, inputFile *os.File) (*os.File, error) {
	if r == nil {
		return nil, fmt.Errorf("range must not be nil")
	}
	// Min must not be greater than Max
	if r.Min != nil && r.Max != nil && *r.Min > *r.Max {
		return nil, fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %d Max: %d", *r.Min, *r.Max)
	}

	// generate a new temporary file path
	newTempFile, err := utils.GenerateNewTempFile("Remove_Range*")
	if err != nil {
		return nil, err
	}

	// buffered writer for performance (1 MB buffer)
	writer := bufio.NewWriterSize(newTempFile, 1024*1024)
	defer writer.Flush()

	// scanner with increased buffer size (default 64 KB → bumped to 1 MB)
	scanner := bufio.NewScanner(inputFile)
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
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// remove and close old file and swap paths
	err = inputFile.Close()
	if err != nil {
		return nil, err
	}
	if err := os.Remove(inputFile.Name()); err != nil {
		return nil, err
	}
	return newTempFile, nil
}

func (wordlist *EditWordlist) FlushFinishedWordlist() error {
	err := utils.MergeWordlists(wordlist.tempFiles, wordlist.outputPath)
	if err != nil {
		return err
	}
	err = utils.RemoveSplitWordlist(wordlist.tempFiles)
	if err != nil {
		return err
	}
	return nil
}
