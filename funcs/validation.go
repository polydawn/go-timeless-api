package funcs

// Validatable is an interface types may implement in order to return validation
// errors in a conventional way.
type Validatable interface {
	Validate() error
}

// MustValidate runs validation on a value and will panic if the validation fails.
// It is preferable to use this vs doing the panic yourself in order to denote
// the semiotics that your call site expects that valid to already be validated.
// (In the future, it should be possible to add a build mode which no-ops this
// func like a C-style 'assert', and there should be no behavior change.)
func MustValidate(v Validatable) {
	if err := v.Validate(); err != nil {
		panic(err)
	}
}
