package repeatrfmt

import (
	"io"

	"github.com/polydawn/refmt/json"

	"go.polydawn.net/go-timeless-api/repeatr"
)

var _ Printer = JsonPrinter{}

type JsonPrinter struct{ stdout io.Writer }

func NewJsonPrinter(stdout io.Writer) *JsonPrinter {
	return &JsonPrinter{stdout}
}

var jsonPrettyOptions = json.EncodeOptions{
	Line:   []byte{'\n'},
	Indent: []byte("    "),
}

func (p JsonPrinter) PrintLog(evt repeatr.Event_Log) {
	if err := json.NewMarshallerAtlased(p.stdout, json.EncodeOptions{}, repeatr.Atlas).Marshal(evt); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}

func (p JsonPrinter) PrintOutput(evt repeatr.Event_Output) {
	if err := json.NewMarshallerAtlased(p.stdout, json.EncodeOptions{}, repeatr.Atlas).Marshal(evt); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}

func (p JsonPrinter) PrintResult(evt repeatr.Event_Result) {
	if err := json.NewMarshallerAtlased(p.stdout, jsonPrettyOptions, repeatr.Atlas).Marshal(evt); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}
