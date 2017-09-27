package api

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/cbor"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestFormulaHashing(t *testing.T) {
	baseFormula := Formula{
		Inputs: map[AbsPath]WareID{
			"/": WareID{"demo", "asdf"},
		},
		Action: FormulaAction{
			Exec: []string{"/bin/hello", "world"},
		},
		Outputs: map[AbsPath]OutputSpec{
			"/saveme": {PackFmt: "tar"},
		},
		FetchUrls: map[AbsPath][]WarehouseAddr{
			"/": []WarehouseAddr{
				"https+ca://ports.polydawn.io/assets/",
				"https+ca://mirror.wahoo.io/timeless/assets/",
			},
		},
		SaveUrls: map[AbsPath][]WarehouseAddr{
			"/saveme": []WarehouseAddr{
				"file+ca://./wares/",
			},
		},
	}

	t.Run("the CAS encoding should have the right stuff", func(t *testing.T) {
		msg, err := refmt.MarshalAtlased(cbor.EncodeOptions{}, baseFormula, FormulaCasAtlas)
		AssertNoError(t, err)

		t.Run("inputs should be present", func(t *testing.T) {
			if !bytes.Contains(msg, []byte("inputs")) {
				t.Errorf("failed to find input spec")
			}
			if !bytes.Contains(msg, []byte("demo:asdf")) {
				t.Errorf("failed to find input wareID")
			}
		})
		t.Run("action should be present", func(t *testing.T) {
			if !bytes.Contains(msg, []byte("exec")) {
				t.Errorf("failed to find action")
			}
			if !bytes.Contains(msg, []byte("/bin/hello")) {
				t.Errorf("failed to find action")
			}
		})
		t.Run("outputs should be present", func(t *testing.T) {
			if !bytes.Contains(msg, []byte("packfmt")) {
				t.Errorf("failed to find output spec")
			}
			if !bytes.Contains(msg, []byte("/saveme")) {
				t.Errorf("failed to find output spec")
			}
		})
		t.Run("fetchUrls should be absent", func(t *testing.T) {
			if bytes.Contains(msg, []byte("fetch")) {
				t.Errorf("should not find fetchUrls")
			}
			if bytes.Contains(msg, []byte("https+ca")) {
				t.Errorf("should not find fetchUrls")
			}
		})
		t.Run("saveUrls should be absent", func(t *testing.T) {
			if bytes.Contains(msg, []byte("saveUrl")) {
				t.Errorf("should not find saveUrls")
			}
			if bytes.Contains(msg, []byte("file+ca")) {
				t.Errorf("should not find saveUrls")
			}
		})
	})
	t.Run("Formula.Clone should DTRT", func(t *testing.T) {
		altFormula := baseFormula.Clone()
		if !reflect.DeepEqual(altFormula, baseFormula) {
			t.Errorf("clone method must yield an equivalent object")
		}
	})
}
