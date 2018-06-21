package repeatr

import (
	"context"
	"time"

	"go.polydawn.net/go-timeless-api"
)

type RunFunc func(
	ctx context.Context,
	boundOp api.BoundOperation, // What formula to run.
	wareSourcing api.WareSourcing, // Suggestions on where to get wares.
	input InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (*api.OperationRecord, error)

type (
	InputControl struct {
		Chan <-chan string
	}

	Monitor struct {
		Chan chan<- Event
	}
)

type (
	// A "union" type of all the kinds of event that may be generated in the
	// course of any of the functions.
	//
	// (The 'Result' message seen on the wire, but converted into returns;
	// it is never sent to the Monitor.Chan.)
	Event interface {
		_Event()
	}

	// Logs from repeatr code.
	// May include logs proxied up from rio.
	Event_Log struct {
		Time   time.Time   `refmt:"t"`
		Level  int8        `refmt:"lvl"`
		Msg    string      `refmt:"msg"`
		Detail [][2]string `refmt:"detail,omitempty"`
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
		Record *api.OperationRecord `refmt:"runRecord,omitEmpty"`
		Error  error                `refmt:",omitEmpty"`
	}
)

func (Event_Log) _Event()    {}
func (Event_Output) _Event() {}
func (Event_Result) _Event() {}
