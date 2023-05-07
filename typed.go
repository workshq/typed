package typed

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/valyala/fastjson"
)

// TODO: Add Library to awesome-go
// Good non-struct tag based JSON parsing
// https://github.com/avelino/awesome-go#validation
// https://awesome-go.com/validation/

type Typed interface {
	parseValue(val *fastjson.Value) error
	getJsonVal() *fastjson.Value
	setJsonVal(val *fastjson.Value)
	Type() Type
	Check() error
}

func Parse(json string, typed Typed) error {
	// logger.Log.Printf("%+v", typed)
	var p fastjson.Parser
	val, err := p.Parse(json)
	if err != nil {
		return fmt.Errorf("error parsing: %w", err)
	}
	return parse(val, typed)
}

func ParseBytes(json []byte, typed Typed) error {
	// logger.Log.Printf("%+v", typed)
	var p fastjson.Parser
	val, err := p.ParseBytes(json)
	if err != nil {
		return fmt.Errorf("error parsing: %w", err)
	}
	return parse(val, typed)
}

// TODO: Format object and point to error if error
//   - I could wrap fastjson.Value in a custom type
//   - This type could store the tree structure / parent relationship.
//   - I could use this tree to create a representation of the JSON
//     and point to the correct value.
//   - Should be user called since recreation would probably be slowish.
//   - Would be nice if you could style this as JSON or YAML so it's
//     more accurate to original imported data.
//   - It's just a representation not the actually data so the format could be anything.
func parse(val *fastjson.Value, typed Typed) error {
	t := typed.Type()

	var parseErr error
	switch t {
	case TypeObject:
		parseErr = parseObject(val, typed.(iObject))
	case TypeArray:
		parseErr = parseArray(val, typed.(iArray))
	case TypeRecord:
		parseErr = parseRecord(val, typed.(iRecord))
	case TypeString, TypeNumber, TypeBoolean, TypeOptional:
		typed.setJsonVal(val)
		parseErr = typed.parseValue(val)
	default:
		return fmt.Errorf("unsupported type %s from %T", t.String(), typed)
	}

	if parseErr != nil {
		return parseErr
	}
	// TODO: Validate and collect all errors, don't fail after the first.
	// Similar to this package, should collect errors
	// https://github.com/gobuffalo/validate
	// Probably need to wrap fastjson.Object in the custom Value (with Parent connections)
	// I can add an error collector to the struct.
	return typed.Check()
}

func parseObject(raw *fastjson.Value, typed iObject) error {
	val, err := raw.Object()
	if err != nil {
		return err
	}
	typed.setJsonVal(raw)
	props := typed.anyProps()
	// logger.Log.Printf("Props %+v", props)
	r := reflect.ValueOf(props).Elem()
	fields := make(map[string]reflect.StructField, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		field := r.Type().Field(i)
		fields[field.Name] = field
	}
	// logger.Log.Printf("%+v", fields)
	// print.PrettyPrint(fields)

	for _, field := range fields {
		fieldName := field.Name
		fieldValue := r.FieldByName(fieldName)
		// Trying to support pointers
		// if fieldValue.Kind() == reflect.Ptr {
		//   fieldValue = fieldValue.Elem()
		// }
		// logger.Log.Printf("Field %s %T", fieldName, fieldValue)
		keyName := field.Tag.Get("name")
		// Check if empty
		if keyName == "" {
			keyName = lowerInitial(fieldName)
		}
		data := val.Get(keyName)

		ttype, ok := fieldValue.Addr().Interface().(Typed)
		if !ok {
			return fmt.Errorf("unable to convert to Typed interface %s", fieldName)
		}
		// Throw error if key is absent if required
		// Checking type first is ~60 microseconds faster
		// TODO: Show a nice inline pointer to where the error is in the JSON
		if ttype.Type() != TypeOptional && data == nil {
			return fmt.Errorf("unable to find required field %s with key \"%s\"", fieldName, keyName)
		}
		// logger.Log.Println("convert", ttype, ok)
		err := parse(data, ttype)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseArray(raw *fastjson.Value, typed iArray) error {
	val, err := raw.Array()
	if err != nil {
		return err
	}
	typed.setJsonVal(raw)
	item := typed.anyItem()
	// logger.Log.Printf("Array %+v", item)
	// items := typed.allItems()
	// logger.Log.Printf("all items %+v", items)

	innerVal := reflect.ValueOf(item).Elem()
	// logger.Log.Printf("innerVal %+v", innerVal)
	// innerVal := reflect.ValueOf(innerType)
	// logger.Log.Printf("innerVal %+v", innerVal)
	// r.Field(0)
	// logger.Log.Println(innerVal, innerVal.Kind())
	_, ok := innerVal.Addr().Interface().(Typed)
	if !ok {
		return fmt.Errorf("unable to convert to Typed interface %T", innerVal)
	}
	// logger.Log.Println("type", ttype, ttype.Type())
	// typed.parseValue(val)

	// TODO: Prealloc array. It seems to fail when done.
	var data []reflect.Value
	// data := make([]reflect.Value, 1)
	for _, v := range val {
		// Create a new Typed container to hold value
		n := reflect.New(innerVal.Type())
		ttype := n.Interface().(Typed)
		err := parse(v, ttype)
		if err != nil {
			return err
		}
		data = append(data, n)
	}
	return typed.saveItems(&data)
}

func parseRecord(raw *fastjson.Value, typed iRecord) error {
	val, err := raw.Object()
	if err != nil {
		return err
	}
	typed.setJsonVal(raw)
	item := typed.anyItem()
	// logger.Log.Printf("Record %+v", item)
	innerVal := reflect.ValueOf(item).Elem()
	// logger.Log.Printf("innerVal %+v", innerVal)
	_, ok := innerVal.Addr().Interface().(Typed)
	if !ok {
		return fmt.Errorf("unable to convert to Typed interface %T", innerVal)
	}

	data := make(map[string]reflect.Value, val.Len())
	err = nil
	val.Visit(func(k []byte, v *fastjson.Value) {
		// Stop process if there is an error
		if err != nil {
			return
		}
		key := string(k)
		// logger.Log.Printf("key %s %+v", key, v)
		r := reflect.New(innerVal.Type())
		ttype := r.Interface().(Typed)
		// logger.Log.Println("type", ttype, ttype.Type())
		err = parse(v, ttype)

		data[key] = r
	})
	if err != nil {
		return err
	}
	return typed.saveItems(&data)
}

func lowerInitial(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}

// TODO: New doesn't store any value yet
func New[T any](val any) T {
	// var typed T
	typed := new(T)
	// TODO: Need to export fastjson.Value properties
	// so I can create them inline without parsing.
	// typed.setVal()
	// typed.parseValue()
	return *typed
}

// Serialized the Typed nodes to JSON
// TODO: I could try breaking out Typed types from Typed Schemas.
// Schemas would contain the logic for parsing and serialization.
// Typed types could be fairly dumb and allow for struct inlining.
// "NewJsonSchema" could wrap you top level Typed and handle json.
// Also could make it flexible to support something else like yaml.
//
// Typed types could have all method pointers removed.
// This could allow "Typed" type constraints without needing pointers everywhere.
// Would need to move parsing logic back into parser, not the types.
// Could be tricky for some more complex types with .saveItems()
//
// Looks like fastjson.Arena can be used to create new types
// https://github.com/valyala/fastjson/issues/69
// Get and Set json val shared types could be removed.
// This could be computed on the fly similar to parse()
func Serialize(typed Typed) []byte {
	val := typed.getJsonVal()
	return val.MarshalTo([]byte{})
}
