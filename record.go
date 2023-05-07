package typed

import (
	"reflect"

	"github.com/valyala/fastjson"
)

type iRecord interface {
	// Typed
	anyItem() any
	saveItems(*map[string]reflect.Value) error
}

// Record Object is a mixed bag of properties.
// Keys are strings and value is user defined.
type Record[T any] struct {
  // iRecord
	Items map[string]T
}

func (*Record[T]) Type() Type {
	return TypeRecord
}

func (*Record[T]) Check() error {
	return nil
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
