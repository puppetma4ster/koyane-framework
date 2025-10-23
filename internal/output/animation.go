package output

import (
	"fmt"
	"time"
)

// Spinner shows a rotating spinner animation together with a growing dot sequence.
// The spinner rotates every 200ms, while the dots increase only every 400ms.
// Example output: "| Loading.", "/ Loading..", "- Loading...", "\ Loading...."
func Spinner(message string, stop <-chan struct{}) {
	frames := []string{"|", "/", "-", "\\"} // Spinner characters for rotation
	frameIndex := 0                         // Current spinner frame index
	points := 1                             // Counter for the number of dots
	step := 0                               // Counts spinner updates to slow down dot growth

	for {
		select {
		case <-stop: // Stop signal received
			fmt.Print("\r\033[K") // Clear line and stop spinner
			fmt.Println()
			return
		default:
			// Print the spinner message:
			// \r moves to the beginning of the current line
			// \033[K clears the rest of the line to avoid leftover characters
			// frames[frameIndex] chooses the current spinner frame
			// makeDots(points) generates a string of '.' characters
			fmt.Print("\r\033[K")
			fmt.Printf("%s %s%s", frames[frameIndex], message, string(makeDots(points)))

			// Wait before updating (spinner speed = 200ms)
			time.Sleep(200 * time.Millisecond)

			// Move to the next spinner frame
			frameIndex = (frameIndex + 1) % len(frames)

			// Update the dots only every 2nd spinner step (400ms)
			step++
			if step%4 == 0 {
				points++
				if points > 5 { // Reset the dots after reaching 5
					points = 1
				}
			}
		}
	}
}

// makeDots creates a slice of runes containing 'n' dots.
// This is converted to a string when used in the spinner.
func makeDots(n int) []rune {
	dots := make([]rune, n)
	for i := 0; i < n; i++ {
		dots[i] = '.'
	}
	return dots
}
