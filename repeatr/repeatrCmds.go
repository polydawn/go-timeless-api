/*
	Interfaces of repeatr commands.

	The real repeatr logic implements these;
	so do the various proxy tools (e.g. r2k8s);
	so do some mocks, which are useful for testing.
*/
package repeatr

import (
	"context"
	"time"

	"github.com/polydawn/go-errcat"

	"go.polydawn.net/go-timeless-api"
)

type RunFunc func(
	ctx context.Context, // Long-running call.  Cancellable.
	formula api.Formula, // What formula to run.  Comes complete with fetch and save paths.
	input InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (*api.RunRecord, error)

/*
	Holder for input control.  The zero value disable input.
	(Disabled input is the norm: streaming inputs are only used
	by 'twerk' mode, e.g. when you've given up on reproducible action.)
*/
type InputControl struct {
	Chan <-chan string
}

/*
	Monitoring configuration structs, and message types used.
*/
type (
	Monitor struct {
		Chan chan<- Event
	}

	// A "union" type of all the kinds of event that may be generated in the
	// course of any of the functions.
	//
	// (The 'Result' message seen on the wire, but converted into returns;
	// it is never sent to the Monitor.Chan.)
	Event struct {
		Log    *Event_Log    `refmt:"log,omitempty"`
		Output *Event_Output `refmt:"txt,omitempty"`
		Result *Event_Result `refmt:"result,omitempty"`
	}

	// Logs from repeatr code.
	Event_Log struct {
		Time  time.Time   `refmt:"t"`
		Level int         `refmt:"lvl"`
		Msg   string      `refmt:"msg"`
		Ctx   [][2]string `refmt:"ctx,omitempty"`
	}

	// Output from the contained process (stdout/stderr conjoined).
	// Stderr/stdout are conjoined so their ordering does not slip.
	// There is no guarantee of buffering (especially not line buffering);
	// in other words, `printch('.')` may indeed flush.
	Event_Output struct {
		Time time.Time `refmt:"t"`
		Msg  string    `refmt:"msg"`
	}

	// Final results.  (Converted into returns; not sent to Monitor.)
	Event_Result struct {
		RunRecord api.RunRecord `refmt:",omityEmpty"`
		Error     *Error        `refmt:",omityEmpty"`
	}
)

/*
	`errcat.Error` implementor with `repeatr.ErrorCategory` concrete category.

	This is necessary for deserialization in client APIs to yield correctly typed categories.
	Note that it does not enforce a check that the error category is from the package's
	enumerated constants when deserializing.
*/
type Error struct {
	Category_ ErrorCategory     `json:"category"          refmt:"category"`
	Message_  string            `json:"message"           refmt:"message"`
	Details_  map[string]string `json:"details,omitempty" refmt:"details,omitempty"`
}

func (e *Error) Category() interface{}      { return e.Category_ }
func (e *Error) Message() string            { return e.Message_ }
func (e *Error) Details() map[string]string { return e.Details_ }
func (e *Error) Error() string              { return e.Message_ }

/*
	Helper to set the Error field of the result message structure,
	handling type conversion checks.
*/
func (r *Event_Result) SetError(err error) {
	if err == nil {
		r.Error = nil
		return
	}
	r.Error = &Error{}
	if e2, ok := err.(errcat.Error); ok {
		r.Error.Category_ = errcat.Category(err).(ErrorCategory)
		r.Error.Message_ = e2.Message()
		r.Error.Details_ = e2.Details()
	} else {
		r.Error.Category_ = ErrRPCBreakdown // :/
		r.Error.Message_ = err.Error()
	}
}

type ErrorCategory string

const (
	ErrUsage             = ErrorCategory("repeatr-usage-error")         // Indicates some piece of user input to a command was invalid and unrunnable.
	ErrJobUnsuccessful   = ErrorCategory("repeatr-job-unsuccessful")    // Not an error -- indicates that the contained process exited nonzero.  TODO review if this needs an error category or just a reserved space in the exit code table.
	ErrJobInvalid        = ErrorCategory("repeatr-job-invalid")         // Indicates the container could not be launched because some part of its specification was invalid -- for example, the CWD requested is not a dir, or the command to exec is not an executable.  (The whole filesystem may have been necessary to set up before this can be detected.)
	ErrLocalCacheProblem = ErrorCategory("repeatr-local-cache-problem") // Indicates an error while while handling internal filesystem paths (for example, if an executor can't mkdir its workspace dirs).
	ErrExecutor          = ErrorCategory("repeatr-executor-problem")    // Indicates an error occured while launching containment or handling the child processes.  Should be seen rarely -- comes up for esotera like "out of file handles".

	ErrWarehouseUnavailable = ErrorCategory("rio-warehouse-unavailable") // The corresponding rio error halted execution.
	ErrWarehouseUnwritable  = ErrorCategory("rio-warehouse-unwritable")  // The corresponding rio error halted execution.
	ErrWareNotFound         = ErrorCategory("rio-ware-not-found")        // The corresponding rio error halted execution.
	ErrWareCorrupt          = ErrorCategory("rio-ware-corrupt")          // The corresponding rio error halted execution.
	ErrWareHashMismatch     = ErrorCategory("rio-hash-mismatch")         // The corresponding rio error halted execution.
	ErrCancelled            = ErrorCategory("rio-cancelled")             // The corresponding rio error halted execution.
	ErrRioCacheProblem      = ErrorCategory("rio-local-cache-problem")   // The corresponding rio error halted execution.
	ErrAssemblyInvalid      = ErrorCategory("rio-assembly-invalid")      // The corresponding rio error halted execution.
	ErrPackInvalid          = ErrorCategory("rio-pack-invalid")          // The corresponding rio error halted execution.
	ErrRPCBreakdown         = ErrorCategory("repeatr-rpc-breakdown")     // Raised when running a remote process and the control channel is lost, the process fails to start, or unrecognized messages are received.
)
