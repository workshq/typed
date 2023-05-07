package typed

import (
	"log"
	"testing"
	"time"
)

type TypedImage struct {
	String
}

// TODO: Support custom validation and value formatting
// Order should be
// 1. Validate and catch any errors
//   - Validate is called during parsing
//
// 2. Value() will return data type
//   - A custom type can be used here (e.g. EID).
func (*TypedImage) Validate() error {
	return nil
}

func (s *TypedImage) Value() string {
	return s.RawVal
}

// TODO: Should all common options from json schema
// I'll need this to build the json schema deceleration.
// Either can pass this data in tags or use
// some structured data similar to zod or typebox.
//
// TODO: Making a field a pointer should make it optional in json-schema
type TestObj = Object[struct {
	Id            String `name:"id"`
	Image         TypedImage
	IntNumber     Number[int8]
	IsProduction  Boolean `name:"is_production"`
	IsDevelopment Boolean `name:"is_development"`
	IncludedStr   Optional[String]
	ExcludedStr   Optional[String]
	Env           Object[struct {
		Key   String
		Value String
	}]
	Envs Array[Object[struct {
		Key   String
		Value String
	}]]
	Tags     Array[String]
	Services Record[Object[struct {
		Name String
		Port Number[int]
	}]]
	MoreTags Record[String]
}]

func TestTyped(t *testing.T) {
	// items := new(Array[Object[struct {
	// 	Key   String
	// 	Value String
	// }]])
	// it := &items.Items
	// r := reflect.ValueOf(it).Elem().Type()
	// // r.
	// // r.
	// log.Println(r)
	// log.Println(r.Kind())
	// log.Println(r.Elem())
	// return;

	start := time.Now()
	schema := new(TestObj)
	data := `
  {
    "id": "gabemx",
    "image": "nginx:10",
    "intNumber": 100,
		"is_production": true,
		"is_development": false,
    "includedStr": "included",
    "env": {
      "key": "thing",
      "value": "stuff"
    },
    "envs": [
      {
        "key": "another",
        "value": "one"
      }
    ],
    "tags": ["production", "server"],
    "services": {
      "service1": {
        "name": "service1",
        "port": 8000
      },
      "service2": {
        "name": "service2",
        "port": 8080
      }
    },
    "moreTags": {
      "type": "thing",
      "something": "else"
    }
  }`
	err := Parse(data, schema)
	duration := time.Since(start)
	log.Printf("Parsed in %fs or %d microseconds", duration.Seconds(), duration.Microseconds())
	if err != nil {
		t.Error(err)
		return
	}
	println()
	log.Printf("%+v", schema.Props)

	// Basic
	expect(t, schema.Props.Id.Value(), "gabemx")
	expect(t, schema.Props.Image.Value(), "nginx:10")
	expect(t, schema.Props.IntNumber.Value(), 100)
	expect(t, schema.Props.Env.Props.Key.Value(), "thing")
	expect(t, schema.Props.Env.Props.Value.Value(), "stuff")
	expect(t, schema.Props.IsProduction.Value(), true)
	expect(t, schema.Props.IsDevelopment.Value(), false)
	// Array
	expect(t, schema.Props.Envs.Items[0].Props.Key.Value(), "another")
	expect(t, schema.Props.Envs.Items[0].Props.Value.Value(), "one")
	expect(t, schema.Props.Tags.Items[0].Value(), "production")
	expect(t, schema.Props.Tags.Items[1].Value(), "server")
	// Record Object
	svc := schema.Props.Services.Items["service1"]
	expect(t, svc.Props.Name.Value(), "service1")
	expect(t, svc.Props.Port.Value(), 8000)
	svc = schema.Props.Services.Items["service2"]
	expect(t, svc.Props.Name.Value(), "service2")
	expect(t, svc.Props.Port.Value(), 8080)
	// Record String
	s := schema.Props.MoreTags.Items["type"]
	expect(t, s.Value(), "thing")
	s = schema.Props.MoreTags.Items["something"]
	expect(t, s.Value(), "else")
	// Optional
	expect(t, schema.Props.IncludedStr.IsPresent(), true)
	expect(t, schema.Props.IncludedStr.IsAbsent(), false)
	expect(t, schema.Props.IncludedStr.Maybe().Value(), "included")
	expect(t, schema.Props.IncludedStr.OrElse(String{RawVal: "fallback"}).Value(), "included")

	expect(t, schema.Props.ExcludedStr.IsPresent(), false)
	expect(t, schema.Props.ExcludedStr.IsAbsent(), true)
	expect(t, schema.Props.ExcludedStr.Maybe(), nil)
	expect(t, schema.Props.ExcludedStr.OrElse(String{RawVal: "fallback"}).Value(), "fallback")

	// t.props.env.props.name
  start = time.Now()
  json := ToJSON(schema)
  duration = time.Since(start)
	log.Printf("Serialized in %fs or %d microseconds", duration.Seconds(), duration.Microseconds())
  log.Printf("json: %s", json)
}

// Testing type definitions instead of type aliases
func XTestTypeDef(t *testing.T) {
	// TODO: This doesn't work
	// Looks like the interface was lying for this syntax
	type Test2 String
	// But this does
	type Test3 struct{ String }
	c := Test3{}
	log.Printf("type %s", c.Type())
	// x := new(Test2)
	// log.Printf("type %s", x.Type())
}

func expect[T comparable](t *testing.T, got T, want T) {
	if got != want {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}
