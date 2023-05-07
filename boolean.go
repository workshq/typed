package typed

import "github.com/valyala/fastjson"

type Boolean struct {
	typedShared
	val bool
}

func (*Boolean) Type() Type {
	return TypeBoolean
}

func (s *Boolean) Value() bool {
	return s.val
}

func (s *Boolean) parseValue(val *fastjson.Value) error {
	b, err := val.Bool()
	if err != nil {
		return err
	}
	// TODO: setJsonVal and remove `setJsonVal` from typedShared
	// logger.Log.Println("PARSING", val)
	s.val = b
	return nil
}

func NewBoolean(val bool) Boolean {
	// return *New[*Boolean](val)
	return Boolean{val: val}
}