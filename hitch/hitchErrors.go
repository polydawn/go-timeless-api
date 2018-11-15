package hitch

type ErrorCategory string

const (
	ErrUsage         ErrorCategory = ("hitch-usage-error")
	ErrCorruptState  ErrorCategory = ("hitch-corrupt-state")  // Indicates saved state is corrupt somehow (does not parse, or fails invariant checks).
	ErrNameCollision ErrorCategory = ("hitch-name-collision") // Indicates some mutation could not be performed because it tried to add data under some name that's already used.
)

type LookupError string

const (
	ErrNoSuchCatalog LookupError = ("no-such-catalog")
	ErrNoSuchRelease LookupError = ("no-such-release")
	ErrNoSuchItem    LookupError = ("no-such-item")
)
