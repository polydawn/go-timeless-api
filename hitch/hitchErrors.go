package hitch

type ErrorCategory string

const (
	ErrUsage        ErrorCategory = ("hitch-usage-error")
	ErrCorruptState ErrorCategory = ("hitch-corrupt-state") // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).

)

type LookupError string

const (
	ErrNoSuchCatalog LookupError = ("no-such-catalog")
	ErrNoSuchRelease LookupError = ("no-such-release")
	ErrNoSuchItem    LookupError = ("no-such-item")
)
