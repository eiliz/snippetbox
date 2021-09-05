package forms

// This type will hold validations errors with the name of the field as a the
// key in the map.
type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	messages, ok := e[field]

	if !ok {
		return ""
	}

	return messages[0]
}
