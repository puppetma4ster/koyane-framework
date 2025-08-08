package analyzer

import (
	"bufio"
	"math"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

type AnalyzerContent struct {
	WordLines            uint64
	SmallestWordLen      int
	BiggestWordLen       int
	AvWordLen            float64
	CharCount            map[rune]uint64
	AvEntropy            float64
	HasDuplicates        bool
	DuplicateWords       []string
	WordsWDigits         float32
	WordsWUpper          float32
	WordsWSpecChar       float32
	WordsWDigitUpper     float32
	WordsWDigitSpec      float32
	WordsWUpperSpec      float32
	WordsWDigitUpperSpec float32
}

func NewAnalyzerContent(inputPath string, count, minMax, avLength, charFreq, avEntropy, duplicate, percStats bool) (*AnalyzerContent, error) {
	var wordlist AnalyzerContent = AnalyzerContent{
		WordLines:            0,
		SmallestWordLen:      0,
		BiggestWordLen:       0,
		AvWordLen:            0.0,
		CharCount:            make(map[rune]uint64),
		HasDuplicates:        false,
		DuplicateWords:       []string{},
		WordsWDigits:         0,
		WordsWUpper:          0,
		WordsWSpecChar:       0,
		WordsWDigitUpper:     0,
		WordsWDigitSpec:      0,
		WordsWUpperSpec:      0,
		WordsWDigitUpperSpec: 0,
	}
	absolutePath, err := utils.ResolvePath(inputPath)
	if err != nil {
		return nil, err
	}
	if duplicate {
		newTempPath, err := utils.GenerateRandomTempPath()
		if err != nil {
			return nil, err
		}
		err = utils.ExternalSort(absolutePath, newTempPath)
		if err != nil {
			return nil, err
		}
		absolutePath = newTempPath
	}
	file, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var totalWordLen uint64 = 0
	var totalEntropy float64 = 0.0
	var lastWord string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var word string = scanner.Text()

		if count || avLength {
			wordlist.WordLines += 1
		}
		if minMax {
			wordlist.passwordMinMaxInfo(word)
		}
		if avLength {
			totalWordLen += uint64(utf8.RuneCountInString(word))
		}
		if avEntropy {
			totalEntropy += calculateEntropy(word)
		}
		if charFreq {
			wordlist.charFrequency(word)
		}
		if duplicate {
			wordlist.duplicates(lastWord, word)
		}
		if percStats {
			wordlist.wordStats(word)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	if avLength {
		wordlist.AvWordLen = float64(totalWordLen / wordlist.WordLines)
	}
	if avEntropy {
		wordlist.AvEntropy = totalEntropy / float64(wordlist.WordLines)
	}
	if percStats {
		wordlist.statsInPercent()
	}
	if duplicate {
		err = os.Remove(absolutePath)
		if err != nil {
			return nil, err
		}
	}
	return &wordlist, nil
}
func NewContentDummy() *AnalyzerContent {
	return &AnalyzerContent{
		WordLines:            0,
		SmallestWordLen:      0,
		BiggestWordLen:       0,
		AvWordLen:            0.0,
		CharCount:            make(map[rune]uint64),
		HasDuplicates:        false,
		DuplicateWords:       []string{},
		WordsWDigits:         0,
		WordsWUpper:          0,
		WordsWSpecChar:       0,
		WordsWDigitUpper:     0,
		WordsWDigitSpec:      0,
		WordsWUpperSpec:      0,
		WordsWDigitUpperSpec: 0,
	}
}

func (wordlist *AnalyzerContent) passwordMinMaxInfo(word string) {
	var wordlength int = utf8.RuneCountInString(word)
	if wordlist.SmallestWordLen == 0 || wordlist.BiggestWordLen == 0 {
		wordlist.SmallestWordLen = wordlength
		wordlist.BiggestWordLen = wordlength
	}
	if wordlength < wordlist.SmallestWordLen {
		wordlist.SmallestWordLen = wordlength
	}
	if wordlength > wordlist.BiggestWordLen {
		wordlist.BiggestWordLen = wordlength
	}
}

func (wordlist *AnalyzerContent) charFrequency(word string) {
	for _, char := range word {
		if _, ok := wordlist.CharCount[char]; ok {
			wordlist.CharCount[char] += 1
		} else {
			wordlist.CharCount[char] = 1
		}
	}
}

func calculateEntropy(word string) float64 {
	freq := make(map[rune]float64)
	length := float64(len(word))

	// Zeichenhäufigkeiten zählen
	for _, ch := range word {
		freq[ch]++
	}

	var ent float64 = 0.0
	for _, count := range freq {
		p := count / length
		ent += -p * math.Log2(p)
	}

	return ent
}

func (wordlist *AnalyzerContent) duplicates(lastWord, currentWord string) {
	if lastWord == currentWord {
		wordlist.HasDuplicates = true
		if utils.NotInSlice(wordlist.DuplicateWords, currentWord) {
			wordlist.DuplicateWords = append(wordlist.DuplicateWords, currentWord)
		}
	}
}

func (wordlist *AnalyzerContent) wordStats(word string) {
	var digit bool = false
	var upper bool = false
	var specsign bool = false

	for _, char := range word {
		if unicode.IsDigit(char) {
			digit = true
		} else if unicode.IsUpper(char) {
			upper = true
		} else if unicode.IsSymbol(char) {
			specsign = true
		} else if unicode.IsSpace(char) {
			specsign = true
		} else if unicode.IsPunct(char) {
			specsign = true
		}
	}
	if digit && upper && specsign {
		wordlist.WordsWDigitUpperSpec += 1
	} else if digit && upper {
		wordlist.WordsWDigitUpper += 1
	} else if digit && specsign {
		wordlist.WordsWUpperSpec += 1
	} else if digit {
		wordlist.WordsWDigits += 1
	} else if upper {
		wordlist.WordsWUpper += 1
	} else if specsign {
		wordlist.WordsWSpecChar += 1
	}
}

func (wordlist *AnalyzerContent) statsInPercent() {
	const percentMultiplier float32 = 100
	wordlist.WordsWDigits = (wordlist.WordsWDigits / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWUpper = (wordlist.WordsWUpper / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWSpecChar = (wordlist.WordsWSpecChar / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitUpper = (wordlist.WordsWDigitUpper / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitSpec = (wordlist.WordsWDigitSpec / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWUpperSpec = (wordlist.WordsWUpperSpec / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitUpperSpec = (wordlist.WordsWDigitUpperSpec / float32(wordlist.WordLines)) * percentMultiplier

}
