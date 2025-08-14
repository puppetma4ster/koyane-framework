package generator

import "testing"

func TestMaskCharDigit(t *testing.T) {
	mask, err := NewMaskChar("?d", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "0123456789"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSmallLetters(t *testing.T) {
	mask, err := NewMaskChar("?l", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "abcdefghijklmnopqrstuvwxyz"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharBigLetters(t *testing.T) {
	mask, err := NewMaskChar("?L", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSmallVowels(t *testing.T) {
	mask, err := NewMaskChar("?v", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "aeiou"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharBigVowels(t *testing.T) {
	mask, err := NewMaskChar("?V", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "AEIOU"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSmallConsonats(t *testing.T) {
	mask, err := NewMaskChar("?c", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "bcdfghjklmnpqrstvwxyz"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharBigConsonats(t *testing.T) {
	mask, err := NewMaskChar("?C", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "BCDFGHJKLMNPQRSTVWXYZ"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSpCharMostUsed(t *testing.T) {
	mask, err := NewMaskChar("?f", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "!@#$%^&*()-_+=?"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSpCharPoints(t *testing.T) {
	mask, err := NewMaskChar("?p", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := ".,:;"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharSpCharBrachelet(t *testing.T) {
	mask, err := NewMaskChar("?b", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "()[]{}"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharAllSpChar(t *testing.T) {
	mask, err := NewMaskChar("?s", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "<>|^°!\"§$%&/()=?´{}[]\\¸`+~*#'-_.:,;@€"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharBigConsonantsAndDigits(t *testing.T) {
	mask, err := NewMaskChar("?Cd", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "BCDFGHJKLMNPQRSTVWXYZ0123456789"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharBigConsonantsAndDigitsReversed(t *testing.T) {
	mask, err := NewMaskChar("?dC", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "0123456789BCDFGHJKLMNPQRSTVWXYZ"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharFixed1(t *testing.T) {
	mask, err := NewMaskChar("!C", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "C"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharFixed2(t *testing.T) {
	mask, err := NewMaskChar("!ü", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "ü"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

func TestMaskCharFixed3(t *testing.T) {
	mask, err := NewMaskChar("!lrn390_:", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "lrn390_:"
	if mask.PermittedCharacters != expected {
		t.Fatalf("got %s, want %s", mask.PermittedCharacters, expected)
	}
}

// Error Testing

func TestMaskCharError(t *testing.T) {
	_, err := NewMaskChar("?z", true)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}

}
func TestMaskCharErrorMulti(t *testing.T) {
	_, err := NewMaskChar("?dz", true)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}

}
