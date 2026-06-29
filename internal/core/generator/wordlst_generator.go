package generator

import (
	"bufio"
	"fmt"
	"io"
	"runtime"
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
		err := utils.RemoveSplitWordlist(files)
		if err != nil {
			return err
		}
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

	writer := bufio.NewWriterSize(outputFile, 16*1024*1024) // 16 MiB buffer
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

// ExtractHashCatPotfile
// Is a method for extracting passwords from hashcat potfiles.
// It extracts the password by locating the last “:” character and extracting everything to the right of it.
func ExtractHashCatPotfile(inputHcListPath string, outputFile string) error {
	in, err := utils.ResolvePath(inputHcListPath) // resolving relative path to an absolute path
	if err != nil {
		return err
	}

	var absOutputPath string = outputFile
	if outputFile != "" { // Checking for stdout or writing to a file
		absOutputPath, err = utils.ListPath(outputFile)
		if err != nil {
			return err
		}
	}

	hcFile, err := os.Open(in) // open the hashcat potfile
	if err != nil {
		return err
	}
	defer hcFile.Close()

	extractPlain := func(hashAndVal string) string { // extracting passwords
		idx := strings.LastIndex(hashAndVal, ":")
		if idx == -1 {
			return ""
		}
		return hashAndVal[idx+1:]
	}
	// Create a buffer for stdout or write to a file
	var writer *bufio.Writer
	if absOutputPath == "" { // if no path is given stdout results
		writer = bufio.NewWriterSize(os.Stdout, 1024*1024) // 1 MB buffer
	} else { // write file
		openOutputFile, err := os.Create(absOutputPath)
		if err != nil {
			return err
		}
		writer = bufio.NewWriterSize(openOutputFile, 1024*1024) // 1 MB buffer
		defer openOutputFile.Close()
	}
	scanner := bufio.NewScanner(hcFile) // reading input file

	for scanner.Scan() { // write passwords
		_, err := writer.WriteString(extractPlain(scanner.Text()) + "\n")
		if err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

// ConcurrentPermutation is a high-performance streaming permutation engine for large wordlists.
//
// It reads a newline-separated input file containing words and generates all possible
// permutations of these words up to a given maximum length. The generation process is
// fully concurrent and designed for large-scale datasets (GB-sized inputs and very large outputs).
//
// Architecture:
//   - Input Stage: Words are loaded from disk using a streaming-safe reader.
//   - Job Distribution: Each word index is used as a starting point and distributed via a job channel.
//   - Worker Pool: Multiple goroutines (workers) execute permutation generation in parallel.
//   - DFS Generation: Each worker performs a depth-first search (DFS) permutation expansion
//     without duplicating state across workers except for local tracking.
//   - Output Stage: All generated permutations are streamed through a shared channel to a
//     single writer goroutine.
//
// Concurrency Model:
//   - Workers are stateless except for local recursion state.
//   - A buffered jobs channel distributes work across workers.
//   - A buffered output channel collects results.
//   - A single writer goroutine ensures safe sequential writes to disk without locking.
//
// Memory Behavior:
//   - No full permutation tree is stored in memory.
//   - Results are streamed immediately as they are generated.
//   - Memory usage is bounded by channel buffers and recursion depth.
//
// Scaling:
//   - If workers is set to 0, runtime.NumCPU() is used automatically.
//   - Performance scales with available CPU cores, but is ultimately limited by disk I/O.
//
// This function is intended for high-throughput wordlist generation and security research
// tooling where controlled combinatorial expansion is required.
func ConcurrentPermutation(inputFile string, outputFile string, minLen int, maxLen int, workers int) error {
	absInputPath, err := utils.ResolvePath(inputFile)
	if err != nil {
		return err
	}
	var absOutputPath string = ""
	if outputFile != "" {
		absOutputPath, err = utils.ListPath(outputFile)
		if err != nil {
			return err
		}
	}
	words, err := loadWords(absInputPath)
	if err != nil {
		return err
	}

	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	jobs := make(chan int, len(words))
	out := make(chan string, 100000)
	done := make(chan struct{})

	go func() {
		var writer *bufio.Writer
		var file *os.File

		if outputFile == "" {
			writer = bufio.NewWriterSize(os.Stdout, 16*1024*1024)
		} else {
			f, err := os.Create(absOutputPath)
			if err != nil {
				panic(err)
			}
			file = f
			writer = bufio.NewWriterSize(f, 16*1024*1024)
		}

		for line := range out {
			writer.WriteString(line)
			writer.WriteByte('\n')
		}

		writer.Flush()

		if file != nil {
			file.Close()
		}

		done <- struct{}{}
	}()

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go permutationWorker(i, words, minLen, maxLen, jobs, out, &wg)
	}

	go func() {
		for i := range words {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	<-done
	return nil
}

// loadWords reads a newline-separated word list from the given file path and returns
// all entries as a string slice. It streams the file line by line to support large files
// efficiently without loading the entire file into memory at once.
func loadWords(path string) ([]string, error) {
	f, err := os.Open(path) // path is already resolved in master function
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// support long lines
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}

// permuteDFS recursively builds permutations of the given word list starting from a prefix.
// It performs a depth-first search over all unused words, appending them to the current prefix
// as long as the resulting string does not exceed the maximum allowed length.
// Generated permutations are streamed to the output channel immediately to minimize memory usage.
func permuteDFS(
	words []string,
	prefix string,
	used []bool,
	minLen int,
	maxLen int,
	out chan<- string,
) {
	if len(prefix) > maxLen {
		return
	}

	for i := 0; i < len(words); i++ {
		if used[i] {
			continue
		}

		var sb strings.Builder
		sb.WriteString(prefix)
		sb.WriteString(words[i])
		next := sb.String()

		if len(next) > maxLen {
			continue
		}

		if len(next) >= minLen && len(next) <= maxLen {
			out <- next
		}

		used[i] = true
		permuteDFS(words, next, used, minLen, maxLen, out)
		used[i] = false
	}
}

// permutationWorker processes permutation jobs in parallel.
// Each worker receives a starting index from the job channel and generates all valid
// permutations starting from that word index up to the configured maximum length.
// Results are streamed to the shared output channel, while synchronization is handled
// via a WaitGroup to ensure proper shutdown of the pipeline.
func permutationWorker(
	id int,
	words []string,
	minLen int,
	maxLen int,
	jobs <-chan int,
	out chan<- string,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for start := range jobs {
		used := make([]bool, len(words))

		prefix := words[start]
		used[start] = true

		if len(prefix) >= minLen && len(prefix) <= maxLen {
			out <- prefix
		}

		permuteDFS(words, prefix, used, minLen, maxLen, out)
	}
}

func GenerateWordlist(maskArg []string, permutationArg []string, hcPotExtractArg []string, outputPathArg string,
	minLenArg int, maxLenArg int, quietArg bool, stdoutArg bool) error {

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
