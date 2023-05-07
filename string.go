package typed

import (
	"github.com/valyala/fastjson"
)

type String struct {
	typedShared
	Val string
}

func (*String) Type() Type {
	return TypeString
}

func (s *String) Value() string {
	return s.Val
}

func (s *String) parseValue(val *fastjson.Value) error {
	bytes, err := val.StringBytes()
	if err != nil {
		return err
	}
  data := string(bytes)
	// logger.Log.Printf("PARSING %s %T", data, data)
	s.Val = data
	return nil
}

func NewString(val string) String {
	// return New[String](val)
	return String{Val: val}
}
