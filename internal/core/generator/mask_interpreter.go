package generator

import (
	"strings"
)

type MaskInterpreter struct {
	Mask         string
	MaskSegments []MaskChar
}

func NewMaskInterpreter(msk string) (*MaskInterpreter, error) {
	const wildcard rune = '?'
	const fixed rune = '!'
	var mask string = strings.TrimSpace(msk)
	var maskSegments []MaskChar

	var segment strings.Builder
	var insideMask bool = false

	for _, char := range mask {
		if char == wildcard || char == fixed {
			if segment.Len() > 0 {
				mc, err := NewMaskChar(segment.String(), insideMask)
				if err != nil {
					return nil, err
				}
				maskSegments = append(maskSegments, *mc)
				segment.Reset()
			}
			segment.WriteRune(char)
			insideMask = (char == wildcard)
		} else {
			segment.WriteRune(char)
		}
	}
	// Append the last segment
	if segment.Len() > 0 {
		mc, err := NewMaskChar(segment.String(), insideMask)
		if err != nil {
			return nil, err
		}
		maskSegments = append(maskSegments, *mc)
	}

	return &MaskInterpreter{
		Mask:         mask,
		MaskSegments: maskSegments,
	}, nil
}

func (mask *MaskInterpreter) MatchesWord(word string) bool {
	wordRunes := []rune(word)

	if len(wordRunes) != len(mask.MaskSegments) {
		return false
	}
	for i, segment := range mask.MaskSegments {
		var match bool = false

		for _, char := range segment.PermittedCharacters {
			if wordRunes[i] == char {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	return true
}

// MaskSplitter Split masks into several partial masks to ensure parallelization.
//
// Returns:
//   - []*MaskInterpreter: array with partial masks of the original
func (mask *MaskInterpreter) MaskSplitter() []*MaskInterpreter {
	var masks []*MaskInterpreter

	for _, char := range mask.MaskSegments[0].PermittedCharacters {
		copyVal := *mask
		copyVal.MaskSegments = append([]MaskChar(nil), mask.MaskSegments...)
		copyVal.MaskSegments[0].PermittedCharacters = string(char)
		newMask := copyVal
		masks = append(masks, &newMask)
	}

	return masks
}
