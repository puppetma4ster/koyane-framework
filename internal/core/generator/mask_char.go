package generator

import (
	"fmt"
	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

// MaskChar – Parsing individual mask segments to determine allowed characters
//
// This file defines the `MaskChar` class, which parses individual mask segments
// (e.g., "?l", "?d", "?v", etc.) and determines which character groups are permitted
// for password or wordlist generation.
//
// Each instance of `MaskChar` interprets a mask segment, sets flags, and compiles
// a string of allowed characters to be used later for string generation.
//
// Masks can represent lowercase letters, uppercase letters, vowels, consonants,
// digits, and various classes of special characters.
//
// Example:
// MaskChar("?lV") allows all lowercase letters plus all uppercase vowels.
//
// Author: puppetm4ster /*
type MaskChar struct {
	Segment             string
	IsWildcard          bool
	PermittedCharacters string
}

func NewMaskChar(seg string, isWildcard bool) (*MaskChar, error) {
	var permittedChars string = ""

	if isWildcard {
		for _, char := range seg {
			if char == '?' {
				continue
			}
			switch char {
			case 'l':
				permittedChars += utils.LowerCaseCharacters
			case 'L':
				permittedChars += utils.UpperCaseCharacters
			case 'v':
				permittedChars += utils.LowerCaseVowels
			case 'V':
				permittedChars += utils.UpperCaseVowels
			case 'c':
				permittedChars += utils.LowerCaseConsonants
			case 'C':
				permittedChars += utils.UpperCaseConsonants
			case 'd':
				permittedChars += utils.Digits
			case 's':
				permittedChars += utils.SpecialCharacters
			case 'f':
				permittedChars += utils.SpecialCharactersMostUsed
			case 'p':
				permittedChars += utils.SpecialCharactersPoints
			case 'b':
				permittedChars += utils.SpecialCharactersBracelet
			default:
				return nil, fmt.Errorf("invalid Wildcard character %#U in ?%s", char, seg)
			}
		}
	} else {
		for _, char := range seg {
			if char == '!' {
				continue
			}
			permittedChars += string(char)
		}
	}
	return &MaskChar{
		Segment:             seg,
		IsWildcard:          isWildcard,
		PermittedCharacters: permittedChars,
	}, nil

}
