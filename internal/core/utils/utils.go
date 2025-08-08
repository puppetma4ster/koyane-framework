package utils

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const LowerCaseCharacters string = "abcdefghijklmnopqrstuvwxyz" //?l
const UpperCaseCharacters string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" //?L
const LowerCaseVowels string = "aeiou"                          //?v
const UpperCaseVowels string = "AEIOU"                          //?V
const LowerCaseConsonants = "bcdfghjklmnpqrstvwxyz"             //?c
const UpperCaseConsonants = "BCDFGHJKLMNPQRSTVWXYZ"             //?C

const Digits string = "0123456789"                                         //?d
const SpecialCharactersMostUsed string = "!@#$%^&*()-_+=?"                 //?f
const SpecialCharactersPoints = ".,:;"                                     //?p
const SpecialCharactersBracelet = "()[]{}"                                 //?b
const SpecialCharacters string = "<>|^°!\"§$%&/()=?´{}[]\\¸`+~*#'-_.:,;@€" //?s

const ListSuffix string = ".klst"
const TempSuffix string = ".ktmp"

const tempDir = "/tmp/koyane_framework_tmp"
const chunkSize = 100000

func ExternalSort(inputPath, outputPath string) error {
	var tempFiles []string

	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return err
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	chunkIndex := 0

	for {
		var lines []string
		for len(lines) < chunkSize && scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if len(lines) == 0 {
			break
		}

		sort.Strings(lines) // Unicode

		chunkFile := filepath.Join(tempDir, fmt.Sprintf("chunk_%d%s", chunkIndex, TempSuffix))
		f, err := os.Create(chunkFile)
		if err != nil {
			return err
		}
		writer := bufio.NewWriter(f)
		for _, line := range lines {
			_, _ = writer.WriteString(line + "\n")
		}
		writer.Flush()
		f.Close()

		tempFiles = append(tempFiles, chunkFile)
		chunkIndex++
	}

	//  Merge
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	outWriter := bufio.NewWriter(outputFile)

	files := make([]*os.File, len(tempFiles))
	scanners := make([]*bufio.Scanner, len(tempFiles))
	currentLines := make([]string, len(tempFiles))

	for i, path := range tempFiles {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		files[i] = f
		scanners[i] = bufio.NewScanner(f)
		if scanners[i].Scan() {
			currentLines[i] = scanners[i].Text()
		} else {
			currentLines[i] = ""
		}
	}

	for {
		minIdx := -1
		for i, line := range currentLines {
			if line == "" {
				continue
			}
			if minIdx == -1 || strings.Compare(line, currentLines[minIdx]) < 0 {
				minIdx = i
			}
		}
		if minIdx == -1 {
			break
		}
		outWriter.WriteString(currentLines[minIdx] + "\n")
		if scanners[minIdx].Scan() {
			currentLines[minIdx] = scanners[minIdx].Text()
		} else {
			currentLines[minIdx] = ""
			files[minIdx].Close()
			os.Remove(tempFiles[minIdx])
		}
	}

	outWriter.Flush()
	return nil
}

// ResolvePath converts a given path string into an absolute, normalized file path.
// It supports:
// - Home directory expansion (e.g., "~" → "/home/user")
// - Relative paths (e.g., "./file.txt")
// - Parent directories (e.g., "../")
// - Absolute paths (unchanged)
func ResolvePath(input string) (string, error) {
	// Expand ~ to the user's home directory
	if strings.HasPrefix(input, "~/") || input == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		input = filepath.Join(home, strings.TrimPrefix(input, "~"))
	}

	// Convert to absolute path
	abs, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}

	// Clean the path (resolve "..", ".", etc.)
	clean := filepath.Clean(abs)

	return clean, nil
}

func ListPath(path string) (string, error) {
	absolutePath, err := ResolvePath(path)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(absolutePath)
	base := strings.TrimSuffix(path, ext)
	return base + ListSuffix, nil
}

func TempPath(path string) (string, error) {
	absolutePath, err := ResolvePath(path)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(absolutePath)
	base := strings.TrimSuffix(path, ext)
	return base + TempSuffix, nil
}

// CreateTempDir
// creates a directory for temporary files
// /*
func CreateTempDir() error {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("folder for temporary files could not be created: %w", err)
	}
	return nil
}

func HumanReadableBytes(bytes uint64) string {
	const unit = 1024
	unitsArr := [6]string{"B", "KB", "MB", "GB", "TB", "PB"}

	f := float64(bytes)
	for i, u := range unitsArr {

		if i == len(unitsArr)-1 || f < float64(unit) {
			return fmt.Sprintf("%.2f %s", f, u)
		}
		f /= float64(unit)
	}

	return fmt.Sprintf("%.2f PB", f)
}

func GenerateRandomTempPath() (string, error) {
	const randRange = 10_000_000
	const maxAttempts = 10000
	src := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(src)

	for i := 0; i < maxAttempts; i++ {
		tempNr := randomGen.Intn(randRange)
		tempPath := tempDir + "/tempFile_" + strconv.Itoa(tempNr) + TempSuffix
		_, err := os.Stat(tempPath)
		if os.IsNotExist(err) {
			return tempPath, nil
		}
	}
	return "", fmt.Errorf("couldn't create a temporary file")
}

func CopyFileToTemp(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}
	return outputFile.Sync()
}

func NotInSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return false
		}
	}
	return true
}
