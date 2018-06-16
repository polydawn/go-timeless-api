package api

import (
	"fmt"
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestModuleNameValidation(t *testing.T) {
	type tcase struct {
		Value ModuleName
		Error error
	}
	for _, tr := range []tcase{
		{"", fmt.Errorf("a moduleName cannot be an empty string")},
		{"yes", nil},
		{"yes.yes", nil},
		{"yes.yes/yes", nil},
		{"yes.yes/yes.yes", nil},
		{"yes.yes/yes..yes", nil},
		{"yes.yes/yes..yes/still", nil},
		{"no..no", fmtMatchError("moduleName", validation_dns1123Subdomain_msg)},
		{"no..no/narp", fmtMatchError("moduleName first segment", validation_dns1123Subdomain_msg)},
		{"sure/../butno", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
		{"sure/./alsono", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
		{"sure/NoCasePlz", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
		{"NoCaseDomains/nope", fmtMatchError("moduleName first segment", validation_dns1123Subdomain_msg)},
		{"noskipping//anysegments", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
		{"notrailing/", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
		{"notrailing/never/", fmtMatchError("moduleName path segment", validation_moduleNamePathHunk_msg)},
	} {
		t.Run(string(tr.Value), func(t *testing.T) {
			Wish(t, tr.Value.Validate(), ShouldEqual, tr.Error)
		})
	}
}
