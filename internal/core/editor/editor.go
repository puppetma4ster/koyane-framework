package editor

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"runtime"
	"sync"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/generator"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

// Chunk struct for managing individual chunks
// Index - chunk number to ensure correct sequence
// Lines - the words in the chunk
type Chunk struct {
	Index uint64
	Lines []string
}

// ReadChunks Reads a file and splits it into chunks
// input - file input stream
// chunkSize - how many lines does the chunk get
// out - piping the chunks
func readChunksWorker(input io.Reader, chunkSize int, out chan<- Chunk) error {
	scanner := bufio.NewScanner(input)

	chunk := Chunk{
		Index: 0,
		Lines: make([]string, 0, chunkSize),
	}

	for scanner.Scan() {
		chunk.Lines = append(chunk.Lines, scanner.Text())

		if len(chunk.Lines) == chunkSize { // if chunk is full

			out <- chunk // write chunk in channel

			chunk = Chunk{ // generate new empty chunk
				Index: chunk.Index + 1,
				Lines: make([]string, 0, chunkSize),
			}
		}
	}
	if err := scanner.Err(); err != nil { // checking for errors in scanner.Scan() for loop
		return err
	}

	if len(chunk.Lines) > 0 { // flushing last chunk
		out <- chunk
	}

	close(out) // closing last channel
	return nil // return nil - no errors where encountered
}

// startReader creates the first stage of the pipeline.
// It reads the input stream, splits it into chunks
// and returns the output channel.
func startReader(
	input io.Reader,
	chunkSize int,
) <-chan Chunk {

	out := make(chan Chunk)

	go func() {

		err := readChunksWorker(
			input,
			chunkSize,
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
func filterRangeWorker(r *utils.Uint64Range, in <-chan Chunk, out chan<- Chunk) {
	for chunk := range in { // loop through all the given chunks

		filtered := Chunk{ // create a new empty chunk
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines { // loop through the lines of the chunk to be filtered

			l := uint64(utf8.RuneCountInString(line)) // count the characters in the line

			// NOTE: Here, inversion can be easily implemented using if-else statements.
			if r.Min != nil && l < *r.Min { // the line character is less than the minimum range
				continue // skip
			}

			if r.Max != nil && l > *r.Max { // if the line length exceeds the maximum range
				continue // skip
			}

			filtered.Lines = append( // Write lines that match the range in the new chunk
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
	r *utils.Uint64Range,
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

// filterMaskWorker Specifies the logic for how masks are filtered
// msk - processed mask arg for filtering
// in - the chunks to be filtered
// out - the filtered chunks
func filterMaskWorker(msk *generator.MaskInterpreter, in <-chan Chunk, out chan<- Chunk) {
	for chunk := range in { // loop through all input chunks
		filtered := Chunk{ // create a new empty chunk
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines { // loop through all lines in the input chunk
			if generator.MatchesWord(msk, line) { // Filter out words that match the mask
				continue
			}
			filtered.Lines = append( // Write lines that match the range in the new chunk
				filtered.Lines,
				line,
			)
		}
		out <- filtered // Write the filtered chunk to the output channel
	}
}

// startFilterMask handles the go routines for the mask filter
// msk - processed mask arg for filtering
// workerCount - how many threads (workers) should be created
// in - the chunks to be filtered
func startFilterMask(msk *generator.MaskInterpreter, workerCount int, in <-chan Chunk) <-chan Chunk {
	out := make(chan Chunk) // implement new output channel

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines
		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterMaskWorker( // worker function
				msk,
				in,
				out)
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

func regexFilterWorker(
	re *regexp.Regexp,
	in <-chan Chunk,
	out chan<- Chunk,
) {

	for chunk := range in {

		filtered := Chunk{
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines {

			if re.MatchString(line) {
				continue
			}

			filtered.Lines = append(
				filtered.Lines,
				line,
			)
		}

		out <- filtered
	}
}

func startRegexFilter(
	re *regexp.Regexp,
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

func filterWordlistWorker() {

}

func startFilterWordlist() {

}

// EditWordlist Process all word list edits in a logical order and, if necessary, pass the chunks on to the next filter
// inputPath - the path to the wordlist being processed, or stdin if “”
// outputPath - the path where the word list should be saved. stdout if “”
// filterRangeStr - specifies how the range should be filtered
// frilterMaskStr - the masks that are to be filtered
func EditWordlist(inputPath string, outputPath string, filterRangeStr string, filterMaskStr string, filterRegExArg string) error {
	var input io.Reader
	// INPUT VALIDATION
	// Input is read from stdin or a file --------------------------------------------------------------
	if inputPath == "" { // input is stdin
		input = os.Stdin
	} else { // file
		absInputPath, err := utils.ResolvePath(inputPath) // resolving path to an absolute path
		if err != nil {
			return err
		}
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
	var output io.Writer
	if outputPath == "" { // output ist stdout
		output = os.Stdout
	} else { // outputfile ist given
		absOutputPath, err := utils.ListPath(outputPath) // makes the path absolute and adds the correct suffix to the output file
		if err != nil {
			return err
		}
		openOutputFile, err := os.Create(absOutputPath)
		if err != nil {
			return err
		}
		defer openOutputFile.Close()
		output = openOutputFile
	}
	// OUTPUT VALIDATION END ------------------------------------------------

	// If range filters are to be applied, prepare them...
	var lineRange *utils.Uint64Range
	if filterRangeStr != "" {
		var err error
		lineRange, err = utils.NewUint64Range(filterRangeStr)
		if err != nil {
			return err
		}
	}

	// If mask filters are to be applied, prepare them...
	var msk *generator.MaskInterpreter
	if filterMaskStr != "" {
		var err error
		msk, err = generator.NewMaskInterpreter(filterMaskStr)
		if err != nil {
			return err
		}
	}

	// if regex filters are applied
	var regExFilter *regexp.Regexp
	if filterRegExArg != "" {
		var err error
		regExFilter, err = regexp.Compile(filterRegExArg)
		if err != nil {
			return err
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

	current := startReader( // channel who get read and changed
		input,
		cfg.General.ChunkLineSize,
	)

	if lineRange != nil {

		current = startFilterRange(
			lineRange,
			runtime.NumCPU(),
			current,
		)
	}
	if msk != nil {

		current = startFilterMask(
			msk,
			runtime.NumCPU(),
			current,
		)
	}
	if regExFilter != nil {

		current = startRegexFilter(
			regExFilter,
			runtime.NumCPU(),
			current,
		)
	}

	return writeChunks(
		output,
		current,
	)
}
