package rio

import (
	"github.com/warpfork/go-errcat"
)

type ErrorCategory string

const (
	ErrUsage                = ErrorCategory("rio-usage-error")           // Indicates some piece of user input to a command was invalid and unrunnable.
	ErrWarehouseUnavailable = ErrorCategory("rio-warehouse-unavailable") // Warehouse 404.
	ErrWarehouseUnwritable  = ErrorCategory("rio-warehouse-unwritable")  // Indicates a warehouse failed to accept a write operation.  The Warehouse is having a bad day.  ("unauthorized" is a different error category.)
	ErrWareNotFound         = ErrorCategory("rio-ware-not-found")        // Ware 404 -- warehouse appeared online, but the requested ware isn't in it.
	ErrWareCorrupt          = ErrorCategory("rio-ware-corrupt")          // Incidates a ware retreival started, but during unpacking it was found to be malformed.
	ErrWareHashMismatch     = ErrorCategory("rio-hash-mismatch")         // Returned when fetching and unpacking a ware gets results in a different content hash than we requested.  (This is distinct from ErrWareCorrupt because a full fileset *was* able to be unpacked; it's just not the one we asked for.)
	ErrCancelled            = ErrorCategory("rio-cancelled")             // The operation timed out or was cancelled
	ErrLocalCacheProblem    = ErrorCategory("rio-local-cache-problem")   // Indicates an error while either reading or writing to rio's local fileset caches.
	ErrAssemblyInvalid      = ErrorCategory("rio-assembly-invalid")      // Indicates an error in unpack or tree-unpack where the requested set of unpacks cannot assemble cleanly (e.g. a tree where a /file is a file and another unpack tries to seat something at /file/dir; this assembly is impossible).
	ErrPackInvalid          = ErrorCategory("rio-pack-invalid")          // Indicates a pack could not be performed, perhaps because you tried to pack a file using a format that must start with dirs, or because of permission errors or other misfortunes during the pack.
	ErrInoperablePath       = ErrorCategory("rio-inoperable-path")       // Indicates pack or unpack failed while reading or writing the target local filesystem path (permissions errors, etc, are likely causes).
	ErrFilterRejection      = ErrorCategory("rio-filter-rejection")      // Indicates filters rejected some part of a fileset, e.g. `rio pack tar /dev --filter dev=reject` will return this error.
	ErrRPCBreakdown         = ErrorCategory("rio-rpc-breakdown")         // Raised when running a remote rio process and the control channel is lost, the process fails to start, or unrecognized messages are received.
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
	{ExitCode: 12 /* */, RioError: ErrInoperablePath},
	{ExitCode: 13 /* */, RioError: ErrFilterRejection},
	{ExitCode: 120 /**/, RioError: ErrRPCBreakdown},
}

// ExitCodeForError translates an error into a numeric exit code, looking up
// a code based on the errcat category of the error.
func ExitCodeForError(err error) int {
	if err == nil {
		return 0
	}
	return ExitCodeForCategory(errcat.Category(err))
}

// ExitCodeForCategory translates an errcat category into a numeric exit code.
func ExitCodeForCategory(category interface{}) int {
	for _, row := range ErrorTable {
		if category == row.RioError {
			return row.ExitCode
		}
	}
	panic(errcat.Errorf(ErrRPCBreakdown, "no exit code mapping for error category %q", category))
}
