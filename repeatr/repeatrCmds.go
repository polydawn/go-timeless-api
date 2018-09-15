package repeatr

import (
	"context"
	"time"

	"go.polydawn.net/go-timeless-api"
)

type RunFunc func(
	ctx context.Context,
	boundOp api.Formula, // What formula to run.
	formulaContext FormulaContext, // Additional information required to run (e.g. fetch and save warehouse addrs).
	input InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (*api.FormulaRunRecord, error)

type (
	// FormulaContext bundles information which is important to being *able* to
	// run the computation, but *immaterial* to the output.  Specifically, this
	// means URLs for where to get things, and where to save them.
	//
	// This type definition is hiding in the repeatr-specific package rather than
	// escalated into the group of main API types because it's useful precisely
	// at the level of repeatr.  Lower levels, e.g. Rio, don't need these
	// groupings; and higher levels, e.g. Stellar, also do something more
	// opinionated (see: WareSourcing).
	FormulaContext struct {
		FetchUrls map[api.AbsPath][]api.WarehouseLocation
		SaveUrls  map[api.AbsPath]api.WarehouseLocation
	}
)

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
		Level  LogLevel    `refmt:"lvl"`
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
		Record *api.FormulaRunRecord `refmt:"runRecord,omitEmpty"`
		Error  error                 `refmt:",omitEmpty"`
	}
)

func (Event_Log) _Event()    {}
func (Event_Output) _Event() {}
func (Event_Result) _Event() {}

type LogLevel int8

const (
	LogError LogLevel = 4 // Error log lines, if used, mean the program is on its way to exiting non-zero.  If used more than once, all but the first are other serious failures to clean up gracefully.
	LogWarn  LogLevel = 3 // Warning logs are for systems which have failed, but in acceptable ways; for example, a warehouse that's not online (but a fallback is, so overall we proceeded happily).
	LogInfo  LogLevel = 2 // Info logs are statements about control flow; for example, which warehouses have been tried in what order.
	LogDebug LogLevel = 1 // Debug logs are off by default.  They may get down to the resolution of called per-file in a transmat, for example.
)
