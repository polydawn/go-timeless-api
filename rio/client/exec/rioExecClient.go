package rioclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-errcat"

	"github.com/polydawn/go-timeless-api"
	"github.com/polydawn/go-timeless-api/rio"
)

var (
	_ rio.UnpackFunc = UnpackFunc
	_ rio.PackFunc   = PackFunc
)

func UnpackFunc(
	ctx context.Context,
	wareID api.WareID,
	path string,
	filt api.FilesetUnpackFilter,
	placementMode rio.PlacementMode,
	warehouses []api.WarehouseLocation,
	monitor rio.Monitor,
) (gotWareID api.WareID, err error) {
	// Marshal args.
	args, err := UnpackArgs(wareID, path, filt, placementMode, warehouses, monitor)
	if err != nil {
		return api.WareID{}, err
	}
	// Bulk of invoking and handling process messages is shared code.
	return packOrUnpack(ctx, args, monitor)
}

func PackFunc(
	ctx context.Context,
	packType api.PackType,
	path string,
	filt api.FilesetPackFilter,
	warehouse api.WarehouseLocation,
	monitor rio.Monitor,
) (api.WareID, error) {
	// Marshal args.
	args, err := PackArgs(packType, path, filt, warehouse, monitor)
	if err != nil {
		return api.WareID{}, err
	}
	// Bulk of invoking and handling process messages is shared code.
	return packOrUnpack(ctx, args, monitor)
}

// internal implementation of message parsing for both pack and unpack.
// (they "conincidentally" have the same API.)
func packOrUnpack(
	ctx context.Context,
	args []string,
	monitor rio.Monitor,
) (api.WareID, error) {
	if monitor.Chan != nil {
		defer close(monitor.Chan)
	}

	// Spawn process.
	cmd := exec.Command("rio", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: failed to start: %s", err)
	}
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	if err = cmd.Start(); err != nil {
		return api.WareID{}, Errorf(rio.ErrRPCBreakdown, "fork rio: failed to start: %s", err)
	}

	// Set up reaction to ctx.done: send a sig to the child proc.
	//  (No, you couldn't set this up without a goroutine -- you can't select with the IO we're about to do;
	//  and No, you couldn't do it until after cmd.Start -- the Process handle doesn't exist until then.)
	go func() {
		<-ctx.Done() // FIXME goroutine leak occurs when the process ends gracefully
		cmd.Process.Signal(os.Interrupt)
		time.Sleep(100 * time.Millisecond)
		cmd.Process.Signal(os.Kill)
	}()

	// Consume stdout, converting it to Monitor.Chan sends.
	//  When exiting because the child sent its 'result' message correctly, the
	//  msgSlot will hold the final data (or error); we'll return it at the end.
	//  (We're relying on the child proc getting signal'd to close the stdout pipe
	//  and in turn release us here in case of ctx.done.)
	unmarshaller := refmt.NewUnmarshallerAtlased(json.DecodeOptions{}, stdout, rio.Atlas)
	var msgSlot rio.Event
	for {
		// Peel off a message.
		if err := unmarshaller.Unmarshal(&msgSlot); err != nil {
			if err == io.EOF {
				// In case of unexpected EOF, there must have been a panic on the other side;
				//  it'll be more informative to break here and return the error from Wait,
				//  which will include the stderr capture.
				waitRes := cmd.Wait()
				if waitRes == nil {
					waitRes = fmt.Errorf("bizarre zero exit code")
				}
				return api.WareID{}, fmt.Errorf("fork rio: unexpected halt: %s\n\tstderr follows:\n%s\n\n", waitRes, stderrBuf.String())
			}
			return api.WareID{}, fmt.Errorf("fork rio: API parse error: %s", err)
		}

		// Handle it based on type.
		switch msg := msgSlot.(type) {
		case rio.Event_Result: // Result messages are the last ones.  Process and break.
			monitor.Send(msg)
			cmd.Wait()
			// Be careful not to return an irritating typed-nil.
			//  (msg.Error has a concrete type, and we return an interface,
			//   so this is something we have to watch out for.)
			if msg.Error == nil {
				return msg.WareID, nil
			}
			return msg.WareID, msg.Error
		case rio.Event_Log:
			monitor.Send(msg)
		case rio.Event_Progress:
			monitor.Send(msg)
		default:
			panic(fmt.Errorf("unhandled message type %T", msg))
		}
	}
}
