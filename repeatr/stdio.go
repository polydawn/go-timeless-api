package repeatr

import (
	"io"
)

/*
	Use this to handle `Event_Output` sent to a `repeatr.Monitor.Chan`,
	turning them into `io.Writer` calls.

	You can hand all `Event`s to this; it will no-op on the irrelevant ones.
*/
func CopyOut(evt Event, into io.Writer) error {
	if evt.Output == nil {
		return nil
	}
	_, err := into.Write([]byte(evt.Output.Msg))
	return err
}
