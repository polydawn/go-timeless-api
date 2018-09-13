package rio

import (
	"context"
	"time"

	"go.polydawn.net/go-timeless-api"
)

type UnpackFunc func(
	ctx context.Context,
	wareID api.WareID, // What wareID to fetch for unpacking.
	path string, // Where to unpack the fileset (absolute path).
	filters api.FilesetUnpackFilter, // Optionally: filters we should apply while unpacking.
	placementMode PlacementMode, // Optionally: a placement mode specifying how the files should be put in the target path.  (Default is "copy".)
	fetchFrom []api.WarehouseLocation, // Suggestions on where to get wares.
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type PackFunc func(
	ctx context.Context,
	packType api.PackType, // The name of pack format to use.  (Most PackFunc impls support exactly one; a demux impl exists, and can route based on this string.)
	path string, // The fileset to scan and pack (absolute path).
	filters api.FilesetPackFilter, // Optionally: filters we should apply while unpacking.
	saveTo api.WarehouseLocation, // Warehouse to save into (or blank to just scan).
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type ScanFunc func(
	ctx context.Context,
	packType api.PackType, // The name of pack format.
	filters api.FilesetUnpackFilter, // Optionally: filters we should apply as if we were unpacking.
	placementMode PlacementMode, // For scanning only "None" (cache; the default) and "Direct" (don't cache) are valid.
	addr api.WarehouseLocation, // The *one* warehouse to fetch from.  Must be a monowarehouse (not a CA-mode).
	monitor Monitor, // Optionally: callbacks for progress monitoring.
) (api.WareID, error)

type MirrorFunc func(
	ctx context.Context,
	wareID api.WareID, // What wareID to mirror.
	saveTo api.WarehouseLocation, // Warehouse to ensure the ware is mirrored into.
	fetchFrom []api.WarehouseLocation, // Warehouses we can try to fetch from.
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

type (
	// Monitor hold a channel which is used for event reporting.
	// Logs will be sent to the channel if one is provided.
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

	// Logs from rio code.
	Event_Log struct {
		Time   time.Time   `refmt:"t"`
		Level  LogLevel    `refmt:"lvl"`
		Msg    string      `refmt:"msg"`
		Detail [][2]string `refmt:"detail,omitempty"`
	}

	//
	// Notifications about progress updates.
	//
	// Imagine it being used to draw the following:
	//
	//		Frobnozing (145/290kb): [=====>    ]  50%
	//
	// The 'totalProg' and 'totalWork' ints are expected to be a percentage;
	// when they equal, a "done" state should be up next.
	// A value of 'totalProg' greater than 'totalWork' is nonsensical.
	//
	// The 'phase' and 'desc' args are freetext;
	// Typically, 'phase' will remain the same for many calls in a row, while
	// 'desc' is used to communicate a more specific contextual info
	// than the 'total*' ints and like the ints may likely change on each call.
	Event_Progress struct {
		Phase, Desc          string
		TotalProg, TotalWork int
	}

	// Final results.  (Also converted into returns.)
	Event_Result struct {
		WareID api.WareID `refmt:",omitEmpty"`
		Error  error      `refmt:",omitEmpty"`
	}
)

func (Event_Log) _Event()      {}
func (Event_Progress) _Event() {}
func (Event_Result) _Event()   {}

type LogLevel int8

const (
	LogError LogLevel = 4 // Error log lines, if used, mean the program is on its way to exiting non-zero.  If used more than once, all but the first are other serious failures to clean up gracefully.
	LogWarn  LogLevel = 3 // Warning logs are for systems which have failed, but in acceptable ways; for example, a warehouse that's not online (but a fallback is, so overall we proceeded happily).
	LogInfo  LogLevel = 2 // Info logs are statements about control flow; for example, which warehouses have been tried in what order.
	LogDebug LogLevel = 1 // Debug logs are off by default.  They may get down to the resolution of called per-file in a transmat, for example.
)
