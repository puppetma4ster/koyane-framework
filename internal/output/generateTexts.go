package output

var GenerateMessages = map[string]string{
	"use":       "generate",
	"short":     "generating wordlists using different methods",
	"long":      "...",
	"minLength": "Specifies the minimum word length.",
	"maxLength": "Specifies the maximum word length.",
	"compress":  "compresses wordlist into a .tar.xz archive after generation",
	"mask": `Generate wordlist from a pattern mask.
A mask consists of segments starting with '?' followed by letters that define the character type.
Example: '?ld?d?f' generates a word with a lowercase letter, a digit, and a special character.

Available wildcards:
  l = lowercase letter
  L = uppercase letter
  v = lowercase vowel
  V = uppercase vowel
  c = lowercase consonant
  C = uppercase consonant
  d = digit
  s = any special character
  f = common special characters
  p = dot special characters
  b = bracket special characters`,
}

var GenerateEditHelpTexts = map[string]string{
	"use":   "edit",
	"short": "A toolkit for editing word lists",
	"long":  "...",
	"sort":  "Sorts a wordlist with unicode",
	"invert": `Inverts the words in the word list.
If no PATTERN is specified, all words are inverted.
If PATTERN is specified (as a regular expression), only the words that match this expression are inverted.`,
	"removeFIle": `Removes all entries from the current wordlist that also appear in the specified <file>.
                        "The file must be a plain text list with one word per line.`,
	"removeMaske": `Removes all entries from the wordlist that match the given mask pattern.
                        The mask must follow the Koyane-Framework mask syntax (?d?d?d -> removes any 3-digit number).`,
}

var AnalyzeHelpTexts = map[string]string{
	"use":     "Analyze",
	"short":   "Analyze wordlists",
	"long":    "...",
	"all":     "Prints all gathered word list information",
	"general": "Prints all collected word list information belonging to the General Information category.",
	"content": "Prints all collected word list information belonging to the Content Information category.",
}

var GenerateRootHelpTexts = map[string]string{
	"use":   "koyane-framework",
	"short": "...",
	"long":  "Koyane-Framework :: wordlist forge & analysis toolkit made by Puppetm4ster",
}
