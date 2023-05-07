package typed

import "github.com/valyala/fastjson"

type iObject interface {
	// Typed
	anyProps() any
}

type Object[T any] struct {
  // iObject
	Props T
}

func (*Object[T]) Type() Type {
	return TypeObject
}

func (*Object[T]) Check() error {
	return nil
}

func (s *Object[T]) anyProps() any {
	return &s.Props
}

func (s *Object[T]) parseValue(val *fastjson.Value) error {
	return nil
}
