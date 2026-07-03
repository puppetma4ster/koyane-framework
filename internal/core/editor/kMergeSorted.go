package editor

import (
	"bufio"
	"container/heap"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"os"
)

// fileLine represents a single line read from a temporary
// sorted file together with the index of the file it came from.
//
// idx is needed so the merge process knows which scanner
// should provide the next line after this line has been consumed.
type fileLine struct {
	line string
	idx  int
}

// lineHeap implements a min-heap of fileLine objects.
//
// The heap always keeps the lexicographically smallest
// available line at the top.
//
// Example:
//
//	file1: apple, zebra
//	file2: banana, wolf
//
// Heap:
//
//	apple
//	banana
//
// The smallest entry can always be retrieved efficiently.
type lineHeap []fileLine

// Len returns the current number of elements in the heap.
func (h lineHeap) Len() int {
	return len(h)
}

// Less defines heap ordering.
//
// Smaller strings have higher priority and move towards
// the top of the heap.
func (h lineHeap) Less(i, j int) bool {
	return h[i].line < h[j].line
}

// Swap exchanges two heap elements.
func (h lineHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Push inserts a new element into the heap.
func (h *lineHeap) Push(x any) {
	*h = append(*h, x.(fileLine))
}

// Pop removes and returns the smallest element.
func (h *lineHeap) Pop() any {
	old := *h

	n := len(old)

	x := old[n-1]

	*h = old[:n-1]

	return x
}

// mergeSortedChunks performs a k-way merge over multiple
// already sorted temporary files.
//
// Each file must already be sorted.
//
// The merged result is emitted as Chunk objects through
// the output channel.
//
// If removeDuplicates is enabled, duplicate lines are
// removed during the merge process.
func mergeSortedChunks(
	paths []string,
	chunkSize int,
	out chan<- utils.Chunk,
	removeDuplicates bool,
) error {

	// Duplicate tracking.
	var lastLine string
	first := true

	// Open all temporary files.
	files := make([]*os.File, len(paths))

	for i, path := range paths {

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		files[i] = f

		// Cleanup:
		// Close and remove the temporary file when the
		// merge operation finishes.
		defer func(f *os.File) {
			f.Close()
			os.Remove(f.Name())
		}(f)
	}

	// Create one scanner per temporary file.
	scanners := make([]*bufio.Scanner, len(files))

	// Create and initialize the min-heap.
	h := &lineHeap{}

	heap.Init(h)

	// Load the first line of every file into the heap.
	//
	// After this step the heap contains the smallest
	// currently available line from each file.
	for i, f := range files {

		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}

		scanners[i] = bufio.NewScanner(f)

		if scanners[i].Scan() {

			heap.Push(
				h,
				fileLine{
					line: scanners[i].Text(),
					idx:  i,
				},
			)
		}
	}

	// Current output chunk.
	//
	// Lines are collected here until chunkSize is reached,
	// then the chunk is emitted through the output channel.
	chunk := utils.Chunk{
		Index: 0,
		Lines: make(
			[]string,
			0,
			chunkSize,
		),
	}

	// Continue merging until all files are exhausted.
	for h.Len() > 0 {

		// Retrieve the globally smallest line currently
		// available across all files.
		smallest := heap.Pop(h).(fileLine)

		// Remove duplicates if requested.
		if removeDuplicates {

			if !first && smallest.line == lastLine {

				// Even when skipping a duplicate we must
				// still advance the scanner belonging to
				// this file.
				if scanners[smallest.idx].Scan() {

					heap.Push(
						h,
						fileLine{
							line: scanners[smallest.idx].Text(),
							idx:  smallest.idx,
						},
					)
				}

				continue
			}

			first = false
			lastLine = smallest.line
		}

		// Add the line to the current output chunk.
		chunk.Lines = append(
			chunk.Lines,
			smallest.line,
		)

		// Flush the chunk once it reaches capacity.
		if len(chunk.Lines) >= chunkSize {

			out <- chunk

			chunk = utils.Chunk{
				Index: chunk.Index + 1,
				Lines: make(
					[]string,
					0,
					chunkSize,
				),
			}
		}

		// Read the next line from the same file
		// the popped line originated from.
		//
		// This keeps exactly one active line from each
		// file inside the heap.
		if scanners[smallest.idx].Scan() {

			heap.Push(
				h,
				fileLine{
					line: scanners[smallest.idx].Text(),
					idx:  smallest.idx,
				},
			)
		}
	}

	// Flush the final partially filled chunk.
	if len(chunk.Lines) > 0 {
		out <- chunk
	}

	return nil
}
