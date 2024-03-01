package reqparser

import (
	"errors"
	"mime/multipart"
)

type FiberInterface interface {
	BodyParser(out interface{}) error
	QueryParser(out interface{}) error
	FormValue(key string, defaultValue ...string) string
	Params(key string, defaultValue ...string) string
	Query(key string, defaultValue ...string) string
	FormFile(key string) (*multipart.FileHeader, error)
}

type RulesFunc func() RuleSet

type RuleSet struct {
	t string // type of parser
	k string // key body/query/params
	r *bool  // required or not
}

type Parser[T any] struct { // target interface
	c FiberInterface // fiber context
	t string         // type of parser
	k string         // key body/query/params
	r bool           // required or not
}

func New[T any](c FiberInterface) *Parser[T] {
	return &Parser[T]{
		c: c,
	}
}

func (p *Parser[T]) Parse(rules ...RulesFunc) (v T, err error) {
	if len(rules) <= 0 {
		return v, errors.New("error: please provide at least one request parser rule")
	}

	for _, rule := range rules {
		rule := rule()
		if rule.t != "" {
			p.t = rule.t
		}
		if rule.k != "" {
			p.k = rule.k
		}
		if rule.r != nil {
			p.r = *rule.r
		}
	}

	return p.parseAll()
}

func BodyParser() RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "body",
		}
	}
}

func Form(key string, required bool) RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "keyform",
			k: key,
			r: &required,
		}
	}
}

func QueryParser() RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "query",
		}
	}
}

func Query(key string, required bool) RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "keyquery",
			k: key,
			r: &required,
		}
	}
}

func Params(key string, required bool) RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "params",
			k: key,
			r: &required,
		}
	}
}

func File(key string, required bool) RulesFunc {
	return func() RuleSet {
		return RuleSet{
			t: "fileform",
			k: key,
			r: &required,
		}
	}
}
