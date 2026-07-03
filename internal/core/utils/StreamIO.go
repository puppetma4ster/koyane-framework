package utils

// StreamIO implements the streaming reader and writer
// used by the editor pipeline.
//
// The reader splits the input stream into byte-limited
// chunks while the writer consumes processed chunks and
// writes them sequentially to the output.

import (
	"bufio"
	"io"
)

const (
	STDIN  = "stdin"
	STDOUT = "stdout"
)

// Chunk struct for managing individual chunks
// Index - chunk number to ensure correct sequence
// Lines - the words in the chunk
type Chunk struct {
	Index uint64
	Lines []string
}

// readChunksWorker reads the input streamline by line,
// splits the input into chunks and sends them through
// the output channel.
//
// input          - input stream (stdin or file)
// chunkByteSize  - maximum number of bytes per chunk
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

// StartReader creates the first stage of the pipeline.
// It reads the input stream, splits it into chunks
// and returns the output channel.
func StartReader(
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

// WriteChunks writes each processed chunk to stdout or to a file
// output - where to write (open file or stdout)
// in - channel containing the finished chunks
func WriteChunks(output io.Writer, in <-chan Chunk) error {
	// TODO: IMPLEMENT FILE SPLITTING IN SEPARATE FILES
	writer := bufio.NewWriter(output) // Create the Writer

	for chunk := range in {

		for _, line := range chunk.Lines {

			if _, err := writer.WriteString(line + "\n"); err != nil {
				return err
			}
		}

		if err := writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func ValidateInput(path string) (*io.Reader, error) {
	//var input io.Reader
	return nil, nil
}

func ValidateOutput(path string) (*io.Writer, error) {
	return nil, nil
}
