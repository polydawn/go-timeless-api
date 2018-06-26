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

	// Monitor hold a channel which is used for event reporting.
	// Logs, contained process output, and results will all be sent to the
	// channel if one is provided.
	//
	// The same channel and monitor may be used for multiple runs.
	// The channel is not closed when the job is done; any goroutine with
	// blocking service to a channel used for a single job should return
	// after recieving an Event_Result, as it will be the final event.
	Monitor struct {
		Chan chan<- Event
	}
)

func (m Monitor) Send(evt Event) {
	if m.Chan != nil {
		m.Chan <- evt
	}
}

type (
	// A "union" type of all the kinds of event that may be generated in the
	// course of any of the functions.
	Event interface {
		_Event()
	}

	// Logs from repeatr code.
	// May include logs proxied up from rio.
	Event_Log struct {
		Time   time.Time   `refmt:"t"`
		Level  byte        `refmt:"lvl"`
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

	// Final results.  (Also converted into returns.)
	Event_Result struct {
		Record *api.OperationRecord `refmt:"runRecord,omitEmpty"`
		Error  error                `refmt:",omitEmpty"`
	}
)

func (Event_Log) _Event()    {}
func (Event_Output) _Event() {}
func (Event_Result) _Event() {}
