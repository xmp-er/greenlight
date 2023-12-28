package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}
func (v *Validator) Valid() bool { //if error map is empty then it is valid
	return len(v.Errors) == 0
}

func (v *Validator) addError(key, message string) {
	if _, isKeyExists := v.Errors[key]; isKeyExists == false {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) { //if ok flag is false the this will be added to error
	if ok == false {
		v.addError(key, message)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool { //to find value in given data which supports comparison operator
	for _, v := range permittedValues {
		if value == v {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool { //check if all are unique elements in the passed array
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(uniqueValues) == len(values)
}
