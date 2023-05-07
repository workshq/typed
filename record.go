package typed

import (
	"reflect"

	"github.com/valyala/fastjson"
)

type iRecord interface {
	Typed
	anyItem() any
	saveItems(*map[string]reflect.Value) error
}

// Record Object is a mixed bag of properties.
// Keys are strings and value is user defined.
type Record[T any] struct {
  typedShared
	Items map[string]T
}

func (*Record[T]) Type() Type {
	return TypeRecord
}

func (s *Record[T]) anyItem() any {
	var item T
	return &item
}

func (s *Record[T]) parseValue(val *fastjson.Value) error {
	return nil
}

func (s *Record[T]) saveItems(val *map[string]reflect.Value) error {
	// logger.Log.Println("saving", *val)
  s.Items = make(map[string]T, len(*val))
	for key, v := range *val {
    ttype := v.Elem().Interface().(T)
		s.Items[key] = ttype
	}
  return nil
}

func NewRecord[T any](items map[string]T) Record[T] {
	// return New[*Record[T]](items)
	return Record[T]{Items: items}
}