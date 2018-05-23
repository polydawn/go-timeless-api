package hitch

type ErrorCategory string

// REVIEW: funcs.LookupError should maybe be replaced with these hitch error categories.

const (
	ErrUsage                 = ErrorCategory("hitch-usage-error")
	ErrModuleCatalogNotExist = ErrorCategory("hitch-modulecatalog-not-exist")
)
