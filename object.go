package typed

import "github.com/valyala/fastjson"

type iObject interface {
	Typed
	anyProps() any
}

type Object[T any] struct {
  typedShared
	Props T
}

func (*Object[T]) Type() Type {
	return TypeObject
}

func (s *Object[T]) anyProps() any {
	return &s.Props
}

func (s *Object[T]) parseValue(val *fastjson.Value) error {
	return nil
}

func NewObject[T any](props T) Object[T] {
	return Object[T]{Props: props}
}
