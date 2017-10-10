/*
	Interfaces of rio commands.

	The heuristic for the rio callable library API is that essentially
	all information must be racked up in the call already: the assumption
	is that the side doing the real work might not share a filesystem with
	you (well, in rio's case, it probably does!  but it might be a subset,
	translated through chroots and bind mounts), doesn't share env vars, etc.
	So, general rule of thumb: the caller is going to have already handled
	all config loading and parsing, and those objects are params in this funcs.
*/
package rio

import (
	"context"
	"time"

	"go.polydawn.net/go-timeless-api"
)

type UnpackFunc func(
	ctx context.Context, // Long-running call.  Cancellable.
	wareID api.WareID, // What wareID to fetch for unpacking.
	path string, // Where to unpack the fileset (absolute path).
	filters api.FilesetFilters, // Optionally: filters we should apply while unpacking.
	placementMode PlacementMode, // Optionally: a placement mode specifying how the files should be put in the target path.  (Default is "copy".)
	warehouses []api.WarehouseAddr, // Warehouses we can try to fetch from.
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type PackFunc func(
	ctx context.Context, // Long-running call.  Cancellable.
	packType api.PackType, // The name of pack format to use.  (Most PackFunc impls support exactly one; a demux impl exists, and can route based on this string.)
	path string, // The fileset to scan and pack (absolute path).
	filters api.FilesetFilters, // Optionally: filters we should apply while unpacking.
	warehouse api.WarehouseAddr, // Warehouse to save into (or blank to just scan).
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type ScanFunc func(
	ctx context.Context, // Long-running call.  Cancellable.
	packType api.PackType, // The name of pack format.
	filters api.FilesetFilters, // Optionally: filters we should apply while unpacking.
	placementMode PlacementMode, // For scanning only "None" (cache; the default) and "Direct" (don't cache) are valid.
	addr api.WarehouseAddr, // The *one* warehouse to fetch from.  Must be a monowarehouse (not a CA-mode).
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type MirrorFunc func(
	ctx context.Context, // Long-running call.  Cancellable.
	wareID api.WareID, // What wareID to mirror.
	target api.WarehouseAddr, // Warehouse to ensure the ware is mirrored into.
	sources []api.WarehouseAddr, // Warehouses we can try to fetch from.
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type PlacementMode string

const (
	// "none" mode instructs unpack not to place the files at all -- but it
	// still updates the fileset cache.  So, you can use this to warm up the
	// cache.  The target path argument to unpack will be ignored.
	Placement_None PlacementMode = "none"
	// "copy" mode -- the default -- instructs unpack to use the cache of
	// already unpacked filesets (unpacking there in case of cache miss), and
	// then place the files in their final location by a plain file copy.
	Placement_Copy PlacementMode = "copy"
	// "mount" mode instructs unpack to use the fileset cache (same as "copy"),
	// then place the files in their final location by using some sort of mount.
	// Whether "mount" means "bind", "aufs", "overlayfs", etc is left to
	// interpretation, but regardless it A) should be faster than "copy" and
	// B) since it's a mount, may be slightly harder to rmdir :)
	Placement_Mount PlacementMode = "mount"
	// "direct" mode instructs unpack to skip the cache and work directly in
	// the target path.  (It will still fall back to copy mode if the requested
	// ware is already in the fileset cache, "direct" is the one mode that
	// will not *populate* the fileset cache if empty.)
	Placement_Direct PlacementMode = "direct"
)

/*
	Monitoring configuration structs, and message types used.
*/
type (
	// REVIEW ... it's rather generalizing to use the same monitor and event union
	//  for all these different functions, isn't it?

	/*
		Configuration for what intermediate progress reports a process should send,
		and slot for the channel the caller wishes them to be sent to.
	*/
	Monitor struct {
		// FUTURE: may add options for how many things we'd like to be sent to us

		// Channel to which events will be sent as the process proceeds.
		// The channel will be closed when the process is done or cancelled.
		// A nil channel will disable all intermediate progress reporting.
		Chan chan<- Event
	}

	/*
		A "union" type of all the kinds of event that may be generated in the
		course of any of the functions.

		The "Result" message is never sent to Monitor.Chan --
		its values are converted into the function returns --
		but *is* seen in the serial form on the wire.

		(This type may be replaced by an interface in the future when the refmt
		library's union message support becomes available.)
	*/
	Event struct {
		Log      *Event_Log      `refmt:"log,omitempty"`
		Progress *Event_Progress `refmt:"prog,omitempty"`
		Result   *Event_Result   `refmt:"result,omitempty"`
	}

	/*
		Logs of major events.  (Progress bars are a separate thing; this is
		for e.g. "tried warehouse; failed" or "warehouse %x: ware %x 404".)
	*/
	Event_Log struct {
		Time   time.Time   `refmt:"t"`
		Level  LogLevel    `refmt:"lvl"`
		Msg    string      `refmt:"msg"`
		Detail [][2]string `refmt:"detail,omitempty"`
	}

	/*
		Notifications about progress updates.

		Imagine it being used to draw the following:

			Frobnozing (145/290kb): [=====>    ]  50%

		The 'totalProg' and 'totalWork' ints are expected to be a percentage;
		when they equal, a "done" state should be up next.
		A value of 'totalProg' greater than 'totalWork' is nonsensical.

		The 'phase' and 'desc' args are freetext;
		Typically, 'phase' will remain the same for many calls in a row, while
		'desc' is used to communicate a more specific contextual info
		than the 'total*' ints and like the ints may likely change on each call.
	*/
	Event_Progress struct {
		Phase, Desc          string
		TotalProg, TotalWork int
	}

	Event_Result struct {
		WareID api.WareID `refmt:",omitempty"`
		Error  *Error     `refmt:",omitempty"`
	}
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
	ErrNotImplemented       = ErrorCategory("rio-not-implemented")       // The operation is not implemented, PRs welcome
	ErrRPCBreakdown         = ErrorCategory("rio-rpc-breakdown")         // Raised when running a remote rio process and the control channel is lost, the process fails to start, or unrecognized messages are received.
)

type LogLevel int8

const (
	LogError LogLevel = 4 // Error log lines, if used, mean the program is on its way to exiting non-zero.  If used more than once, all but the first are other serious failures to clean up gracefully.
	LogWarn  LogLevel = 3 // Warning logs are for systems which have failed, but in acceptable ways; for example, a warehouse that's not online (but a fallback is, so overall we proceeded happily).
	LogInfo  LogLevel = 2 // Info logs are statements about control flow, for exmaple, which warehouses have been tried in what order.
	LogDebug LogLevel = 1 // Debug logs are off by default.  They may get down to the resolution of called per-file in a transmat, for example.
)

func (ll LogLevel) String() string {
	switch ll {
	case LogError:
		return "error"
	case LogWarn:
		return "warn"
	case LogInfo:
		return "info"
	case LogDebug:
		return "debug"
	default:
		return "invalid"
	}
}
