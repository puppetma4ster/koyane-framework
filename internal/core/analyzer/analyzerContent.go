package analyzer

import (
	"bufio"
	"math"
	"os"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
)

type AnalyzerContent struct {
	WordLines                   uint64
	SmallestWordLen             int
	SmallestWordStr             string
	BiggestWordLen              int
	BiggestWordStr              string
	AvWordLen                   float64
	CharCount                   map[rune]uint64
	AvEntropy                   float64
	HasDuplicates               bool
	DuplicateWords              []string
	WordsWDigits                float32
	WordsWDigitsPercent         float32
	WordsWUpper                 float32
	WordsWUpperPercent          float32
	WordsWSpecChar              float32
	WordsWSpecCharPercent       float32
	WordsWDigitUpper            float32
	WordsWDigitUpperPercent     float32
	WordsWDigitSpec             float32
	WordsWDigitSpecPercent      float32
	WordsWUpperSpec             float32
	WordsWUpperSpecPercent      float32
	WordsWDigitUpperSpec        float32
	WordsWDigitUpperSpecPercent float32
}

func NewAnalyzerContent(inputPath string, count, minMax, avLength, charFreq, avEntropy, duplicate, percStats bool) (*AnalyzerContent, error) {
	var wordlist AnalyzerContent = *NewContentDummy()

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
			lastWord = word
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
		WordLines:                   0,
		SmallestWordLen:             0,
		SmallestWordStr:             "",
		BiggestWordLen:              0,
		BiggestWordStr:              "",
		AvWordLen:                   0.0,
		CharCount:                   make(map[rune]uint64),
		HasDuplicates:               false,
		DuplicateWords:              []string{},
		WordsWDigits:                0.0,
		WordsWDigitsPercent:         0.0,
		WordsWUpper:                 0.0,
		WordsWUpperPercent:          0.0,
		WordsWSpecChar:              0.0,
		WordsWSpecCharPercent:       0.0,
		WordsWDigitUpper:            0.0,
		WordsWDigitUpperPercent:     0.0,
		WordsWDigitSpec:             0.0,
		WordsWDigitSpecPercent:      0.0,
		WordsWUpperSpec:             0.0,
		WordsWUpperSpecPercent:      0.0,
		WordsWDigitUpperSpec:        0.0,
		WordsWDigitUpperSpecPercent: 0.0,
	}
}

func (wordlist *AnalyzerContent) passwordMinMaxInfo(word string) {
	var wordlength int = utf8.RuneCountInString(word)
	if wordlist.SmallestWordLen == 0 || wordlist.BiggestWordLen == 0 { // is not initialized
		wordlist.SmallestWordLen = wordlength
		wordlist.BiggestWordLen = wordlength

		wordlist.SmallestWordStr = word
		wordlist.BiggestWordStr = word
	}
	if wordlength < wordlist.SmallestWordLen { // is smaler
		wordlist.SmallestWordLen = wordlength
		wordlist.SmallestWordStr = word

	}
	if wordlength > wordlist.BiggestWordLen { // is bigger
		wordlist.BiggestWordLen = wordlength
		wordlist.BiggestWordStr = word

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
		wordlist.WordsWDigitUpperSpec += 1.0
	} else if digit && upper {
		wordlist.WordsWDigitUpper += 1.0
	} else if digit && specsign {
		wordlist.WordsWDigitSpec += 1.0
	} else if upper && specsign {
		wordlist.WordsWUpperSpec += 1.0
	} else if digit {
		wordlist.WordsWDigits += 1.0
	} else if upper {
		wordlist.WordsWUpper += 1.0
	} else if specsign {
		wordlist.WordsWSpecChar += 1.0
	}
}

func (wordlist *AnalyzerContent) statsInPercent() {
	const percentMultiplier float32 = 100
	wordlist.WordsWDigitsPercent = (wordlist.WordsWDigits / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWUpperPercent = (wordlist.WordsWUpper / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWSpecCharPercent = (wordlist.WordsWSpecChar / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitUpperPercent = (wordlist.WordsWDigitUpper / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitSpecPercent = (wordlist.WordsWDigitSpec / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWUpperSpecPercent = (wordlist.WordsWUpperSpec / float32(wordlist.WordLines)) * percentMultiplier
	wordlist.WordsWDigitUpperSpecPercent = (wordlist.WordsWDigitUpperSpec / float32(wordlist.WordLines)) * percentMultiplier

}
func mergeContentAnalyzers(wordlist1, wordlist2 *AnalyzerContent) *AnalyzerContent {
	var smallLen = min(wordlist1.SmallestWordLen, wordlist2.SmallestWordLen)
	var smallStr string = ""
	if smallLen == wordlist1.SmallestWordLen {
		smallStr = wordlist1.SmallestWordStr
	} else {
		smallStr = wordlist2.SmallestWordStr
	}

	var bigLen = min(wordlist1.SmallestWordLen, wordlist2.SmallestWordLen)
	var bigStr string = ""
	if bigLen == wordlist1.BiggestWordLen {
		bigStr = wordlist1.BiggestWordStr
	} else {
		bigStr = wordlist2.BiggestWordStr
	}

	var hasDup bool = false
	if wordlist1.HasDuplicates || wordlist2.HasDuplicates {
		hasDup = true
	}
	var dupSlice []string
	dupSlice = append(wordlist1.DuplicateWords, wordlist2.DuplicateWords...)

	for char, freq := range wordlist2.CharCount {
		wordlist1.CharCount[char] += freq
	}
	mergedContent := &AnalyzerContent{
		WordLines:                   wordlist1.WordLines + wordlist2.WordLines,
		SmallestWordLen:             smallLen,
		SmallestWordStr:             smallStr,
		BiggestWordLen:              bigLen,
		BiggestWordStr:              bigStr,
		AvWordLen:                   (wordlist1.AvWordLen*float64(wordlist1.WordLines) + wordlist2.AvWordLen*float64(wordlist2.WordLines)) / (float64(wordlist1.WordLines) + float64(wordlist2.WordLines)),
		CharCount:                   wordlist1.CharCount,
		AvEntropy:                   (wordlist1.AvEntropy*float64(wordlist1.WordLines) + wordlist2.AvEntropy*float64(wordlist2.WordLines)) / (float64(wordlist1.WordLines) + float64(wordlist2.WordLines)),
		HasDuplicates:               hasDup,
		DuplicateWords:              dupSlice,
		WordsWDigits:                wordlist1.WordsWDigits + wordlist2.WordsWDigits,
		WordsWDigitsPercent:         0,
		WordsWUpper:                 wordlist1.WordsWUpper + wordlist2.WordsWUpper,
		WordsWUpperPercent:          0,
		WordsWSpecChar:              wordlist1.WordsWSpecChar + wordlist2.WordsWSpecChar,
		WordsWSpecCharPercent:       0,
		WordsWDigitUpper:            wordlist1.WordsWDigitUpper + wordlist2.WordsWDigitUpper,
		WordsWDigitUpperPercent:     0,
		WordsWDigitSpec:             wordlist1.WordsWDigitSpec + wordlist2.WordsWDigitSpec,
		WordsWDigitSpecPercent:      0,
		WordsWUpperSpec:             wordlist1.WordsWUpperSpec + wordlist2.WordsWUpperSpec,
		WordsWUpperSpecPercent:      0,
		WordsWDigitUpperSpec:        wordlist1.WordsWDigitUpperSpec + wordlist2.WordsWDigitUpperSpec,
		WordsWDigitUpperSpecPercent: 0,
	}
	mergedContent.statsInPercent()

	return mergedContent
}

func ConcurrentContentAnalyzer(inputPath string, count, minMax, avLength, charFreq, avEntropy, duplicate, percStats bool) (*AnalyzerContent, error) {
	tempPaths, err := utils.SplitWordlist(inputPath) // splitting wordlist into pices
	if err != nil {
		return nil, err
	}
	var threadPool sync.WaitGroup
	channelContent := make(chan *AnalyzerContent)

	for _, tempPath := range tempPaths { // starting GoRoutines (threads)
		threadPool.Add(1)
		go func(path string) {
			defer threadPool.Done()
			conResult, err := NewAnalyzerContent(path, count, minMax, avLength, charFreq, avEntropy, duplicate, percStats)
			if err != nil {
				panic(err)
			}
			channelContent <- conResult
		}(tempPath)
	}

	go func() { // wait till every GoRoutine is finished
		threadPool.Wait()
		close(channelContent)
	}()

	var contentResults []*AnalyzerContent
	for content := range channelContent { // coping channel pointer returns into slice
		contentResults = append(contentResults, content)

	}
	err = utils.RemoveSplitWordlist(tempPaths) // delete old tempfiles
	if err != nil {
		return nil, err
	}

	// merging results
	var wordlist *AnalyzerContent = NewContentDummy()
	var firstEntry bool = true
	for _, content := range contentResults {
		if firstEntry {
			wordlist = content
			firstEntry = false
		} else {
			wordlist = mergeContentAnalyzers(wordlist, content)
		}
	}

	return wordlist, err
}
