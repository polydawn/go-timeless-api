package repeatrfmt

import (
	"io"

	"github.com/polydawn/refmt/json"

	"github.com/polydawn/go-timeless-api/repeatr"
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
	box := repeatr.Event(evt)
	if err := json.NewMarshallerAtlased(p.stdout, json.EncodeOptions{}, repeatr.Atlas).Marshal(&box); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}

func (p JsonPrinter) PrintOutput(evt repeatr.Event_Output) {
	box := repeatr.Event(evt)
	if err := json.NewMarshallerAtlased(p.stdout, json.EncodeOptions{}, repeatr.Atlas).Marshal(&box); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}

func (p JsonPrinter) PrintResult(evt repeatr.Event_Result) {
	box := repeatr.Event(evt)
	if err := json.NewMarshallerAtlased(p.stdout, jsonPrettyOptions, repeatr.Atlas).Marshal(&box); err != nil {
		panic(err)
	}
	p.stdout.Write([]byte{'\n'})
}
