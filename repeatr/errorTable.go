package repeatr

import (
	"github.com/polydawn/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
)

/*
	Error mapping table.

	Where both a repeatr and a rio variant are listed in the same row,
	they are considered equivalent, and the rio error may be mapped to
	to the repeatr category (this mapping should occur at all edge points).

	Where an exit code and a rio category are listed, but no repeatr category,
	it is to leave an explicit gap in the exit codes of repeatr for clarity.
*/
var ErrorTable = []struct {
	ExitCode     int
	RepeatrError ErrorCategory
	RioError     rio.ErrorCategory
}{
	// Codes 1 and 2 are generic.
	{ExitCode: 1 /*  */, RepeatrError: ErrUsage /*                */, RioError: ""}, // Rio ErrUsage not remapped, because internal calls shouldn't get usage errors!
	{ExitCode: 2 /*  */, RepeatrError: "" /*                      */, RioError: ""}, // Reserved for panics and crashes.
	// The exit code ranges from rio, we keep same.
	{ExitCode: 3 /*  */, RepeatrError: ErrWarehouseUnavailable /* */, RioError: rio.ErrWarehouseUnavailable}, // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 4 /*  */, RepeatrError: ErrWarehouseUnwritable /*  */, RioError: rio.ErrWarehouseUnwritable},  // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 5 /*  */, RepeatrError: ErrWareNotFound /*         */, RioError: rio.ErrWareNotFound},         // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 6 /*  */, RepeatrError: ErrWareCorrupt /*          */, RioError: rio.ErrWareCorrupt},          // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 7 /*  */, RepeatrError: ErrWareHashMismatch /*     */, RioError: rio.ErrWareHashMismatch},     // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 8 /*  */, RepeatrError: ErrCancelled /*            */, RioError: rio.ErrCancelled},            // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 9 /*  */, RepeatrError: ErrRioCacheProblem /*      */, RioError: rio.ErrLocalCacheProblem},    // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 10 /* */, RepeatrError: ErrAssemblyInvalid /*      */, RioError: rio.ErrAssemblyInvalid},      // Translation of rio errors.  (Strings even stay same; types do not.)
	{ExitCode: 11 /* */, RepeatrError: ErrPackInvalid /*          */, RioError: rio.ErrPackInvalid},          // Translation of rio errors.  (Strings even stay same; types do not.)
	// Jump a few numbers, then repeatr exit codes begin:
	{ExitCode: 20 /* */, RepeatrError: ErrLocalCacheProblem},
	{ExitCode: 21 /* */, RepeatrError: ErrExecutor},
	// Let's give user job exit a nice round number:
	{ExitCode: 31 /* */, RepeatrError: ErrJobInvalid},
	{ExitCode: 32 /* */, RepeatrError: ErrJobUnsuccessful},
	// Numbers do a big jump as we get into "you really shouldn't see these" territory...
	{ExitCode: 110 /**/, RepeatrError: "" /*              */, RioError: rio.ErrNotImplemented},
	{ExitCode: 120 /**/, RepeatrError: "" /*              */, RioError: rio.ErrRPCBreakdown},
	{ExitCode: 121 /**/, RepeatrError: ErrRPCBreakdown /* */, RioError: ""},
}

/*
	Utility function for Repeatr.

	Returns the exit code for a given ErrorCategory.
*/
func ExitCodeForError(err error) int {
	if err == nil {
		return 0
	}
	return ExitCodeForCategory(errcat.Category(err))
}

/*
	Utility function for Repeatr.

	Returns the exit code for a given ErrorCategory.
*/
func ExitCodeForCategory(category interface{}) int {
	for _, row := range ErrorTable {
		if category == row.RepeatrError {
			return row.ExitCode
		}
	}
	panic(errcat.Errorf(ErrRPCBreakdown, "no exit code mapping for error category %q", category))
}

/*
	Utility function for Repeatr.

	Filter errors from rio into the corresponding repeatr.ErrorCategory.
	Returns repeatr.ErrRPCBreakdown for unexpected errors.
*/
func ReboxRioError(err error) error {
	category := errcat.Category(err)
	switch category.(type) {
	case nil:
		return nil
	case rio.ErrorCategory:
		for _, row := range ErrorTable {
			if category != row.RioError {
				continue
			}
			if row.RepeatrError == "" {
				continue
			}
			return errcat.Recategorize(row.RepeatrError, err)
		}
		return errcat.Errorf(ErrRPCBreakdown, "protocol error: unexpected error category %q from rio (error was: %s)", category, err)
	default:
		return errcat.Errorf(ErrRPCBreakdown, "protocol error: unexpected error category type %T from rio (error was: %s)", category, err)
	}
}
