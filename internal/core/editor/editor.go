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
	"github.com/puppetma4ster/koyane-framework/internal/output"
)

// Chunk struct for managing individual chunks
// Index - chunk number to ensure correct sequence
// Lines - the words in the chunk
type Chunk struct {
	Index uint64
	Lines []string
}

// readChunksWorker reads the input stream line by line,
// splits the input into chunks and sends them through
// the output channel.
//
// input          - input stream (stdin or file)
// chunkSize      - maximum number of lines per chunk
// byteLineLimit  - maximum allowed line size in bytes (0 = unlimited)
// out            - output channel for processed chunks
func readChunksWorker(input io.Reader, chunkSize int, byteLineLimit int, out chan<- Chunk) error {

	// Create a buffered reader
	reader := bufio.NewReader(input)

	// Create the first chunk
	chunk := Chunk{
		Index: 0,
		Lines: make([]string, 0, chunkSize),
	}

	for {

		// Read one line including the trailing '\n'
		line, err := reader.ReadBytes('\n')

		// Return unexpected errors
		if err != nil && err != io.EOF {
			return err
		}

		// Process the line if any bytes were read
		if len(line) > 0 {

			// Remove trailing newline
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}

			// Remove Windows carriage return
			if len(line) > 0 &&
				line[len(line)-1] == '\r' {

				line = line[:len(line)-1]
			}

			// Skip lines that exceed the configured limit
			if byteLineLimit > 0 && len(line) > byteLineLimit {

				if err == io.EOF {
					break
				}
				// TODO: implement warning log/output print
				continue
			}

			// Append the line to the current chunk
			chunk.Lines = append(
				chunk.Lines,
				string(line),
			)

			// Send the chunk if it is full
			if len(chunk.Lines) == chunkSize {

				out <- chunk

				// Create a new empty chunk
				chunk = Chunk{
					Index: chunk.Index + 1,
					Lines: make(
						[]string,
						0,
						chunkSize,
					),
				}
			}
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
	chunkSize int,
	byteLineLimit int,
) <-chan Chunk {

	out := make(chan Chunk)

	go func() {

		err := readChunksWorker(
			input,
			chunkSize,
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
func filterRangeWorker(r *utils.Uint64Range, in <-chan Chunk, out chan<- Chunk, invert bool) {
	for chunk := range in { // loop through all the given chunks

		filtered := Chunk{ // create a new empty chunk
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines { // loop through the lines of the chunk to be filtered

			l := uint64(utf8.RuneCountInString(line)) // count the characters in the line

			match := true

			if r.Min != nil && l < *r.Min {
				match = false
			}

			if r.Max != nil && l > *r.Max {
				match = false
			}

			if invert {
				match = !match
			}

			if match {
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

// startFilterRange manages the go routines to be executed for the range filter
// r - the processed range arg with min and max
// workerCount - how many threads (workers) should be created
// in - input stream of the chunks to be filtered
func startFilterRange(
	r *utils.Uint64Range,
	workerCount int,
	in <-chan Chunk,
	invert bool,
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
				invert,
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
func filterMaskWorker(msk *generator.MaskInterpreter, in <-chan Chunk, out chan<- Chunk, invert bool) {
	for chunk := range in { // loop through all input chunks
		filtered := Chunk{ // create a new empty chunk
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines { // loop through all lines in the input chunk

			match := generator.MatchesWord(msk, line)

			if invert {
				match = !match
			}

			if match { // Filter out words that match the mask
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
func startFilterMask(msk *generator.MaskInterpreter, workerCount int, in <-chan Chunk, invert bool) <-chan Chunk {
	out := make(chan Chunk) // implement new output channel

	var wg sync.WaitGroup // wait until the threads are finished

	for i := 0; i < workerCount; i++ { // generates goroutines
		wg.Add(1) // adds a goroutine to the workgroup

		go func() {
			defer wg.Done() // signal know that the thread is finished

			filterMaskWorker( // worker function
				msk,
				in,
				out,
				invert,
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

func regexFilterWorker(
	re *regexp.Regexp,
	in <-chan Chunk,
	out chan<- Chunk,
	invert bool,
) {

	for chunk := range in {

		filtered := Chunk{
			Index: chunk.Index,
			Lines: make([]string, 0, len(chunk.Lines)),
		}

		for _, line := range chunk.Lines {

			match := re.MatchString(line)

			if invert {
				match = !match
			}

			if match {
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
	invert bool,
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
				invert,
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
func EditWordlist(
	inputPath string, outputPath string, muteStatusMessages bool,
	filterRangeStr string, invRangeArg string, // range filter args
	filterMaskStr string, invMaskArg string, //
	filterRegExArg string, invRegexArg string) error {
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
	var lineRange *utils.Uint64Range
	if filterRangeStr != "" {
		var err error
		lineRange, err = utils.NewUint64Range(filterRangeStr)
		if err != nil {
			return err
		}
	}
	var invLineRange *utils.Uint64Range
	if invRangeArg != "" {
		var err error
		invLineRange, err = utils.NewUint64Range(filterRangeStr)
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
	var invMaskFilter *generator.MaskInterpreter
	if invMaskArg != "" {
		var err error
		invMaskFilter, err = generator.NewMaskInterpreter(invMaskArg)
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
	var invRegExFilter *regexp.Regexp
	if invRegexArg != "" {
		var err error
		invRegExFilter, err = regexp.Compile(invRegexArg)
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

	if !muteStatusMessages {
		output.PrintStatus("statusEditor", "readWordlist", statusInputPath)
	}
	current := startReader( // channel who get read and changed
		input,
		cfg.General.ChunkLineSize,
		cfg.General.ByteLineLimit,
	)

	if lineRange != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "rangeFilter")
		}

		current = startFilterRange(
			lineRange,
			runtime.NumCPU(),
			current,
			false,
		)
	}
	if invLineRange != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "rangeFilter")
		}

		current = startFilterRange(
			lineRange,
			runtime.NumCPU(),
			current,
			true,
		)
	}

	if msk != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "maskFilter")
		}

		current = startFilterMask(
			msk,
			runtime.NumCPU(),
			current,
			false,
		)
	}
	if invMaskFilter != nil {
		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "maskFilter")
		}

		current = startFilterMask(
			msk,
			runtime.NumCPU(),
			current,
			true,
		)
	}

	if regExFilter != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "regExFilter")
		}

		current = startRegexFilter(
			regExFilter,
			runtime.NumCPU(),
			current,
			false,
		)
	}
	if invRegExFilter != nil {

		if !muteStatusMessages {
			output.PrintStatus("statusEditor", "regExFilter")
		}

		current = startRegexFilter(
			regExFilter,
			runtime.NumCPU(),
			current,
			true,
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
