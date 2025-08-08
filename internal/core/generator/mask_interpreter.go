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
	// Letztes Segment anhängen
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

func MatchesWord(mask *MaskInterpreter, word string) bool {
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
