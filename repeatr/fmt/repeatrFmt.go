// repeatrfmt contains translators for writing repeatr.Event to an io.Writer,
// in both human-readable and API-friendly variants.
package repeatrfmt

import (
	"github.com/polydawn/go-timeless-api/repeatr"
)

type Printer interface {
	PrintLog(repeatr.Event_Log)
	PrintOutput(repeatr.Event_Output)
	PrintResult(repeatr.Event_Result)
}

func ServeMonitor(p Printer) (mon repeatr.Monitor, waitCh <-chan struct{}) {
	ch := make(chan repeatr.Event)
	sigCh := make(chan struct{})
	go serveMonitor(p, ch, sigCh)
	return repeatr.Monitor{ch}, sigCh
}

func serveMonitor(p Printer, evtCh <-chan repeatr.Event, doneCh chan<- struct{}) {
	for {
		evt, ok := <-evtCh
		if !ok {
			close(doneCh)
			return
		}
		switch evt2 := evt.(type) {
		case repeatr.Event_Log:
			p.PrintLog(evt2)
		case repeatr.Event_Output:
			p.PrintOutput(evt2)
		case repeatr.Event_Result:
			p.PrintResult(evt2)
			close(doneCh)
			return
		}
	}
}
