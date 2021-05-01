package domain

type Validatable interface {
	Errors() map[string][]string
	AddError(name string, msg string)
	ClearErrors()
	HasError() bool
}
