package typed

import "github.com/valyala/fastjson"

type Boolean struct {
	// Typed
	RawVal bool
}

func (*Boolean) Type() Type {
	return TypeBoolean
}

func (*Boolean) Check() error {
	return nil
}

func (s *Boolean) Value() bool {
	return s.RawVal
}

func (s *Boolean) parseValue(val *fastjson.Value) error {
	b, err := val.Bool()
	if err != nil {
		return err
	}
	// logger.Log.Println("PARSING", val)
	s.RawVal = b
	return nil
}
