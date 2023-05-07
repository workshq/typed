package typed

import (
	"github.com/valyala/fastjson"
)

type String struct {
	typedShared
	RawVal string
}

func (*String) Check() error {
	return nil
}

func (*String) Type() Type {
	return TypeString
}

func (s *String) Value() string {
	return s.RawVal
}

func (s *String) parseValue(val *fastjson.Value) error {
	bytes, err := val.StringBytes()
	if err != nil {
		return err
	}
  data := string(bytes)
	// logger.Log.Printf("PARSING %s %T", data, data)
	s.RawVal = data
	return nil
}
