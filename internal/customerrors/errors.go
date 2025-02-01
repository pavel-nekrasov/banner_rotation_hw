package customerrors

import "fmt"

type (
	NotFound struct {
		Message string
	}

	ParamError struct {
		Param string
		Err   error
	}

	ValidationError struct {
		Field string
		Err   error
	}
)

func (e NotFound) Error() string {
	return e.Message
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

func (v ParamError) Error() string {
	return fmt.Sprintf("%s: %s", v.Param, v.Err.Error())
}
