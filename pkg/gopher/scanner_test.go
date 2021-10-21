package gopher

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestScannerCodeAndFields(t *testing.T) {
	testCases := map[string][]struct {
		Code   string
		Fields []string
	}{
		"empty": {},
		"simple": {
			{"1", []string{"these", "are", "the", "fields"}},
			{"0", []string{"more", "fields"}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// build test text
			var builder strings.Builder
			for _, r := range tc {
				builder.WriteString(r.Code + strings.Join(r.Fields, "\t") + "\n")
			}
			builder.WriteString(".\n")

			// scan test text and compare results
			scanner := NewScanner(bytes.NewBufferString(builder.String()))

			for _, r := range tc {
				if !scanner.Scan() {
					t.Fatal("failed to scan code")
				}
				if scanner.Code() != r.Code || !reflect.DeepEqual(scanner.Fields(), r.Fields) {
					t.Fatal("failed to scan fields")
				}
			}
		})
	}
}
