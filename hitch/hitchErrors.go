package hitch

type ErrorCategory string

const (
	ErrUsage ErrorCategory = ("hitch-usage-error")
)

type LookupError string

const (
	ErrNoSuchCatalog LookupError = ("no-such-catalog")
	ErrNoSuchRelease LookupError = ("no-such-release")
	ErrNoSuchItem    LookupError = ("no-such-item")
)
