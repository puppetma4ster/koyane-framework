package editor

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

// Structure representing a single line from a file along with the index of its file
type fileLine struct {
	line string // the actual text line
	idx  int    // index of the file in the slice of files
}

// lineHeap implements a min-heap for fileLine nodes
type lineHeap []fileLine

func (h lineHeap) Len() int           { return len(h) }
func (h lineHeap) Less(i, j int) bool { return h[i].line < h[j].line } // min-heap based on line content
func (h lineHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *lineHeap) Push(x any)        { *h = append(*h, x.(fileLine)) }
func (h *lineHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MergeSortedFiles merges multiple already sorted files into a single sorted output file
func MergeSortedFiles(files []*os.File, outFile string) error {
	// Create the output file
	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	scanners := make([]*bufio.Scanner, len(files)) // scanner for each input file
	h := &lineHeap{}                               // min-heap to track the smallest line
	heap.Init(h)

	// Push the first line of each file into the heap
	for i, f := range files {
		_, err := f.Seek(0, 0) // rewind file to the beginning
		if err != nil {
			return err
		}
		scanners[i] = bufio.NewScanner(f)
		if scanners[i].Scan() { // read first line
			heap.Push(h, fileLine{line: scanners[i].Text(), idx: i})
		}
	}

	// Continue until the heap is empty
	for h.Len() > 0 {
		smallest := heap.Pop(h).(fileLine) // get the smallest line across all files
		fmt.Fprintln(out, smallest.line)   // write it to the output file

		// Push the next line from the same file into the heap
		if scanners[smallest.idx].Scan() {
			heap.Push(h, fileLine{line: scanners[smallest.idx].Text(), idx: smallest.idx})
		}
	}

	return nil
}
