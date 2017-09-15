package api

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestLol(t *testing.T) {
	t.Logf("%s", cmp.Diff("asdf\nsad", "asdf"))
}

func TestDiff(t *testing.T) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain("asdf\nsad\nnear match\noverlap", "wow\nasdf\nne arm atch\noverlap", false)
	t.Logf("%v", diffs) // so like you could just iterate this and print it quotedly -.-
	t.Logf("%s", dmp.DiffToDelta(diffs))
	t.Logf("%v", dmp.DiffLevenshtein(diffs))
	t.Logf("%v", dmp.DiffPrettyText(diffs))
}

//func TextualDiff(diffs []Diff) string {
//	var text bytes.Buffer
//	for _, aDiff := range diffs {
//		switch aDiff.Type {
//		case diffmatchpatch.DiffInsert:
//			_, _ = text.WriteString("+")
//			_, _ = text.WriteString(strings.Replace(url.QueryEscape(aDiff.Text), "+", " ", -1))
//			_, _ = text.WriteString("\t")
//			break
//		case diffmatchpatch.DiffDelete:
//			_, _ = text.WriteString("-")
//			_, _ = text.WriteString(strconv.Itoa(utf8.RuneCountInString(aDiff.Text)))
//			_, _ = text.WriteString("\t")
//			break
//		case diffmatchpatch.DiffEqual:
//			_, _ = text.WriteString("=")
//			_, _ = text.WriteString(strconv.Itoa(utf8.RuneCountInString(aDiff.Text)))
//			_, _ = text.WriteString("\t")
//			break
//		}
//	}
//	delta := text.String()
//	if len(delta) != 0 {
//		// Strip off trailing tab character.
//		delta = delta[0 : utf8.RuneCountInString(delta)-1]
//		delta = unescaper.Replace(delta)
//	}
//	return delta
//}
