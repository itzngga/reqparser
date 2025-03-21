package reqparser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func (p *Parser[T]) parseAll() (v T, err error) {
	switch p.t {
	// parse body type
	case "body":
		err := p.c.BodyParser(&v)
		if err != nil {
			return v, err
		}

		return v, nil
	// parse query type
	case "query":
		err := p.c.QueryParser(&v)
		if err != nil {
			return v, err
		}

		return v, nil
	// parse keyform type
	case "keyform":
		val := p.c.FormValue(p.k)
		if p.r && val == "" {
			return v, NewCommonError(p.toSnakeCase(p.k), "NOT_BLANK")
		}

		if !p.r && val == "" {
			return
		}

		return p.parseType(val)
	// parse keyquery type
	case "keyquery":
		val := p.c.Query(p.k)
		if p.r && val == "" {
			return v, NewCommonError(p.toSnakeCase(p.k), "NOT_BLANK")
		}

		if !p.r && val == "" {
			return
		}

		return p.parseType(val)
	// parse params type
	case "params":
		val := p.c.Params(p.k)
		if p.r && val == "" {
			return v, NewCommonError(p.toSnakeCase(p.k), "NOT_BLANK")
		}

		if !p.r && val == "" {
			return
		}

		return p.parseType(val)
	// parse fileform type
	case "fileform":
		file, err := SaveFileToStorage(p.c, p.k, p.r)
		if err != nil {
			return v, err
		}

		return p.parseType(file)
	default:
		return v, errors.New("error: unsupported request type")
	}
}

func (p *Parser[T]) parseType(val string) (v T, err error) {
	switch any(v).(type) {
	case int:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return v, NewCommonError(p.toSnakeCase(p.k), "MUST_NUMBER")
		}

		var valint any = int(intVal)
		if valint, ok := valint.(T); ok {
			return valint, nil
		}
	case string:
		var vstr any = val
		if vstr, ok := vstr.(T); ok {
			return vstr, nil
		}
	case int8:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return v, NewCommonError(p.toSnakeCase(p.k), "MUST_NUMBER")
		}

		var valint any = int8(intVal)
		if valint, ok := valint.(T); ok {
			return valint, nil
		}
	case int64:
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return v, NewCommonError(p.toSnakeCase(p.k), "MUST_NUMBER")
		}

		var valint any = intVal
		if valint, ok := valint.(T); ok {
			return valint, nil
		}
	case bool:
		intVal, err := strconv.ParseBool(val)
		if err != nil {
			return v, NewCommonError(p.toSnakeCase(p.k), "NOT_VALID")
		}

		var valint any = intVal
		if valint, ok := valint.(T); ok {
			return valint, nil
		}
	case float64:
		intVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return v, NewCommonError(p.toSnakeCase(p.k), "MUST_NUMBER")
		}

		var valint any = intVal
		if valint, ok := valint.(T); ok {
			return valint, nil
		}
	default:
		return v, errors.New("error: unsupported parser type")
	}

	return v, nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func (p *Parser[T]) toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
