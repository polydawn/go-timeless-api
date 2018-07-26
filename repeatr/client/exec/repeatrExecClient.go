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
	"github.com/polydawn/refmt/obj/atlas/common"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/repeatr"
)

var _ repeatr.RunFunc = Run

func Run(
	ctx context.Context,
	boundOp api.BoundOperation, // What formula to run.
	wareSourcing api.WareSourcing, // Suggestions on where to get wares; note only the WareID index will be honored, so flip everything to that before calling.
	wareStaging api.WareStaging, // Instructions on where to store output wares.
	input repeatr.InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor repeatr.Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (record *api.OperationRecord, err error) {
	// traverse operation, flipping inputs and outputs to legacy format
	frm := formula{
		Inputs:  make(map[api.AbsPath]api.WareID),
		Action:  boundOp.Action,
		Outputs: make(map[api.AbsPath]outputSpec),
	}
	for slotRef, pth := range boundOp.Inputs {
		frm.Inputs[pth] = boundOp.InputPins[slotRef]
		if frm.Inputs[pth] == (api.WareID{}) {
			panic("missing pins in op")
		}
	}
	outputReverseMap := map[api.AbsPath]api.SlotName{}
	for slotName, pth := range boundOp.Outputs {
		frm.Outputs[pth] = outputSpec{"tar"}
		outputReverseMap[pth] = slotName
	}

	// traverse operation, convert wareSourcing to formulaContext
	//  (formulaContext will be deprecated and eventually replaced with wareSourcing)
	frmCtx := formulaContext{
		FetchUrls: make(map[api.AbsPath][]api.WarehouseLocation),
		SaveUrls:  make(map[api.AbsPath]api.WarehouseLocation),
	}
	for slotRef, pth := range boundOp.Inputs {
		frmCtx.FetchUrls[pth] = wareSourcing.ByWare[boundOp.InputPins[slotRef]]
	}
	for _, pth := range boundOp.Outputs {
		frmCtx.SaveUrls[pth] = wareStaging.ByPackType[frm.Outputs[pth].PackType]
	}

	// flip all those into formulaPlus, serialize to buffer
	frmPlus := formulaPlus{frm, frmCtx}
	frmPlusBytes, err := refmt.MarshalAtlased(json.EncodeOptions{}, frmPlus, atl_formula)
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
	unmarshaller := refmt.NewUnmarshallerAtlased(json.DecodeOptions{}, stdout, atl_rpc)
	var msgSlot event
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
				return nil, fmt.Errorf("fork repeatr: unexpected halt: %s", err)
			}
			return nil, fmt.Errorf("fork repeatr: API parse error: %s", err)
		}

		// Handle it based on type.
		switch msg := msgSlot.(type) {
		case event_Result: // Result messages are the last ones.  Process and break.
			if msg.Error != nil {
				record, err = nil, err
			} else {
				record, err = &api.OperationRecord{}, nil
				record.Guid = msg.Record.Guid
				record.Time = msg.Record.Time
				record.ExitCode = msg.Record.ExitCode
				record.Results = make(map[api.SlotName]api.WareID)
				for pth, wareID := range msg.Record.Results {
					record.Results[outputReverseMap[pth]] = wareID
				}
			}
			monitor.Send(repeatr.Event_Result{record, err})
			cmd.Wait()
			return
		case event_Log:
			monitor.Send(repeatr.Event_Log(msg))
		case event_Output:
			monitor.Send(repeatr.Event_Output(msg))
		default:
			panic(fmt.Errorf("unhandled message type %T", msg))
		}
	}
}

type (
	formulaPlus struct {
		Formula formula
		Context formulaContext
	}

	formula struct {
		Inputs  map[api.AbsPath]api.WareID
		Action  api.OpAction // "coincidentally" identical to the modern spec so we can just reuse.
		Outputs map[api.AbsPath]outputSpec
	}

	formulaContext struct {
		FetchUrls map[api.AbsPath][]api.WarehouseLocation
		SaveUrls  map[api.AbsPath]api.WarehouseLocation
	}

	outputSpec struct {
		PackType api.PackType `refmt:"packtype"`

		// filter support currently skipped.
	}

	event interface{}

	event_Result struct {
		Record *runRecord `refmt:"runRecord,omitEmpty"`
		Error  error      `refmt:",omitEmpty"`
	}

	event_Log struct {
		Time   time.Time   `refmt:"t"`
		Level  byte        `refmt:"lvl"`
		Msg    string      `refmt:"msg"`
		Detail [][2]string `refmt:"detail,omitempty"`
	}

	// Output from the contained process (stdout/stderr conjoined).
	// Stderr/stdout are conjoined so their ordering does not slip.
	// There is no guarantee of buffering (especially not line buffering);
	// in other words, `printch('.')` may indeed flush.
	event_Output struct {
		Time time.Time `refmt:"t"`
		Msg  string    `refmt:"msg"`
	}

	runRecord struct {
		Guid     string                     // random number, presumed globally unique.
		Time     int64                      // time at start of build.
		ExitCode int                        // exit code of the contained process.
		Results  map[api.AbsPath]api.WareID // wares produced by the run!
	}
)

var (
	formulaPlus_AtlasEntry    = atlas.BuildEntry(formulaPlus{}).StructMap().Autogenerate().Complete()
	formula_AtlasEntry        = atlas.BuildEntry(formula{}).StructMap().Autogenerate().Complete()
	formulaContext_AtlasEntry = atlas.BuildEntry(formulaContext{}).StructMap().Autogenerate().Complete()
	outputSpec_AtlasEntry     = atlas.BuildEntry(outputSpec{}).StructMap().Autogenerate().Complete()
	event_AtlasEntry          = atlas.BuildEntry((*event)(nil)).KeyedUnion().Of(map[string]*atlas.AtlasEntry{
		"result": event_Result_AtlasEntry,
		"log":    event_Log_AtlasEntry,
		"txt":    event_Output_AtlasEntry,
	})
	event_Result_AtlasEntry = atlas.BuildEntry(event_Result{}).StructMap().Autogenerate().Complete()
	event_Log_AtlasEntry    = atlas.BuildEntry(event_Log{}).StructMap().Autogenerate().Complete()
	event_Output_AtlasEntry = atlas.BuildEntry(event_Output{}).StructMap().Autogenerate().Complete()
	runRecord_AtlasEntry    = atlas.BuildEntry(runRecord{}).StructMap().
				Autogenerate().
				IgnoreKey("formulaID").
				IgnoreKey("hostname").
				Complete()
)

var (
	atl_formula = atlas.MustBuild(
		formulaPlus_AtlasEntry,
		formula_AtlasEntry,
		formulaContext_AtlasEntry,
		outputSpec_AtlasEntry,
		api.WareID_AtlasEntry,
		api.OpAction_AtlasEntry,
	)
	atl_rpc = atlas.MustBuild(
		runRecord_AtlasEntry,
		event_AtlasEntry,
		commonatlases.Time_AsUnixInt,
		api.WareID_AtlasEntry,
	)
)
