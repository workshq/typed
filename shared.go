package typed

import "github.com/valyala/fastjson"

// Common default methods for all Typed
type typedShared struct {
	val *fastjson.Value
}

func (s *typedShared) getVal() *fastjson.Value {
	return s.val
}

func (s *typedShared) setVal(val *fastjson.Value) {
	s.val = val
}

func (*typedShared) Check() error {
	return nil
}
