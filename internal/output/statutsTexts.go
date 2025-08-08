package output

type StatusCategory struct {
	Prefix   string
	Messages map[string]string
}

var StatusMessages = map[string]StatusCategory{
	"errors": {
		Prefix: "[-]",
		Messages: map[string]string{
			"error": "An unexpected error occurred:\n %s",
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
			"calculateWords":       "Final wordlist has %d entries.",
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
}
