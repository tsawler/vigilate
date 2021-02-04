package forms

type errors map[string][]string

// Add implements an Add() method to add error messages for a given field to the map.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get implements an Get() method to retrieve the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
