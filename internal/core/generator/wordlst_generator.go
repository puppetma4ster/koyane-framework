package generator

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"

	"log"
	"os"
	"strings"
)

// ConcurrentGenerateMaskWordlist Parallelizes the generation of mask lists
// Partial masks are generated from a mask.
// Each partial mask writes to its own temporary file at the same time.
// Once all temporary files have been written, they are merged into one large file.
// The temporary files are then deleted.
//
// Parameters:
//   - maskArg: The mask that is to be generated in a word list
//   - outPath: The path where the created word list should be stored and the name of the list
//   - minLen: To determine the character length of the generation.
//     If left blank, only the mask length is generated.
//     If filled in, the desired character length up to the mask length is generated.
//
// Returns:
//   - error: if an error occurs during generation
func ConcurrentGenerateMaskWordlist(maskArg, outputPath string, minLen ...int) error {
	msk, err := NewMaskInterpreter(maskArg)
	if err != nil {
		fmt.Println("Error parsing mask:", err)
		return err
	}
	masks := msk.MaskSplitter()

	var files []*os.File

	// Prepare temp files. 1 per sub mask.
	for range masks {
		file, err := utils.GenerateNewTempFile("mask_gen*")
		if err != nil {
			return err
		}
		files = append(files, file)
	}

	if len(files) != len(masks) { // if there are not the same number of sub masks and temp files
		defer utils.CloseAllFiles(files)
		fmt.Errorf("a file was not prepared for each submask")
	}

	var threadPool sync.WaitGroup

	if len(minLen) == 0 { // if no minLen is given
		for i := 0; i < len(masks); i++ {
			threadPool.Add(1)
			go func(m *MaskInterpreter, f *os.File) {
				defer threadPool.Done()
				err := generateMaskWordlist(m, f)
				if err != nil {
					panic(err)
				}
			}(masks[i], files[i])
		}
	} else { // if minLen is given
		for i := 0; i < len(masks); i++ {
			threadPool.Add(1)
			go func(m *MaskInterpreter, f *os.File, l int) {
				defer threadPool.Done()
				err := generateMaskWordlist(m, f, l)
				if err != nil {
					panic(err)
				}
			}(masks[i], files[i], minLen[0])
		}
	}
	threadPool.Wait() // wait till every GoRoutine is finished

	absPath, err := utils.ListPath(outputPath) // The path is made absolute and the correct file extension is appended.
	if err != nil {
		return err
	}
	masterFile, err := os.Create(absPath) //create final file
	if err != nil {
		return err
	}

	// temp files are merged together
	for _, file := range files {
		// Set the cursor of the source file to the beginning
		if _, err := file.Seek(0, 0); err != nil {
			return err
		}
		// Copy everything from src to target
		if _, err := io.Copy(masterFile, file); err != nil {
			return err
		}

		// temp file cleanup...
		err := file.Close()
		if err != nil {
			return err
		}
		err = os.Remove(file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// generateMaskWordlist generates a wordlist based on a given mask.
// If no minimum length is specified, it generates only words of full mask length.
// If a minimum length is provided, it generates all combinations from that length up to the full mask length.
func generateMaskWordlist(msk *MaskInterpreter, outputFile *os.File, minLen ...int) error {
	// preparation of characters per position
	var fullSegments []string
	for _, seg := range msk.MaskSegments {
		fullSegments = append(fullSegments, seg.PermittedCharacters)
	}
	mskLength := len(fullSegments)

	var lengths []int
	if len(minLen) == 0 {
		lengths = []int{mskLength}
	} else {
		start := minLen[0]
		lengths = make([]int, 0, mskLength-start+1)
		for i := start; i <= mskLength; i++ {
			lengths = append(lengths, i)
		}
	}

	writer := bufio.NewWriterSize(outputFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	for _, length := range lengths {
		segments := fullSegments[:length]

		segmentSets := make([][]string, 0, length)

		for _, chars := range segments {
			set := make([]string, 0, len(chars))
			for _, ch := range chars {
				set = append(set, string(ch))
			}
			segmentSets = append(segmentSets, set)
		}

		err := productWriter(segmentSets, writer)
		if err != nil {
			log.Fatal("Write error:", err)
		}
	}

	return nil
}

// productWriter generates all combinations from segmentSets and writes each to the writer.
// This is a memory-efficient, iterative alternative to recursive generation.
func productWriter(segmentSets [][]string, writer *bufio.Writer) error {
	if len(segmentSets) == 0 {
		return nil
	}

	indexes := make([]int, len(segmentSets))
	sizes := make([]int, len(segmentSets))
	for i, set := range segmentSets {
		if len(set) == 0 {
			return nil // no characters to combine
		}
		sizes[i] = len(set)
	}

	buffer := make([]string, len(segmentSets))
	for {
		// Build current combination
		for i := range indexes {
			buffer[i] = segmentSets[i][indexes[i]]
		}
		word := strings.Join(buffer, "")
		if _, err := writer.WriteString(word + "\n"); err != nil {
			return err
		}

		// Increment indexes (like odometer)
		for i := len(indexes) - 1; i >= 0; i-- {
			indexes[i]++
			if indexes[i] < sizes[i] {
				break
			}
			if i == 0 {
				return nil // done
			}
			indexes[i] = 0
		}
	}
}
func CalculateMaskStorage(maskArg string, minLen ...uint16) (uint64, uint64, error) {
	mask, err := NewMaskInterpreter(maskArg)
	if err != nil {
		fmt.Println("Error parsing mask:", err)
		return 0, 0, err
	}
	var maxLen uint16 = uint16(len(mask.MaskSegments))
	var minSize uint16
	if len(minLen) == 0 {
		minSize = maxLen
	} else {
		minSize = minLen[0]
	}
	var totalCombinations uint64 = 0

	for length := minSize; length <= maxLen; length++ {
		var combinations uint64 = 1
		for _, seg := range mask.MaskSegments[:length] {
			combinations *= uint64(len([]rune(seg.PermittedCharacters)))
		}
		totalCombinations += combinations
	}
	var avgLineLength uint64 = 0
	if len(minLen) != 0 {
		avgLineLength = uint64(minLen[0])
	} else {
		avgLineLength = uint64(maxLen)
	}
	var estimatedBytes = totalCombinations * (avgLineLength + 1)
	return totalCombinations, estimatedBytes, nil

}
