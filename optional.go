package typed

import (
	"fmt"

	"github.com/valyala/fastjson"
)

type Optional[T any] struct {
	RawVal *T
}

func (*Optional[T]) Type() Type {
	return TypeOptional
}

func (*Optional[T]) Check() error {
	return nil
}

// True if Optional Value is present
func (s *Optional[T]) IsPresent() bool {
  return s.RawVal != nil
}

// True if Optional Value is absent
func (s *Optional[T]) IsAbsent() bool {
  return s.RawVal == nil
}

func (s *Optional[T]) Expect(format string, a ...any) (*T, error) {
  if s.IsAbsent() {
    return nil, fmt.Errorf(format, a...)
  }
	return s.RawVal, nil
}

// Get the optional value. May be nil
func (s *Optional[T]) Maybe() *T {
	return s.RawVal
}

// Gets value. If absent, uses the fallback.
// The return will not be nil.
func (s *Optional[T]) OrElse(fallback T) *T {
  if s.RawVal == nil {
    return &fallback
  }
	return s.RawVal
}

func (s *Optional[T]) parseValue(val *fastjson.Value) error {
  // Return out if no value
  if val == nil {
    return nil;
  }
  t := new(T)
  s.RawVal = t
  // logger.Log.Printf("parsing %+v %T", t, t)
  typed := any(t).(Typed)
  return parse(val, typed)
}
