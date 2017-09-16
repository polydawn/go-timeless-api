package testutil

import (
	"bytes"
	stdjson "encoding/json"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func WantStringEqual(t *testing.T, a, b string) {
	result, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	})
	if err != nil {
		t.Fatalf("diffing failed: %s", err)
	}
	t.Helper()
	if result != "" {
		t.Errorf("Match failed: diff:\n%s", result)
	}
}

func Dedent(s string) string {
	lines := strings.Split(s, "\n")
	lines = lines[1 : len(lines)-1]
	var prefixLen int
	for i, ch := range lines[0] {
		if ch != '\t' {
			prefixLen = i
			break
		}
	}
	for i := range lines {
		lines[i] = lines[i][prefixLen:]
	}
	return strings.Join(lines, "\n")
}

func PrettifyJson(s []byte) string {
	var out bytes.Buffer
	stdjson.Indent(&out, s, "", "\t")
	return out.String()
}
