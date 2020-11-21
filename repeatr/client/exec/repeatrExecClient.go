package repeatrclient

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
	"github.com/polydawn/refmt/obj/atlas"

	"github.com/polydawn/go-timeless-api"
	"github.com/polydawn/go-timeless-api/repeatr"
)

var _ repeatr.RunFunc = Run

func Run(
	ctx context.Context,
	frm api.Formula, // What formula to run.
	frmCtx repeatr.FormulaContext, // Additional information required to run (e.g. fetch and save warehouse addrs).
	input repeatr.InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor repeatr.Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (record *api.FormulaRunRecord, err error) {
	// Organize all the task specs into formulaPlus, serialize to buffer
	frmPlus := formulaPlus{frm, frmCtx}
	frmPlusBytes, err := refmt.MarshalAtlased(json.EncodeOptions{}, frmPlus, atl_formulaPlus)
	if err != nil {
		panic(err)
	}

	// prepare to exec
	//  args are almost fixed point: 'repeatr run /dev/stdin'
	cmd := exec.Command("repeatr", "--format=json", "run", "/dev/stdin")
	// feed buffer into stdin
	cmd.Stdin = bytes.NewBuffer(frmPlusBytes)
	// connect stderr straight through
	cmd.Stderr = os.Stderr
	// get stdout as pipe for rpc unmarshall to loop on
	stdout, err := cmd.StdoutPipe()
	// launch!
	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("fork repeatr: failed to start: %s", err)
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
	unmarshaller := refmt.NewUnmarshallerAtlased(json.DecodeOptions{}, stdout, repeatr.Atlas)
	var msgSlot repeatr.Event
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
				return nil, fmt.Errorf("fork repeatr: unexpected halt: %s", waitRes)
			}
			return nil, fmt.Errorf("fork repeatr: API parse error: %s", err)
		}

		// Handle it based on type.
		switch msg := msgSlot.(type) {
		case repeatr.Event_Result: // Result messages are the last ones.  Process and break.
			monitor.Send(msg)
			cmd.Wait()
			// Be careful not to return an irritating typed-nil.
			//  (msg.Error has a concrete type, and we return an interface,
			//   so this is something we have to watch out for.)
			if msg.Error == nil {
				return msg.Record, nil
			}
			return msg.Record, msg.Error
		case repeatr.Event_Log:
			monitor.Send(msg)
		case repeatr.Event_Output:
			monitor.Send(msg)
		default:
			panic(fmt.Errorf("unhandled message type %T", msg))
		}
	}
}

type (
	// formulaPlus is the concatenation of a formula and its context, and is
	// useful to serialize both {the thing to do} and {what you need to do it}
	// for sending to a repeatr process as one complete message.
	formulaPlus struct {
		Formula api.Formula
		Context repeatr.FormulaContext
	}
)

var (
	formulaPlus_AtlasEntry = atlas.BuildEntry(formulaPlus{}).StructMap().Autogenerate().Complete()

	atl_formulaPlus = atlas.MustBuild(
		formulaPlus_AtlasEntry,
		api.Formula_AtlasEntry,
		api.FilesetPackFilter_AtlasEntry,
		api.FormulaAction_AtlasEntry,
		api.FormulaUserinfo_AtlasEntry,
		api.FormulaOutputSpec_AtlasEntry,
		api.WareID_AtlasEntry,
		repeatr.FormulaContext_AtlasEntry,
	)
)
