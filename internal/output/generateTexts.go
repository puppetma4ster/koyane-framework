package output

var GenerateMessages = map[string]string{
	"use":   "kyfgen",
	"short": "generating wordlists using different methods",
	"long":  "Generate new word lists for hash, directory, and web fuzzing analyses.",

	"minLength": "Specifies the minimum char length for permutation and mask generation.",
	"maxLength": "Specifies the maximum char length length for permutation.",
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
	"extract-hc-potfile": "extracts the passwords from a Hashcat potfile",
	"output":             "output path for the saved word lists",
	"permutation":        "enter a path to a word list that is to be combined in all possible permutations. The strings must be separated by line breaks.",
}

var GenerateEditHelpTexts = map[string]string{
	"use":   "kyfedit [OPTIONS] [inputlist] [outputlist]",
	"short": "A toolkit for editing word lists",
	"long": "Edit an existing password list or configuration file. " +
		"This command allows you to update, remove, or add entries interactively or via command-line options.",
	"sort": "Sorts a wordlist with unicode",
	"invert": `Inverts the words in the word list.
If no PATTERN is specified, all words are inverted.
If PATTERN is specified (as a regular expression), only the words that match this expression are inverted.`,
	"removeFile": `Removes all entries from the current wordlist that also appear in the specified <file>.
                        "The file must be a plain text list with one word per line.`,
	"removeMask": `Removes all entries from the wordlist that match the given mask pattern.
                        The mask must follow the Koyane-Framework mask syntax (?d?d?d -> removes any 3-digit number).`,
	"removeRange": `Filters words by length. Only words that fit within the range will be output. 
						A range is specified as follows -> minChar:maxChar (example: 3:3 -> only words that are 3 characters long). 
						If a value is not to be used, you can use the “*” wildcard 
						(Example:
						-r 3:3 -> only words that are 3 characters long
						-r 3:* -> only words that have at least 3 characters) .`,
	"removeChars": `Deletes all words containing the specified characters that were passed as arguments. 
						You can specify as many characters as you want. 
						(Example: 
						-c a -> all words containing an “a” are deleted, 
						-c 1234567890 -> all words containing numbers are deleted)`,
	"european": `Cleans the list of all non-European characters such as Arabic, Armenian, Chinese, Japanese, including emojis.`,
	"delete":   "deletes the input list ",
}

var AnalyzeHelpTexts = map[string]string{
	"use":      "kyfinfo",
	"short":    "Analyze wordlists",
	"long":     "Analyzes an existing word list and lists information such as character statistics about the list",
	"all":      "Prints all gathered word list information",
	"general":  "Prints all collected word list information belonging to the General Information category.",
	"content":  "Prints all collected word list information belonging to the Content Information category.",
	"stats":    "Prints all gathered statistics about the wordlist",
	"saveFile": "Creates a file with the output print. A path with the name of the file to be created must be specified.",
}

var SearchHelpTexts = map[string]string{
	"use":               "search [flags]",
	"short":             "Search wordlists",
	"long":              "Search and find Word Lists and Rules for Various Purposes",
	"name":              "filters only displays wordlists that match the name. The name must be specified.",
	"words":             "Only displays word lists that contain at least or at most the specified entities. Argument must be in range format (min:max).",
	"smallWord":         "Filter wordlists by the length of the smallest word. Argument must be in range format (min:max).",
	"bigWord":           "Filter wordlists by the length of the biggest word. Argument must be in range format (min:max).",
	"size":              "Filters word lists by memory size. Argument must be in range format (min:max).",
	"avgEntity":         "Filters word lists by average length. Argument must be in range format (min:max) (Float).",
	"avgEntropy":        "Filters word lists by average entropy. Argument must be in range format (min:max) (Float).",
	"digitsPerc":        "Filters word lists by digit percent value. Argument must be in range format (min:max) (Float).",
	"upperPerc":         "Filters word lists by upper case percent value. Argument must be in range format (min:max) (Float).",
	"specialPerc":       "Filters word lists by special Percent percent value. Argument must be in range format (min:max) (Float).",
	"digitsUpperPerc":   "Filters word lists by digit and upper case percent value. Argument must be in range format (min:max) (Float).",
	"digitsSpecialPerc": "Filters word lists by digit and special character percent value. Argument must be in range format (min:max) (Float).",
	"upperSpecialPerc":  "Filters word lists by upper case and special character percent value. Argument must be in range format (min:max) (Float).",
	"allPerc":           "Filters word lists by digit, upper case and special character percent value. Argument must be in range format (min:max) (Float).",
	"encoding":          "Filters word lists by encoding.",
	"language":          "Filters word lists by language.",
	"category":          "Filters word lists by categories.",
	"author":            "Filters word lists by author.",
	"tags":              "Filters word lists by tags.",
	"url":               "Filters word lists by url",
}

var ShowHelpTexts = map[string]string{
	"use":   "show [flags] wordlistID",
	"short": "Show wordlist statistics",
	"long":  "Detailed listing of all attributes and statistics of a word list in the database. To make a selection, the database ID must be entered at the end WITHOUT a flag.",
}

var GetHelpTexts = map[string]string{
	"use":   "kyfdb",
	"short": "Download Wordlists",
	"long":  "Download wordlists from the database using the list ID.",
	"path":  "A custom download path if the path from the configuration YAML should not be used. To make a selection, the database ID must be entered at the end WITHOUT a flag.",
}

var GenerateRootHelpTexts = map[string]string{
	"use":     "kyfdb",
	"short":   "Framework for wordlists",
	"long":    "Koyane-Framework :: wordlist forge & analysis toolkit made by Puppetm4ster",
	"toggle":  "",
	"version": "shows the version of the program",
}
