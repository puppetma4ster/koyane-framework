package utils

import (
	"fmt"
	"strings"
	"testing"
)

// SplitWordlist Test
func TestSplitWordlist(t *testing.T) {
	paths, err := SplitWordlist("/home/puppetm4ster/Downloads/cleadrockyou.txt")
	fmt.Printf("%#v\n", paths) // []string{"vodka", "kahlua", "sahne"}
	if err != nil {
		t.Fatalf("SplitWordlist returned error: %v", err)
	}

	prefix := "/tmp/koyane_framework_tmp/"

	for i, p := range paths {
		if !strings.HasPrefix(p, prefix) {
			t.Errorf("Eintrag %d (%q) startet nicht mit %q", i, p, prefix)
		}
	}
}
