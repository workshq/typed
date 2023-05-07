package typed

import "github.com/valyala/fastjson"

// Common default methods for all Typed
type typedShared struct {
	jsonVal *fastjson.Value
}

func (s *typedShared) getJsonVal() *fastjson.Value {
	return s.jsonVal
}

func (s *typedShared) setJsonVal(val *fastjson.Value) {
	s.jsonVal = val
}

func (*typedShared) Check() error {
	return nil
}
