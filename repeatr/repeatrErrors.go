package repeatr

import (
	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
)

type ErrorCategory string

// Error implements the errcat.Error interface, specifically using
// this package's ErrorCategory as the concrete category,
// and giving us a type to hang custom serialization on.
type Error struct {
	Category_ ErrorCategory     `json:"category"          refmt:"category"`
	Message_  string            `json:"message"           refmt:"message"`
	Details_  map[string]string `json:"details,omitempty" refmt:"details,omitempty"`
}

func (e *Error) Category() interface{}      { return e.Category_ }
func (e *Error) Message() string            { return e.Message_ }
func (e *Error) Details() map[string]string { return e.Details_ }
func (e *Error) Error() string              { return e.Message_ }

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

// ToError converts any arbitrary error into the concrete rio.Error type.
// If it's an errcat.Error and already has a rio.ErrorCategory, this is
// lossless; if it's some other kind of error, we'll panic.
func ToError(err error) *Error {
	if err == nil {
		return nil
	}
	errcat.RequireErrorHasCategoryOrPanic(&err, ErrorCategory(""))
	ec := err.(errcat.Error)
	return &Error{
		ec.Category().(ErrorCategory),
		ec.Message(),
		ec.Details(),
	}
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
		if category == row.RepeatrError {
			return row.ExitCode
		}
	}
	panic(errcat.Errorf(ErrRPCBreakdown, "no exit code mapping for error category %q", category))
}

// ReboxRioError is a utility function for flipping rio.ErrorCategory into
// repeatr.ErrorCategory (or, returning ErrRPCBreakdown for unexpected cases).
func ReboxRioError(err error) error {
	category := errcat.Category(err)
	switch cat2 := category.(type) {
	case nil:
		return nil
	case rio.ErrorCategory:
		for _, row := range ErrorTable {
			if string(row.RepeatrError) == string(cat2) {
				return errcat.Recategorize(row.RepeatrError, err)
			}
		}
		return errcat.Errorf(ErrRPCBreakdown, "protocol error: unexpected error category %q from rio (error was: %s)", category, err)
	default:
		return errcat.Errorf(ErrRPCBreakdown, "protocol error: unexpected error category type %T from rio (error was: %s)", category, err)
	}
}
