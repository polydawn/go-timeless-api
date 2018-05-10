package repeatr

import (
	"fmt"

	"github.com/warpfork/go-errcat"
)

/*
	`errcat.Error` implementor with this package's `ErrorCategory` concrete category.

	This is necessary for deserialization in client APIs to yield correctly typed categories.
	Note that it does not enforce a check that the error category is from the package's
	enumerated constants when deserializing.
*/
type Error struct {
	Category_ ErrorCategory     `json:"category"          refmt:"category"`
	Message_  string            `json:"message"           refmt:"message"`
	Details_  map[string]string `json:"details,omitempty" refmt:"details,omitempty"`
}

func (e *Error) Category() interface{}      { return e.Category_ }
func (e *Error) Message() string            { return e.Message_ }
func (e *Error) Details() map[string]string { return e.Details_ }
func (e *Error) Error() string              { return e.Message_ }

/*
	Helper to set the Error field of the result message structure,
	handling type conversion checks.
*/
func (r Event_Result) SetError(err error) Event_Result {
	if err == nil {
		r.Error = nil
		return r
	}
	r.Error = &Error{}
	if e2, ok := err.(errcat.Error); ok {
		cat := errcat.Category(err)
		r.Error.Category_, ok = cat.(ErrorCategory)
		if !ok {
			r.Error.Category_ = ErrRPCBreakdown
			r.Error.Message_ = fmt.Sprintf("rpc breakdown: unfiltered error category %q -- original message: ", cat)
		}
		r.Error.Message_ += e2.Message()
		r.Error.Details_ = e2.Details()
	} else {
		r.Error.Category_ = ErrRPCBreakdown
		r.Error.Message_ = err.Error()
	}
	return r
}
