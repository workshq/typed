package typed

import "fmt"

type Type int

const (
	TypeObject Type = iota
	TypeArray
	TypeString
	TypeNumber
	TypeBoolean
	TypeRecord
  TypeOptional
)

// Returns a string representation of Type
func (t Type) String() string {
	switch t {
	case TypeObject:
		return "object"
	case TypeArray:
		return "array"
	case TypeString:
		return "string"
	case TypeNumber:
		return "number"
	case TypeRecord:
		return "record"
	case TypeBoolean:
		return "boolean"
  case TypeOptional:
    return "optional"
	default:
		panic(fmt.Errorf("BUG: unknown Typed type: %d", t))
	}
}
