package rio

import (
	"github.com/polydawn/go-errcat"
)

var ErrorTable = []struct {
	ExitCode int
	RioError ErrorCategory
}{
	{ExitCode: 1 /*  */, RioError: ErrUsage},
	{ExitCode: 2 /*  */, RioError: ""}, // Reserved for panics and crashes.
	{ExitCode: 3 /*  */, RioError: ErrWarehouseUnavailable},
	{ExitCode: 4 /*  */, RioError: ErrWarehouseUnwritable},
	{ExitCode: 5 /*  */, RioError: ErrWareNotFound},
	{ExitCode: 6 /*  */, RioError: ErrWareCorrupt},
	{ExitCode: 7 /*  */, RioError: ErrWareHashMismatch},
	{ExitCode: 8 /*  */, RioError: ErrCancelled},
	{ExitCode: 9 /*  */, RioError: ErrLocalCacheProblem},
	{ExitCode: 10 /* */, RioError: ErrAssemblyInvalid},
	{ExitCode: 11 /* */, RioError: ErrPackInvalid},
	{ExitCode: 110 /**/, RioError: ErrNotImplemented},
	{ExitCode: 120 /**/, RioError: ErrRPCBreakdown},
}

/*
	Utility function for Rio.

	Returns the exit code for a given rio.ErrorCategory.
*/
func GetExitCode(err error) int {
	if err == nil {
		return 0
	}
	category := errcat.Category(err)
	for _, row := range ErrorTable {
		if category == row.RioError {
			return row.ExitCode
		}
	}
	panic(errcat.Errorf(ErrRPCBreakdown, "no exit code mapping for error category %q", category))
}
