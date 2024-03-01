package reqparser

type CustomError struct {
	message string
}

func (e CustomError) Error() string {
	return e.message
}

func NewCommonError(field, message string) CustomError {
	return CustomError{message: `"` + field + `":"` + message + `"`}
}
