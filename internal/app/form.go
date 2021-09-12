package app

func NewPostForm(action string) Form {
	return Form{
		Action: action,
		Method: "POST",
	}
}

type Form struct {
	Errors map[string][]string `schema:"-" json:"-"`
	Action string              `schema:"-" json:"-"`
	Method string              `schema:"_method" json:"-"`
}

func (f *Form) IsUpdate() bool {
	return f.Method == "PATCH" || f.Method == "PUT"
}

func (f *Form) IsInsert() bool {
	return f.Method == "POST"
}

func (f *Form) IsCreate() bool {
	return f.IsInsert()
}

func (f *Form) AddError(name string, msg string) {
	if f.Errors == nil {
		f.Errors = make(map[string][]string)
	}

	_, exists := f.Errors[name]
	if !exists {
		f.Errors[name] = make([]string, 0, 1)
	}

	f.Errors[name] = append(f.Errors[name], msg)
}

func (f *Form) ClearErrors() {
	if f.HasError() {
		f.Errors = make(map[string][]string)
	}
}

func (f *Form) HasError() bool {
	if f.Errors == nil {
		return false
	}

	if len(f.Errors) == 0 {
		return false
	}

	return true
}
