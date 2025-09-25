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

	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v3"
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
const AnalyzedWlSaveDir = "~/.koyane_framework_saves"
const chunkSize = 100000

func ExternalSort(inputFile *os.File, outputFile *os.File) error {
	var tempFiles []string

	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(inputFile)
	chunkIndex := 0

	newTempFolder, err := os.MkdirTemp(tempDir, "sort")
	if err != nil {
		return err
	}

	// Split   and  sort
	for {
		var lines []string
		for len(lines) < chunkSize && scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if len(lines) == 0 {
			break
		}

		sort.Strings(lines)
		chunkFile := filepath.Join(newTempFolder, fmt.Sprintf("chunk_%d%s", chunkIndex, TempSuffix))
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
		}
	}

	outWriter.Flush()
	// removeing chunks
	err = os.RemoveAll(newTempFolder)
	if err != nil {
		return err
	}

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

// ListPath generates the final list file path based on a temporary file path.
//
// It performs the following steps:
// 1. Resolves the given path to an absolute path using ResolvePath().
// 2. Extracts the file extension from the absolute path.
// 3. Removes the extension from the original path to get the base name.
// 4. Appends the predefined ListSuffix to the base name to create the final file path.
// 5. Returns the new file path and any error encountered during path resolution.
func ListPath(path string) (string, error) {
	absolutePath, err := ResolvePath(path)
	if err != nil {
		return "", err // Return empty string and error if path cannot be resolved
	}

	ext := filepath.Ext(absolutePath)     // Extract file extension, e.g., ".tmp"
	base := strings.TrimSuffix(path, ext) // Remove extension from original path
	return base + ListSuffix, nil         // Append ListSuffix and return final path
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
	absoluteUserPath, err := ResolvePath(AnalyzedWlSaveDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absoluteUserPath, 0755); err != nil {
		return fmt.Errorf("folder for analyzedWlSaveDir could not be created: %w", err)
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

// Deprecated: GenerateRandomTempPath
// use GenerateNewTempFile instead
// GenerateRandomTempPath proved to be extremely unstable, especially when there were many tempfiles in the directory.
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

func GenerateNewTempFile(name string) (*os.File, error) {
	f, err := os.CreateTemp(tempDir, name+TempSuffix)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func splitRange(arg string) (string, string, error) {
	if arg == "" {
		return "", "", fmt.Errorf("no arguments specified")
	}
	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid range %q (expected min:max)", arg)
	}
	a := strings.TrimSpace(parts[0])
	b := strings.TrimSpace(parts[1])
	return a, b, nil
}

// ---------- UINT64 ----------

type Uint64Range struct {
	Min *uint64
	Max *uint64
}

func NewUint64Range(arg string) (*Uint64Range, error) {
	a, b, err := splitRange(arg)
	if err != nil {
		return nil, err
	}

	var minPtr, maxPtr *uint64
	var minNumber, maxNumber uint64
	var haveMin, haveMax bool

	// validate min
	if a == "" || a == "*" {
		// open
	} else if v, err := strconv.ParseUint(a, 10, 64); err != nil {
		return nil, fmt.Errorf("the Minimal Range is not a number or \"*\": %s", a)
	} else {
		minNumber = v
		minPtr = &minNumber
		haveMin = true
	}

	// validate max
	if b == "" || b == "*" {
		// open
	} else if v, err := strconv.ParseUint(b, 10, 64); err != nil {
		return nil, fmt.Errorf("the Minimal Range is not a number or \"*\": %s", b) // (genau deine Message für max)
	} else {
		maxNumber = v
		maxPtr = &maxNumber
		haveMax = true
	}

	// min <= max
	if haveMin && haveMax && minNumber > maxNumber {
		return nil, fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %d Max: %d", minNumber, maxNumber)
	}

	return &Uint64Range{Min: minPtr, Max: maxPtr}, nil
}

// ---------- FLOAT64 ----------

type Float64Range struct {
	Min *float64
	Max *float64
}

func NewFloat64Range(arg string) (*Float64Range, error) {
	a, b, err := splitRange(arg)
	if err != nil {
		return nil, err
	}

	var minPtr, maxPtr *float64
	var minNumber, maxNumber float64
	var haveMin, haveMax bool

	// validate min
	if a == "" || a == "*" {
		// open
	} else if v, err := strconv.ParseFloat(a, 64); err != nil {
		return nil, fmt.Errorf("the Minimal Range is not a number or \"*\": %s", a)
	} else {
		minNumber = v
		minPtr = &minNumber
		haveMin = true
	}

	// validate max
	if b == "" || b == "*" {
		// open
	} else if v, err := strconv.ParseFloat(b, 64); err != nil {
		return nil, fmt.Errorf("the Minimal Range is not a number or \"*\": %s", b) // gleiche Message wie gewünscht
	} else {
		maxNumber = v
		maxPtr = &maxNumber
		haveMax = true
	}

	// min <= max
	if haveMin && haveMax && minNumber > maxNumber {
		return nil, fmt.Errorf("the minimum variable must not be greater than the maximum! Min: %g Max: %g", minNumber, maxNumber)
	}

	return &Float64Range{Min: minPtr, Max: maxPtr}, nil
}

type Config struct {
	General struct {
		DefaultWordlistPath  string `yaml:"default_wordlist_file_location"`
		UserAgent            string `yaml:"user_agent"`
		DatabasePath         string `yaml:"wordlist_db_path"`
		MultiThreadFileLines uint32 `yaml:"multithread_file_lines"`
	} `yaml:"general"`
	Ui struct {
		Language            string `yaml:"language"`
		disableOutputColors bool   `yaml:"disable_output_colors"`
	} `yaml:"ui"`
	Editor struct {
		MaxWordLen uint16 `yaml:"max_word_length"`
	} `yaml:"editor"`
	Analyzer struct {
		SaveAllAnalyzed bool `yaml:"save_all_analyzed_values"`
	} `yaml:"analyzer"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Settings struct {
	binaryPaths struct {
		CfgPath string `yaml:"cfgString"`
		DbPath  string `yaml:"bdPath"`
	}
}

func LoadSettings(path string) (*Settings, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var set Settings
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&set); err != nil {
		return nil, err
	}
	return &set, nil
}

func SplitWordlist(inputPath string) ([]*os.File, error) {
	// Load config to determine how many lines each chunk file should contain
	cfgPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	maxLines := cfg.General.MultiThreadFileLines

	// Resolve and open the input file
	absInputPath, err := ResolvePath(inputPath)
	if err != nil {
		return nil, err
	}

	mainFile, err := os.Open(absInputPath)
	if err != nil {
		return nil, err
	}
	defer mainFile.Close()

	var (
		tempFile  *os.File
		writer    *bufio.Writer
		retFiles  []*os.File
		lineIndex uint32
	)

	// startNew creates a new temporary file and prepares a buffered writer
	startNew := func() error {
		var err error
		tempFile, err = GenerateNewTempFile("split_chunk*")
		if err != nil {
			return err
		}
		writer = bufio.NewWriterSize(tempFile, 128*1024)
		lineIndex = 0
		return nil
	}

	if err := startNew(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(mainFile)

	for scanner.Scan() {
		// Write the current line into the current chunk
		if _, err := writer.WriteString(scanner.Text() + "\n"); err != nil {
			return nil, err
		}
		lineIndex++

		// If the chunk reached the maximum number of lines → finalize it
		if lineIndex >= maxLines {
			if err := writer.Flush(); err != nil {
				return nil, err
			}
			if err := tempFile.Close(); err != nil {
				return nil, err
			}

			// Reopen in read mode before returning it
			f, err := os.Open(tempFile.Name())
			if err != nil {
				return nil, err
			}
			retFiles = append(retFiles, f)

			if err := startNew(); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Handle the last file (if it has content)
	if lineIndex > 0 {
		if err := writer.Flush(); err != nil {
			return nil, err
		}
		if err := tempFile.Close(); err != nil {
			return nil, err
		}

		f, err := os.Open(tempFile.Name())
		if err != nil {
			return nil, err
		}
		retFiles = append(retFiles, f)
	} else {
		// If the last chunk had no lines, remove it
		os.Remove(tempFile.Name())
	}

	return retFiles, nil
}

func RemoveSplitWordlist(files []*os.File) error {
	for _, file := range files {
		err := file.Close() // Close all files
		if err != nil {
			return err
		}
		err = os.Remove(file.Name()) // delete all files
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseAllFiles(f []*os.File) {
	for _, file := range f {
		file.Close() // funktioniert, weil file vom Typ *os.File ist
	}
}

// MergeWordlists merges multiple open *os.File into a single output file.
// inputFiles: slice of already opened *os.File
// outputPath: path to the final merged file
func MergeWordlists(inputFiles []*os.File, outputPath string) error {
	// Resolve or create the output file
	newPath, err := ListPath(outputPath)
	if err != nil {
		return err
	}

	outFile, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	for _, f := range inputFiles {
		// Reset file offset to start in case it's not at 0
		if _, err := f.Seek(0, 0); err != nil {
			return fmt.Errorf("failed to seek input file %s: %w", f.Name(), err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if _, err := writer.WriteString(line + "\n"); err != nil {
				return fmt.Errorf("failed to write line: %w", err)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading file %s: %w", f.Name(), err)
		}
	}

	return nil
}

// formatIntWithCommas returns an int64 as string with commas as thousands separators.
func formatIntWithCommas(n int64) string {
	// Negative Numbers...
	negative := n < 0
	if negative {
		n = -n
	}

	// int64 -> string
	s := strconv.FormatInt(n, 10)

	// right to left
	out := ""
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		out = string(s[i]) + out
		count++
		if count%3 == 0 && i != 0 {
			out = "," + out
		}
	}

	if negative {
		return "-" + out
	}
	return out
}

// FormatUint64WithCommas returns a uint64 as string with commas as thousands separators.
func FormatUint64WithCommas(n uint64) string {
	s := strconv.FormatUint(n, 10) // uint64 -> string

	out := ""
	count := 0
	for i := len(s) - 1; i >= 0; i-- {
		out = string(s[i]) + out
		count++
		if count%3 == 0 && i != 0 {
			out = "," + out
		}
	}

	return out
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

// TruncateWithTail returns a transformer for go-pretty table columns
// that truncates long strings and appends a custom tail
func TruncateWithTail(tail string, maxLen int) func(interface{}) string {
	return func(val interface{}) string {
		str := fmt.Sprint(val)
		if len(str) <= maxLen {
			return str
		}
		cut := maxLen - len(tail)
		if cut < 0 {
			return tail
		}
		return str[:cut] + tail
	}
}

// PrintDotted Formats output statistics
func PrintDotted(label string, value interface{}, unit ...string) string {
	totalWidth := 25
	dots := totalWidth - len(label)
	if dots < 0 {
		dots = 0
	}

	var valueStr string
	switch v := value.(type) {
	case float32:
		valueStr = fmt.Sprintf("%.2f", v)
	case float64:
		valueStr = fmt.Sprintf("%.2f", v)
	case int:
		valueStr = fmt.Sprintln(humanize.Comma(int64(v)))
	case uint64:
		valueStr = fmt.Sprintln(humanize.Comma(int64(v)))
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	unitStr := ""
	if len(unit) > 0 {
		unitStr = unit[0]
	}

	return fmt.Sprintf("%s%s: %s%s\n", label, strings.Repeat(".", dots), valueStr, unitStr)
}

// PrintRuneMapColumns formats  rune slice for analyze output
func PrintRuneMapColumns(m map[rune]uint64, columns int) string {
	var builder strings.Builder
	totalWidth := 12

	i := 0
	for k, v := range m {
		keyStr := string(k)
		dots := totalWidth - len(keyStr) - len(fmt.Sprintf(": %d", v))
		if dots < 0 {
			dots = 0
		}

		fmtStr := fmt.Sprintf("%s%s: %d    ", keyStr, strings.Repeat(".", dots), v)
		builder.WriteString(fmtStr)

		i++
		if i%columns == 0 {
			builder.WriteString("\n")
		}
	}
	if i%columns != 0 {
		builder.WriteString("\n")
	}

	return builder.String()
}

// fileExists checks if an file exists
//
// Parameters:
//   - path: the path to the file to be tested
//
// Returns:
//   - bool: true if file exists and false if it doesn't
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		// file exists
		return true
	}
	if os.IsNotExist(err) {
		// file dont exists
		return false
	}
	return false
}

// GetConfigPath determines the path to the program's configuration file.
// It checks possible locations in the following order:
// 1. Default system path: "/etc/koyane-framework/config.yaml".
// 2. User-specific path: "$HOME/.config/koyane-framework/settings.yaml", if this file exists.
// 3. Environment variable "KOYANE_CONFIG": if set, it overrides the previous paths.
// Returns:
// - string: the resolved path to the configuration file.
// - error: an error if, for example, the user's home directory cannot be determined.
func GetConfigPath() (string, error) {
	var configPath string = "/etc/koyane-framework/config.yaml"
	var changed bool = false

	// Check the user's home settings YAML
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	settingsPath := filepath.Join(home, ".config/koyane-framework/settings.yaml")
	if fileExists(settingsPath) {
		set, err := LoadSettings(settingsPath)
		if err != nil {
			return "", err
		}
		configPath = set.binaryPaths.DbPath
		changed = true
	}

	// Check if an environment variable is set
	if envPath := os.Getenv("KOYANE_CONFIG"); envPath != "" {
		configPath = envPath
		changed = true
	}
	if !changed {
		if !fileExists(configPath) {
			return "", fmt.Errorf("there is no config path found as environment variable, "+
				"user home and default path: %s", configPath)
		}
	}
	return configPath, nil
}

// GetDatabasePath determines the path to the program's database file.
func GetDatabasePath() (string, error) {
	var databasePath string = "/var/lib/koyane-framework/wordLists.db"
	var changed bool = false

	// Check the user's home settings YAML
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	settingsPath := filepath.Join(home, ".config/koyane-framework/settings.yaml")
	if fileExists(settingsPath) {
		set, err := LoadSettings(settingsPath)
		if err != nil {
			return "", err
		}
		databasePath = set.binaryPaths.DbPath
		changed = true
	}

	// Check if an environment variable is set
	if envPath := os.Getenv("KOYANE_DATABASE"); envPath != "" {
		databasePath = envPath
		changed = true
	}
	if !changed {
		if !fileExists(databasePath) {
			return "", fmt.Errorf("there is no config path found as environment variable, "+
				"user home and default path: %s", databasePath)
		}
	}
	return databasePath, err
}
