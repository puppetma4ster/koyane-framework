package output

type StatusCategory struct {
	Prefix   string
	Messages map[string]string
}

var StatusMessages = map[string]StatusCategory{
	"errors": {
		Prefix: "[-]",
		Messages: map[string]string{
			"error":          "An unexpected error occurred:\n %s",
			"invID":          "An invalid ID was given: %S",
			"wrongMaxLength": "The `max_length` parameter cannot be 0 or negative!",
			"wrongMinLength": "The `minlength` parameter cannot be less than 0 or greater than the `maxlength` parameter!",
			"noOutputPath":   "No output path was specified!",
		},
	},
	"warnings": {
		Prefix: "[!]",
		Messages: map[string]string{
			"invView": "️Invalid view argument: %s, switching to default view \"summary\"",
		},
	},
	"statusRoot": {
		Prefix: "[*]",
		Messages: map[string]string{
			"generateTemp": "Generate temp path if it does not yet exist",
		},
	},
	"statusGenerator": {
		Prefix: "[*]",
		Messages: map[string]string{
			"calculateWords":       "Final wordlist has %s entries.",
			"calculateSize":        "Final wordlist size: %s",
			"buildingMaskWordlist": "Building wordlist using the following mask: '%s'",
			"wordlist_stats":       "The wordlist contains %d entries and is approximately %s bytes in size",
			"compress_wordlist":    "Compress generated wordlist to: %s",
		},
	},
	"successGenerator": {
		Prefix: "[+]",
		Messages: map[string]string{
			"wordlistCreated": "Wordlist successfully created at: %s",
			"archiveCreated":  "Compressed wordlist successfully created at: %s",
		},
	},
	"statusEditor": {
		Prefix: "[*]",
		Messages: map[string]string{
			"readWordlist":   "Read wordlist from: %s",
			"rangeFilter":    "Applying range filter",
			"invRangeFilter": "Applying inverted range filter",
			"maskFilter":     "Applying mask filter",
			"invMaskFilter":  "Applying inverted mask filter",
			"regExFilter":    "Applying RegEx filter",
			"invRegExFilter": "Applying RegEx filter",
			"writeWordlist":  "Write wordlist to",
		},
	},
	"statusGet": {
		Prefix: "[*]",
		Messages: map[string]string{
			"try": "Try to download wordlist from: %s",
		},
	},
	"successGet": {
		Prefix: "[+]",
		Messages: map[string]string{
			"succeeded": "The list was successfully downloaded and saved at path: %s",
		},
	},
}
