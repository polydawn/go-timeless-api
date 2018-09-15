package repeatr

import (
	"go.polydawn.net/go-timeless-api/rio"
)

type ErrorCategory string

const (
	ErrUsage                = ErrorCategory("repeatr-usage-error")         // Indicates some piece of user input to a command was invalid and unrunnable.
	ErrWarehouseUnavailable = ErrorCategory(rio.ErrWarehouseUnavailable)   // The corresponding rio error halted execution.
	ErrWarehouseUnwritable  = ErrorCategory(rio.ErrWarehouseUnwritable)    // The corresponding rio error halted execution.
	ErrWareNotFound         = ErrorCategory(rio.ErrWareNotFound)           // The corresponding rio error halted execution.
	ErrWareCorrupt          = ErrorCategory(rio.ErrWareCorrupt)            // The corresponding rio error halted execution.
	ErrWareHashMismatch     = ErrorCategory(rio.ErrWareHashMismatch)       // The corresponding rio error halted execution.
	ErrCancelled            = ErrorCategory(rio.ErrCancelled)              // The corresponding rio error halted execution.
	ErrRioCacheProblem      = ErrorCategory(rio.ErrLocalCacheProblem)      // The corresponding rio error halted execution.
	ErrAssemblyInvalid      = ErrorCategory(rio.ErrAssemblyInvalid)        // The corresponding rio error halted execution.
	ErrPackInvalid          = ErrorCategory(rio.ErrPackInvalid)            // The corresponding rio error halted execution.
	ErrInoperablePath       = ErrorCategory(rio.ErrInoperablePath)         // The corresponding rio error halted execution.  (This one shouldn't show up much...!  Things like "out of disk" or such could still cause this though.)
	ErrFilterRejection      = ErrorCategory(rio.ErrFilterRejection)        // The corresponding rio error halted execution.  (This one shouldn't show up much!  By default, Repeatr uses "keep all" settings.  But users can still configure rejection filters.)
	ErrJobUnsuccessful      = ErrorCategory("repeatr-job-unsuccessful")    // Not an error -- indicates that the contained process exited nonzero.  TODO review if this needs an error category or just a reserved space in the exit code table.
	ErrJobInvalid           = ErrorCategory("repeatr-job-invalid")         // Indicates the container could not be launched because some part of its specification was invalid -- for example, the CWD requested is not a dir, or the command to exec is not an executable.  (The whole filesystem may have been necessary to set up before this can be detected.)
	ErrLocalCacheProblem    = ErrorCategory("repeatr-local-cache-problem") // Indicates an error while while handling internal filesystem paths (for example, if an executor can't mkdir its workspace dirs).
	ErrExecutor             = ErrorCategory("repeatr-executor-problem")    // Indicates an error occured while launching containment or handling the child processes.  Should be seen rarely -- comes up for esotera like "out of file handles".
	ErrRPCBreakdown         = ErrorCategory("repeatr-rpc-breakdown")       // Raised when running a remote process and the control channel is lost, the process fails to start, or unrecognized messages are received.

)

var ErrorTable = []struct {
	ExitCode     int
	RepeatrError ErrorCategory
}{
	// Codes 1 and 2 are generic.
	{ExitCode: 1 /*  */, RepeatrError: ErrUsage},
	{ExitCode: 2 /*  */, RepeatrError: ""}, // Reserved for panics and crashes.
	// The exit code ranges from rio, we keep same.
	{ExitCode: 3 /*  */, RepeatrError: ErrWarehouseUnavailable},
	{ExitCode: 4 /*  */, RepeatrError: ErrWarehouseUnwritable},
	{ExitCode: 5 /*  */, RepeatrError: ErrWareNotFound},
	{ExitCode: 6 /*  */, RepeatrError: ErrWareCorrupt},
	{ExitCode: 7 /*  */, RepeatrError: ErrWareHashMismatch},
	{ExitCode: 8 /*  */, RepeatrError: ErrCancelled},
	{ExitCode: 9 /*  */, RepeatrError: ErrRioCacheProblem},
	{ExitCode: 10 /* */, RepeatrError: ErrAssemblyInvalid},
	{ExitCode: 11 /* */, RepeatrError: ErrPackInvalid},
	{ExitCode: 12 /* */, RepeatrError: ErrInoperablePath},
	{ExitCode: 13 /* */, RepeatrError: ErrFilterRejection},
	// Let's give user job exit a nice round number:
	{ExitCode: 32 /* */, RepeatrError: ErrJobUnsuccessful},
	// Jump a few numbers, then repeatr exit codes begin:
	{ExitCode: 40 /* */, RepeatrError: ErrJobInvalid},
	{ExitCode: 41 /* */, RepeatrError: ErrLocalCacheProblem},
	{ExitCode: 42 /* */, RepeatrError: ErrExecutor},
	// Numbers do a big jump as we get into "you really shouldn't see these" territory...
	{ExitCode: 120 /**/, RepeatrError: ErrRPCBreakdown},
}
