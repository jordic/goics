package goics

import (
	"reflect"
	"strings"
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		want   string
		length int
	}{
		{"Anders", "And", 3},
		{"åååååå", "ååå", 6},
		// Hiragana a times 4
		{"\u3042\u3042\u3042\u3042", "\u3042", 4},
		{"\U0001F393", "", 1},
		// Continuation bytes
		{"\x80\x80", "", 1},
	}
	for _, test := range tests {
		if got := truncateString(test.input, test.length); got != test.want {
			t.Errorf("expected %q, got %q", test.want, got)
		}
	}
}

func TestSplitLength(t *testing.T) {
	tests := []struct {
		input string
		len   int
		want  []string
	}{
		{
			"AndersSrednaFoobarBazbarX",
			6,
			[]string{"Anders", "Sredna", "Foobar", "Bazbar", "X"},
		},
		{
			"AAAA\u00c5\u00c5\u00c5\u00c5\u3042\u3042\u3042\u3042\U0001F393\U0001F393\U0001F393\U0001F393",
			4,
			[]string{
				"AAAA",                         // 1 byte
				"\u00c5\u00c5", "\u00c5\u00c5", // 2 bytes
				"\u3042", "\u3042", "\u3042", "\u3042", // 3 bytes
				"\U0001F393", "\U0001F393", "\U0001F393", "\U0001F393", // 4 bytes
			},
		},
		{
			"\u3042\u3042\u3042\u3042",
			8,
			[]string{"\u3042\u3042", "\u3042\u3042"},
		},
		{
			"\u3042",
			2,
			nil,
		},
	}
	for _, tt := range tests {
		if got := splitLength(tt.input, tt.len); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("splitLength() = %v, want %v", got, tt.want)
		}
	}
}

func BenchmarkTruncateString(b *testing.B) {
	longString := strings.Repeat("\u3042", 100)
	for i := 0; i < b.N; i++ {
		truncateString(longString, 150)
	}
}
