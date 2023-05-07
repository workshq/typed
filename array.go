package typed

import (
	"reflect"

	"github.com/valyala/fastjson"
)

type iArray interface {
	Typed
	anyItem() any
	saveItems(*[]reflect.Value) error
}

type Array[T any] struct {
  typedShared
	// TODO: Add a safe `.at()` method
	Items []T
}

func (*Array[T]) Type() Type {
	return TypeArray
}

func (s *Array[T]) anyItem() any {
	var item T
	return &item
}

func (s *Array[T]) parseValue(val *fastjson.Value) error {
	return nil
}

func (s *Array[T]) saveItems(val *[]reflect.Value) error {
	// logger.Log.Println("saving", val)
	for _, v := range *val {
    ttype := v.Elem().Interface().(T)
    // g := v.(T)
		// b, err := v.(T)
		// s.Items = append(s.Items, any(v).(T))
		s.Items = append(s.Items, ttype)
	}
  return nil
}

func NewArray[T any](items []T) Array[T] {
	// return New[*Array[T]](items)
	return Array[T]{Items: items}
}