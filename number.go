package typed

import (
	"fmt"

	"github.com/valyala/fastjson"
	"golang.org/x/exp/constraints"
)

type NumberType interface {
	constraints.Integer | constraints.Float
}

type Number[T NumberType] struct {
	typedShared
	RawVal T
}

func (*Number[T]) Type() Type {
	return TypeNumber
}

func (s *Number[T]) Value() T {
	return s.RawVal
}

func (s *Number[T]) parseValue(val *fastjson.Value) error {
	raw, err := parseNumber[T](val)
	if err != nil {
		return err
	}
	s.RawVal = raw

	return nil
}

func parseNumber[T NumberType](jsonVal *fastjson.Value) (T, error) {
	var num T
	var err error
	// Ordered by most common (guessing)
	if checkImplements[int](num) ||
		checkImplements[int64](num) ||
		checkImplements[int32](num) ||
		checkImplements[int16](num) ||
		checkImplements[int8](num) {
		val, e := jsonVal.Int()
		err = e
		num = T(val)
	} else if checkImplements[float64](num) || checkImplements[float32](num) {
		val, e := jsonVal.Float64()
		err = e
		num = T(val)
	} else if checkImplements[uint](num) ||
		checkImplements[uint64](num) ||
		checkImplements[uint32](num) ||
		checkImplements[uint16](num) ||
		checkImplements[uint8](num) {
		val, e := jsonVal.Uint()
		err = e
		num = T(val)
	} else {
		return num, fmt.Errorf("unsupported number type %T", num)
	}

	if err != nil {
		return num, err
	}
	return num, nil
}

// Switch to type checks once available
// See: https://github.com/golang/go/issues/45380
//
// Using runtime checks for now
// From: https://stackoverflow.com/a/75710009/6635914
func checkImplements[I any, T any](check T) bool {
	_, ok := interface{}(check).(I)
	return ok
}

func NewNumber[T NumberType](num T) Number[T] {
	// return *New[*Number[T]](num)
	return Number[T]{RawVal: num}
}