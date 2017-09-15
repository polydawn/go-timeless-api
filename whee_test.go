package api

import (
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func Test(t *testing.T) {
	helper(t, "a\nb\ncc\nmid\nmid\nmid\nmid\nmid\nmid\nmid\nmid\nd\n", "\nb\nc\nmid\nmid\nmid\nmid\nmid\nmid\nmid\nmid\nd")
}

func helper(t *testing.T, a, b string) {
	result, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	})
	if result != "" {
		t.Logf("Match failed: diff:\n%s", result)
	}
}
