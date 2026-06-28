package editor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"sort"
	"sync"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/puppetma4ster/koyane-framework/internal/output"
)

// Chunk struct for managing individual chunks
// Index - chunk number to ensure correct sequence
// Lines - the words in the chunk
type Chunk struct {
	Index uint64
	Lines []string
}

type rangeFilter struct {
	rangeF *utils.Uint64Range
}

func (r rangeFilter) keepRange(length uint64) bool {

	return (r.rangeF.Min == nil || length >= *r.rangeF.Min) &&
		(r.rangeF.Max == nil || length <= *r.rangeF.Max)
}

type maskFilter struct {
	maskF *generator.MaskInterpreter
}

func (m maskFilter) keepMask(word string) bool {
	return m.maskF.MatchesWord(word)
}

type regExFilter struct {
	regEx *regexp.Regexp
}

func (x regExFilter) keepRegEx(word string) bool {
	return x.regEx.MatchString(word)
}

// tempRun represents a single sorted temporary file
// generated during the external sort process.
//
// Path points to the temporary file on disk.
// Err is used to report worker errors.
type tempRun struct {
	Path string
	Err  error
}

// readChunksWorker reads the input stream line by line,
// splits the input into chunks and sends them through
// the output channel.
//
// input          - input stream (stdin or file)
// chunkSize      - maximum number of lines per chunk
// byteLineLimit  - maximum allowed line size in bytes (0 = unlimited)
// out            - output channel for processed chunks
func readChunksWorker(input io.Reader, chunkByteSize int, byteLineLimit int, out chan<- Chunk) error {

	// Create a buffered reader
	reader := bufio.NewReader(input)

	estimatedAvgLineSize := 30 // ~ 30 bytes
	estimatedLines := chunkByteSize / estimatedAvgLineSize
	// Create the first chunk
	chunk := Chunk{
		Index: 0,
		Lines: make([]string, 0, estimatedLines),
	}

	currentBytes := 0 // chunk byte counter
	for {

		// Read one line including the trailing '\n'
		line, err := reader.ReadBytes('\n')

		// Return unexpected errors
		if err != nil && err != io.EOF {
			return err
		}
		lineSize := len(line)
		// Process the line if any bytes were read
		if lineSize > 0 {

			// Remove trailing newline
			if line[lineSize-1] == '\n' {
				line = line[:lineSize-1]
				lineSize--
			}

			// Remove Windows carriage return
			if lineSize > 0 &&
				line[lineSize-1] == '\r' {

				line = line[:lineSize-1]
				lineSize--
			}

			// Skip lines that exceed the configured limit
			if byteLineLimit > 0 && lineSize > byteLineLimit {

				if err == io.EOF {
					break
				}
				// TODO: implement warning log/output print
				continue
			}

			// Would this line exceed the configured chunk size?
			if currentBytes > 0 &&
				currentBytes+lineSize > chunkByteSize {

				out <- chunk

				chunk = Chunk{
					Index: chunk.Index + 1,
					Lines: make([]string, 0, estimatedLines),
				}

				currentBytes = 0
			}

			// Append the line
			chunk.Lines = append(
				chunk.Lines,
				string(line),
			)

			currentBytes += lineSize
		}

		// End of file reached
		if err == io.EOF {
			break
		}
	}

	// Flush the last partially filled chunk
	if len(chunk.Lines) > 0 {
		out <- chunk
	}

	// Close the output channel
	close(out)

	return nil
}

// startReader creates the first stage of the pipeline.
// It reads the input stream, splits it into chunks
// and returns the output channel.
func startReader(
	input io.Reader,
	chunkByteSize int,
	byteLineLimit int,
) <-chan Chunk {

	out := make(chan Chunk)

	go func() {

		err := readChunksWorker(
			input,
			chunkByteSize,
			byteLineLimit,
			out,
		)

		if err != nil {
			panic(err)
		}

	}()

	return out
}

// writeChunks writes each processed chunk to stdout or to a file
// output - where to write (open file or stdout)
// in - channel containing the finished chunks
func writeChunks(output io.Writer, in <-chan Chunk) error {

	writer := bufio.NewWriter(output) // Create the Writer

	for chunk := range in { // loop through all chunks

		for _, line := range chunk.Lines { // loop through all chunk lines

			_, err := writer.WriteString( // Write a line with a line break
				line + "\n",
			)

			if err != nil {
				return err
			}
		}
	}
	if err := writer.Flush(); err != nil { // if errors occurred during the writing process
		return err
	}
	return nil
}

// removeRangeWorker specifies how the range is filtered
// r - range to filter
// in -the chunk being filtered
// out - the finished filtered chunk
func filterRangeWorker(filterRanges []*rangeFilter, invFilterRanges []*rangeFilter, in <-chan Chunk, out chan<- Chunk) {
	for chunk := range in { // loop through all the given chunks

		filtered := Chunk{ // create a new empty chunk
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines {
			l := uint64(utf8.RuneCountInString(line))
			fmt.Printf("line=%q len=%d\n", line, l)
			// Keep filter (-r)
			keep := len(filterRanges) == 0

			for _, r := range filterRanges {
				fmt.Printf(
					"match=%v\n",
					r.keepRange(l),
				)
				if r.keepRange(l) {
					keep = true
					break
				}
			}

			if !keep {
				continue
			}

			// Remove filter (-R)
			remove := false

			for _, r := range invFilterRanges {
				if r.keepRange(l) {
					remove = true
					break
				}
			}

			if remove {
				continue
			}

			filtered.Lines = append(
				filtered.Lines,
				line,
			)
		}

		out <- filtered // Write the filtered chunk to the output channel
	}
}

// startFilterRange manages the go routines to be executed for the range filter
// r - the processed range arg with min and max
// workerCount - how many threads (workers) should be created
// in - input stream of the chunks to be filtered
func startFilterRange(
	r []*rangeFilter,
	rInv []*rangeFilter,
	workerCount int,
	in <-chan Chunk,
) <-chan Chunk {

	out := make(chan Chunk) // implement new output channel

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines

		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterRangeWorker( // worker function
				r,
				rInv,
				in,
				out,
			)
		}()
	}

	go func() { // wait for threads
		wg.Wait()
		close(out)
	}()

	return out // return new channel
}

// filterMaskWorker processes chunks and applies mask-based filtering.
//
// Filtering is performed in two stages:
//
//  1. Keep filters (-m)
//     - If at least one keep filter exists, a line must match
//     at least one of them to remain in the output.
//     - If no keep filters exist, all lines are accepted.
//
//  2. Remove filters (-M)
//     - Lines that match any remove filter are discarded.
//
// This results in the following logic:
//
//	(Keep1 OR Keep2 OR ...)
//	AND NOT
//	(Remove1 OR Remove2 OR ...)
//
// msks     - keep filters
// invMasks - remove filters
// in       - input chunk stream
// out      - output chunk stream
func filterMaskWorker(
	msks []*maskFilter,
	invMasks []*maskFilter,
	in <-chan Chunk,
	out chan<- Chunk,
) {

	// Process all incoming chunks
	for chunk := range in {

		// Create a new chunk for filtered results
		filtered := Chunk{
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		// Process every line in the current chunk
		for _, line := range chunk.Lines {

			// -----------------------------
			// Keep filters (-m)
			// -----------------------------

			// If no keep filters exist,
			// every line is accepted by default.
			keep := len(msks) == 0

			// Check whether the line matches
			// at least one keep filter.
			for _, m := range msks {

				if m.keepMask(line) {

					// One match is sufficient.
					keep = true
					break
				}
			}

			// Skip lines that do not match
			// any keep filter.
			if !keep {
				continue
			}

			// -----------------------------
			// Remove filters (-M)
			// -----------------------------

			// Tracks whether the line should
			// be removed from the output.
			remove := false

			// Check whether the line matches
			// any remove filter.
			for _, m := range invMasks {

				if m.keepMask(line) {

					// One remove match is enough
					// to discard the line.
					remove = true
					break
				}
			}

			// Skip removed lines.
			if remove {
				continue
			}

			// The line passed all filters
			// and can be written to the
			// filtered chunk.
			filtered.Lines = append(
				filtered.Lines,
				line,
			)
		}

		// Send the filtered chunk to the
		// next stage of the pipeline.
		out <- filtered
	}
}

// startFilterMask handles the go routines for the mask filter
// msk - processed mask arg for filtering
// workerCount - how many threads (workers) should be created
// in - the chunks to be filtered
func startFilterMask(msks []*maskFilter, mskInv []*maskFilter, workerCount int, in <-chan Chunk) <-chan Chunk {
	out := make(chan Chunk) // implement new output channel

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines
		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterMaskWorker( // worker function
				msks,
				mskInv,
				in,
				out,
			)
		}()
	}
	go func() { // wait for threads
		wg.Wait()
		close(out)
	}()

	return out // return new channel
}

func filterSubstringsWorker() {

}

func startFilterSubstrings() {

}

// regexFilterWorker processes chunks and applies regex-based filtering.
//
// Filtering is performed in two stages:
//
//  1. Keep filters (-x)
//     - If at least one keep regex exists, a line must match
//     at least one of them to remain in the output.
//     - If no keep regex filters exist, all lines are accepted.
//
//  2. Remove filters (-X)
//     - Lines that match any remove regex are discarded.
//
// This results in the following logic:
//
//	(Keep1 OR Keep2 OR ...)
//	AND NOT
//	(Remove1 OR Remove2 OR ...)
//
// re     - keep regex filters
// reInve - remove regex filters
// in     - input chunk stream
// out    - output chunk stream
func regexFilterWorker(
	re []*regExFilter,
	reInve []*regExFilter,
	in <-chan Chunk,
	out chan<- Chunk,
) {

	// Process all incoming chunks
	for chunk := range in {

		// Create a new chunk for filtered results
		filtered := Chunk{
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		// Process every line in the current chunk
		for _, line := range chunk.Lines {

			// -----------------------------
			// Keep filters (-x)
			// -----------------------------

			// If no keep regex filters exist,
			// every line is accepted by default.
			keep := len(re) == 0

			// Check whether the line matches
			// at least one keep regex.
			for _, x := range re {

				if x.keepRegEx(line) {

					// One match is sufficient.
					keep = true
					break
				}
			}

			// Skip lines that do not match
			// any keep regex.
			if !keep {
				continue
			}

			// -----------------------------
			// Remove filters (-X)
			// -----------------------------

			// Tracks whether the line should
			// be removed from the output.
			remove := false

			// Check whether the line matches
			// any remove regex.
			for _, x := range reInve {

				if x.keepRegEx(line) {

					// One remove match is enough
					// to discard the line.
					remove = true
					break
				}
			}

			// Skip removed lines.
			if remove {
				continue
			}

			// The line passed all filters
			// and can be written to the
			// filtered chunk.
			filtered.Lines = append(
				filtered.Lines,
				line,
			)
		}

		// Send the filtered chunk to the
		// next stage of the pipeline.
		out <- filtered
	}
}
func startRegexFilter(
	re []*regExFilter,
	reInv []*regExFilter,
	workerCount int,
	in <-chan Chunk,
) <-chan Chunk {

	out := make(chan Chunk)

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			regexFilterWorker(
				re,
				reInv,
				in,
				out,
			)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// externalSortWorker receives chunks from the input channel,
// sorts each chunk in memory and writes the sorted result
// into a temporary file.
//
// Each temporary file represents one sorted run that can
// later be merged using a k-way merge.
//
// Input:
//
//	Chunk
//	↓
//	sort.Strings()
//	↓
//	temp file
//
// Output:
//
//	tempChunk{Path: "..."}
func externalSortWorker(
	in <-chan Chunk,
	out chan<- tempRun,
) {

	// Process all incoming chunks.
	for chunk := range in {

		// Sort the chunk in memory.
		sort.Strings(chunk.Lines)
		// Create a new temporary file that will hold
		// the sorted chunk.
		tmpFile, err := os.CreateTemp(
			"",
			"kyf-sort-*",
		)

		if err != nil {
			out <- tempRun{Err: err}
			return
		}

		// Buffered writer improves write performance.
		writer := bufio.NewWriter(tmpFile)

		// Write all sorted lines into the temporary file.
		for _, line := range chunk.Lines {

			_, err := writer.WriteString(
				line + "\n",
			)

			if err != nil {

				tmpFile.Close()

				out <- tempRun{Err: err}
				return
			}
		}

		// Flush remaining buffered data to disk.
		if err := writer.Flush(); err != nil {

			tmpFile.Close()

			out <- tempRun{Err: err}
			return
		}

		// Close the file so it can later be reopened
		// by the merge phase.
		tmpFile.Close()

		// Return the path of the generated sorted run.
		out <- tempRun{
			Path: tmpFile.Name(),
		}
	}
}

// startExternalSort manages the complete external sort pipeline.
//
// Workflow:
//
//	input chunks
//	      ↓
//	multiple sort workers
//	      ↓
//	sorted temporary files
//	      ↓
//	k-way merge
//	      ↓
//	output chunks
//
// workerCount controls how many chunks may be sorted
// concurrently.
//
// chunkSize specifies the size of the merged output chunks.
func startExternalSort(
	workerCount int,
	chunkSize uint32,
	removeDuplicates bool,
	in <-chan Chunk,
) <-chan Chunk {

	// Output channel containing the final merged chunks.
	out := make(chan Chunk)

	go func() {

		// Channel used by workers to report generated
		// temporary files.
		tempFiles := make(chan tempRun)

		// WaitGroup tracks all sorting workers.
		var wg sync.WaitGroup

		// Start the sorting worker pool.
		for i := 0; i < workerCount; i++ {

			wg.Add(1)

			go func() {

				defer wg.Done()

				externalSortWorker(
					in,
					tempFiles,
				)
			}()
		}

		// Close tempFiles once all workers have finished.
		go func() {

			wg.Wait()

			close(tempFiles)
		}()

		// Collect all generated temporary file paths.
		var paths []string

		for result := range tempFiles {

			// Abort immediately if a worker reports an error.
			if result.Err != nil {
				panic(result.Err)
			}

			paths = append(
				paths,
				result.Path,
			)
		}

		// Perform the final k-way merge across all
		// sorted temporary files.
		err := mergeSortedChunks(
			paths,
			int(chunkSize),
			out,
			removeDuplicates,
		)

		if err != nil {
			panic(err)
		}

		// Signal that no more merged chunks will be sent.
		close(out)

	}()

	return out
}

func filterWordlistWorker(
	wordlists []string,
	in <-chan Chunk,
	out chan<- Chunk,
) error {

	// TODO:
	// 1. sortiere jede subtract-list falls nötig
	// 2. merge alle subtract-listen
	// 3. vergleiche mit eingehenden Chunks

	return nil
}

func startFilterWordlist(
	wordlists []string,
	in <-chan Chunk,
) <-chan Chunk {

	out := make(chan Chunk)

	go func() {

		err := filterWordlistWorker(
			wordlists,
			in,
			out,
		)

		if err != nil {
			panic(err)
		}

		close(out)

	}()

	return out
}

// EditWordlist Process all word list edits in a logical order and, if necessary, pass the chunks on to the next filter
// inputPath - the path to the wordlist being processed, or stdin if “”
// outputPath - the path where the word list should be saved. stdout if “”
// filterRangeStr - specifies how the range should be filtered
// filterMaskStr - the masks that are to be filtered
func EditWordlist(
	inputPath string, outputPath string, muteStatusMessages bool,
	filterRangeArg []string, invRangeArg []string, // range filter args
	filterMaskArg []string, invMaskArg []string,   //
	filterRegExArg []string, invRegexArg []string,
	filterSubtractWordlistsArg []string, sortArg bool,
	removeDuplicatesArg bool,
) error {
	var input io.Reader

	var statusInputPath string
	var statusOutputPath string
	// INPUT VALIDATION
	// Input is read from stdin or a file --------------------------------------------------------------
	if inputPath == "" { // input is stdin
		input = os.Stdin
		statusInputPath = "stdin"
	} else { // file
		absInputPath, err := utils.ResolvePath(inputPath) // resolving path to an absolute path
		if err != nil {
			return err
		}
		statusInputPath = absInputPath

		openInputFile, err := os.Open(absInputPath) // open file
		if err != nil {
			return err
		}
		defer openInputFile.Close()
		input = openInputFile
	}
	// INPUT VALIDATION END --------------------------------------------------

	// OUTPUT VALIDATION
	// The output is prepared here for writing to a file or stdout --------------------------------------
	var outputWriter io.Writer
	if outputPath == "" { // output ist stdout
		outputWriter = os.Stdout
		statusOutputPath = "stdout"
	} else { // outputfile ist given
		absOutputPath, err := utils.ListPath(outputPath) // makes the path absolute and adds the correct suffix to the output file
		if err != nil {
			return err
		}
		statusOutputPath = absOutputPath

		openOutputFile, err := os.Create(absOutputPath)
		if err != nil {
			return err
		}
		defer openOutputFile.Close()
		outputWriter = openOutputFile
	}
	// OUTPUT VALIDATION END ------------------------------------------------

	// If range filters are to be applied, prepare them...
	var lineRanges []*rangeFilter
	if filterRangeArg != nil {
		for _, r := range filterRangeArg {
			lineRange, err := utils.NewUint64Range(r)
			if err != nil {
				return err
			}
			var f = &rangeFilter{rangeF: lineRange}
			lineRanges = append(lineRanges, f)
		}
	}
	var invLineRanges []*rangeFilter
	if invRangeArg != nil {
		for _, r := range invRangeArg {
			lineRange, err := utils.NewUint64Range(r)
			if err != nil {
				return err
			}
			var f = &rangeFilter{rangeF: lineRange}
			invLineRanges = append(invLineRanges, f)
		}
	}

	// If mask filters are to be applied, prepare them...
	var maskFilters []*maskFilter
	if filterMaskArg != nil {
		for _, m := range filterMaskArg {
			msk, err := generator.NewMaskInterpreter(m)
			if err != nil {
				return err
			}
			var f = &maskFilter{maskF: msk}
			maskFilters = append(maskFilters, f)
		}
	}
	var invMaskFilters []*maskFilter
	if invMaskArg != nil {
		for _, m := range invMaskArg {
			msk, err := generator.NewMaskInterpreter(m)
			if err != nil {
				return err
			}
			var f = &maskFilter{maskF: msk}
			invMaskFilters = append(invMaskFilters, f)
		}
	}

	// if regex filters are applied
	var regExFilters []*regExFilter
	if filterRegExArg != nil {
		for _, x := range filterRegExArg {
			regEx, err := regexp.Compile(x)
			if err != nil {
				return err
			}
			var f = &regExFilter{regEx: regEx}
			regExFilters = append(regExFilters, f)
		}
	}
	var invRegExFilters []*regExFilter
	if invRegexArg != nil {
		for _, x := range invRegexArg {
			regEx, err := regexp.Compile(x)
			if err != nil {
				return err
			}
			var f = &regExFilter{regEx: regEx}
			invRegExFilters = append(invRegExFilters, f)
		}
	}

	// prepare here more filters

	// loading chunk size from config file
	cfgPath, err := utils.GetConfigPath()
	if err != nil {
		return err
	}
	cfg, err := utils.LoadConfig(cfgPath)
	if err != nil {
		return err
	}

	if !muteStatusMessages {
		output.PrintStatus("statusEditor", "readWordlist", statusInputPath)
	}

	var threads int = runtime.NumCPU() * 4

	current := startReader( // channel who get read and changed
		input,
		cfg.General.ChunkByteSize,
		cfg.General.ByteLineLimit,
	)

	// DEBUG
	// fmt.Println("lineRanges:", len(lineRanges))
	// fmt.Println("invLineRanges:", len(invLineRanges))
	// fmt.Println("maskFilters:", len(maskFilters))
	// fmt.Println("invMaskFilters:", len(invMaskFilters))
	// fmt.Println("regExFilters:", len(regExFilters))
	// fmt.Println("invRegExFilters:", len(invRegExFilters))

	if lineRanges != nil || invLineRanges != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "rangeFilter")
		}

		current = startFilterRange(
			lineRanges,
			invLineRanges,
			threads,
			current,
		)

	}

	if maskFilters != nil || invMaskFilters != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "maskFilter")
		}
		current = startFilterMask(maskFilters, invMaskFilters, threads, current)

	}

	if regExFilters != nil || invRegExFilters != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "regExFilter")
		}
		current = startRegexFilter(regExFilters, invRegExFilters, threads, current)
	}

	if sortArg || removeDuplicatesArg {

		current = startExternalSort(
			threads,
			cfg.General.MultiThreadFileLines,
			removeDuplicatesArg,
			current,
		)
	}

	if !muteStatusMessages {
		output.PrintStatus("statusEditor", "writeWordlist", statusOutputPath)
	}

	return writeChunks(
		outputWriter,
		current,
	)
}
