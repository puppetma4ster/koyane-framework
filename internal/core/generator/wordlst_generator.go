package generator

import (
	"bufio"
	"fmt"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"

	"log"
	"os"
	"strings"
)

// GenerateMaskWordlist generates a wordlist based on a given mask.
// If no minimum length is specified, it generates only words of full mask length.
// If a minimum length is provided, it generates all combinations from that length up to the full mask length.
func GenerateMaskWordlist(maskArg string, outputPath string, minLen ...int) error {
	msk, err := NewMaskInterpreter(maskArg)
	if err != nil {
		fmt.Println("Error parsing mask:", err)
		return err
	}

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

	absPath, err := utils.ListPath(outputPath)
	if err != nil {
		fmt.Println("Path resolution error:", err)
		return err
	}
	outputFile, err := os.Create(absPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriterSize(outputFile, 1024*1024) // 1 MB buffer
	defer writer.Flush()

	for _, length := range lengths {
		segments := fullSegments[:length]

		segmentSets := make([][]string, 0, length) // Frisch anlegen

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
