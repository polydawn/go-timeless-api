package api

import (
	"fmt"
	"regexp"
	"strings"
)

// Validate returns errors if the string is not a valid ModuleName.
// A ModuleName must resemble a domain name (per DNS-1123) with optional
// subsequent '/'-separated path segments, roughly like:
//
//	[[[...]subsubdomain.]subdomain.]domain[/path[/morepath[...]]]
//
// The domain segments are restricted to DNS-1123 characters, and
// and path segements restricted to [TODO:regexp].
// These requirements ensure that mapping module names onto a filesystem
// path is always a reasonable operation.
func (x ModuleName) Validate() error {
	if len(x) == 0 {
		return fmt.Errorf("a moduleName cannot be an empty string")
	}
	hunks := strings.Split(string(x), "/")
	switch len(hunks) {
	case 1:
		return validateDNS1123Subdomain("moduleName", hunks[0])
	default:
		if err := validateDNS1123Subdomain("moduleName first segment", hunks[0]); err != nil {
			return err
		}
		for _, hunk := range hunks[1:] {
			if err := validateModuleNamePathHunk("moduleName path segment", hunk); err != nil {
				return err
			}
		}
		return nil
	}
	panic("unreachable")
}

// similar to dns1123 label hunks, but allows mid-string dots also.
const validation_moduleNamePathHunk_regexpStr string = "[a-z0-9]([-a-z0-9\\.]*[a-z0-9])?"
const validation_moduleNamePathHunk_msg string = "must consist of lower case alphanumeric characters or '-' or '.', and must start and end with an alphanumeric character"
const validation_moduleNamePathHunk_maxlen int = 63

var validation_moduleNamePathHunk_regexp = regexp.MustCompile("^" + validation_moduleNamePathHunk_regexpStr + "$")

func validateModuleNamePathHunk(use string, value string) error {
	if len(value) > validation_moduleNamePathHunk_maxlen {
		return fmtMaxLenError(use, validation_moduleNamePathHunk_maxlen)
	}
	if !validation_moduleNamePathHunk_regexp.MatchString(value) {
		return fmtMatchError(use, validation_moduleNamePathHunk_msg)
	}
	return nil
}

// n.b. RFC 1035 *also* specifies DNS name sections, but only allows alphabetic first char.
const validation_dns1123Label_regexpStr string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
const validation_dns1123Label_msg string = "must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character"
const validation_dns1123Label_maxlen int = 63

var validation_dns1123Label_regexp = regexp.MustCompile("^" + validation_dns1123Label_regexpStr + "$")

func validateDNS1123Label(use string, value string) error {
	if len(value) > validation_dns1123Label_maxlen {
		return fmtMaxLenError(use, validation_dns1123Label_maxlen)
	}
	if !validation_dns1123Label_regexp.MatchString(value) {
		return fmtMatchError(use, validation_dns1123Label_msg)
	}
	return nil
}

const validation_dns1123Subdomain_regexpStr string = validation_dns1123Label_regexpStr + "(\\." + validation_dns1123Label_regexpStr + ")*"
const validation_dns1123Subdomain_msg string = "must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character"
const validation_dns1123Subdomain_maxlen int = 253

var validation_dns1123Subdomain_regexp = regexp.MustCompile("^" + validation_dns1123Subdomain_regexpStr + "$")

func validateDNS1123Subdomain(use string, value string) error {
	if len(value) > validation_dns1123Subdomain_maxlen {
		return fmtMaxLenError(use, validation_dns1123Subdomain_maxlen)
	}
	if !validation_dns1123Subdomain_regexp.MatchString(value) {
		return fmtMatchError(use, validation_dns1123Subdomain_msg)
	}
	return nil
}

func fmtMaxLenError(use string, length int) error {
	return fmt.Errorf("a %s must be no more than %d characters", use, length)
}

func fmtMatchError(use string, msg string) error {
	return fmt.Errorf("a %s %s", use, msg)
}
