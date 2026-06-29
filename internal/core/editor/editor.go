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

// removeRangeWorker specifies how the range is filtered
// r - range to filter
// in -the chunk being filtered
// out - the finished filtered chunk
func filterRangeWorker(filterRanges []*rangeFilter, invFilterRanges []*rangeFilter, in <-chan utils.Chunk, out chan<- utils.Chunk) {
	for chunk := range in { // loop through all the given chunks

		filtered := utils.Chunk{ // create a new empty chunk
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
	rangeFilterArg []string,
	InvRangeFilterArg []string,
	workerCount int,
	in <-chan utils.Chunk,
) (<-chan utils.Chunk, error) {

	out := make(chan utils.Chunk) // implement new output channel

	// If range filters are to be applied, prepare them...
	var lineRanges []*rangeFilter
	if rangeFilterArg != nil {
		for _, r := range rangeFilterArg {
			lineRange, err := utils.NewUint64Range(r)
			if err != nil {
				return nil, err
			}
			var f = &rangeFilter{rangeF: lineRange}
			lineRanges = append(lineRanges, f)
		}
	}
	var invLineRanges []*rangeFilter
	if InvRangeFilterArg != nil {
		for _, r := range InvRangeFilterArg {
			lineRange, err := utils.NewUint64Range(r)
			if err != nil {
				return nil, err
			}
			var f = &rangeFilter{rangeF: lineRange}
			invLineRanges = append(invLineRanges, f)
		}
	}

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines

		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterRangeWorker( // worker function
				lineRanges,
				invLineRanges,
				in,
				out,
			)
		}()
	}

	go func() { // wait for threads
		wg.Wait()
		close(out)
	}()

	return out, nil // return new channel
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
	masksFilter []*maskFilter,
	invMasksFilter []*maskFilter,
	in <-chan utils.Chunk,
	out chan<- utils.Chunk,
) {

	// Process all incoming chunks
	for chunk := range in {

		// Create a new chunk for filtered results
		filtered := utils.Chunk{
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
			keep := len(masksFilter) == 0

			// Check whether the line matches
			// at least one keep filter.
			for _, m := range masksFilter {

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
			for _, m := range invMasksFilter {

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
func startFilterMask(maskFilterArg []string, invMaskFilterArg []string, workerCount int, in <-chan utils.Chunk) (<-chan utils.Chunk, error) {
	out := make(chan utils.Chunk) // implement new output channel

	// If mask filters are to be applied, prepare them...
	var maskFilters []*maskFilter
	if maskFilterArg != nil {
		for _, m := range maskFilterArg {
			msk, err := generator.NewMaskInterpreter(m)
			if err != nil {
				return nil, err
			}
			var f = &maskFilter{maskF: msk}
			maskFilters = append(maskFilters, f)
		}
	}
	var invMaskFilters []*maskFilter
	if invMaskFilterArg != nil {
		for _, m := range invMaskFilterArg {
			msk, err := generator.NewMaskInterpreter(m)
			if err != nil {
				return nil, err
			}
			var f = &maskFilter{maskF: msk}
			invMaskFilters = append(invMaskFilters, f)
		}
	}

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines
		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterMaskWorker( // worker function
				maskFilters,
				invMaskFilters,
				in,
				out,
			)
		}()
	}
	go func() { // wait for threads
		wg.Wait()
		close(out)
	}()

	return out, nil // return new channel
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
	in <-chan utils.Chunk,
	out chan<- utils.Chunk,
) {

	// Process all incoming chunks
	for chunk := range in {

		// Create a new chunk for filtered results
		filtered := utils.Chunk{
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
	reFilterArg []string,
	reInvFilterArg []string,
	workerCount int,
	in <-chan utils.Chunk,
) (<-chan utils.Chunk, error) {

	out := make(chan utils.Chunk)

	// if regex filters are applied
	var regExFilters []*regExFilter
	if reFilterArg != nil {
		for _, x := range reFilterArg {
			regEx, err := regexp.Compile(x)
			if err != nil {
				return nil, err
			}
			var f = &regExFilter{regEx: regEx}
			regExFilters = append(regExFilters, f)
		}
	}
	var invRegExFilters []*regExFilter
	if reInvFilterArg != nil {
		for _, x := range reInvFilterArg {
			regEx, err := regexp.Compile(x)
			if err != nil {
				return nil, err
			}
			var f = &regExFilter{regEx: regEx}
			invRegExFilters = append(invRegExFilters, f)
		}
	}

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			regexFilterWorker(
				regExFilters,
				invRegExFilters,
				in,
				out,
			)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out, nil
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
	in <-chan utils.Chunk,
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
	in <-chan utils.Chunk,
) <-chan utils.Chunk {

	// Output channel containing the final merged chunks.
	out := make(chan utils.Chunk)

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
	in <-chan utils.Chunk,
	out chan<- utils.Chunk,
) error {

	// TODO:
	// 1. sortiere jede subtract-list falls nötig
	// 2. merge alle subtract-listen
	// 3. vergleiche mit eingehenden Chunks

	return nil
}

func startFilterWordlist(
	wordlists []string,
	in <-chan utils.Chunk,
) <-chan utils.Chunk {

	out := make(chan utils.Chunk)

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

	current := utils.StartReader( // channel who get read and changed
		input,
		cfg.General.ChunkByteSize,
		cfg.General.ByteLineLimit,
	)

	if filterRangeArg != nil || invRangeArg != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "rangeFilter")
		}

		current, err = startFilterRange(
			filterRangeArg,
			invRangeArg,
			threads,
			current,
		)
		if err != nil {
			return err
		}

	}

	if filterMaskArg != nil || invMaskArg != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "maskFilter")
		}
		current, err = startFilterMask(filterMaskArg, invMaskArg, threads, current)
		if err != nil {
			return err
		}

	}

	if filterRegExArg != nil || invRegexArg != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "regExFilter")
		}
		current, err = startRegexFilter(filterRegExArg, invRegexArg, threads, current)
		if err != nil {
			return err
		}
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

	return utils.WriteChunks(
		outputWriter,
		current,
	)
}
