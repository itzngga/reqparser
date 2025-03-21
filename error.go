package reqparser

type CustomError struct {
	field   string
	message string
}

func (customError CustomError) Error() string {
	return customError.message
}

func (customError CustomError) Parse() map[string]string {
	str := make(map[string]string)
	str[customError.field] = customError.message
	return str
}

func NewCommonError(field, message string) CustomError {
	return CustomError{field, message}
}
