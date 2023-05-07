package typed

import (
	"log"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
)

type TypedImage struct{ String }

// type TypedName

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
	return s.Val
}

type TestObjEnvProps struct {
	Key   String
	Value String
}
type TestObjEnv = Object[TestObjEnvProps]

type ServiceProps struct {
	Name String
	Port Number[int]
}
type Service = Object[ServiceProps]

type TestObjProps struct {
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
	Envs     Array[TestObjEnv]
	Tags     Array[String]
	Services Record[Service]
	MoreTags Record[String]
}
type TestObj = Object[TestObjProps]

func TestParsing(t *testing.T) {
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
	expect(t, schema.Props.IncludedStr.OrElse(NewString("fallback")).Value(), "included")
	// expect(t, schema.Props.IncludedStr.OrElse(String{Val: "fallback"}).Value(), "included")

	expect(t, schema.Props.ExcludedStr.IsPresent(), false)
	expect(t, schema.Props.ExcludedStr.IsAbsent(), true)
	expect(t, schema.Props.ExcludedStr.Maybe(), nil)
	expect(t, schema.Props.ExcludedStr.OrElse(NewString("fallback")).Value(), "fallback")
	// expect(t, schema.Props.ExcludedStr.OrElse(String{Val: "fallback"}).Value(), "fallback")

	// Test Serialization
	start = time.Now()
	json := Serialize(schema)
	duration = time.Since(start)
	log.Printf("Serialized in %fs or %d microseconds", duration.Seconds(), duration.Microseconds())
	snaps.MatchJSON(t, json)
}

// Test for writing typed declarations inline.
func TestInlined(t *testing.T) {
	start := time.Now()
	data := NewObject(TestObjProps{
		Id:            NewString("gabe"),
		Image:         TypedImage{NewString("nginx:10")},
		IntNumber:     NewNumber[int8](123),
		IsProduction:  NewBoolean(true),
		IsDevelopment: NewBoolean(false),
		IncludedStr:   NewOptional(NewString("included")),
		Env: NewObject(struct {
			Key   String
			Value String
		}{
			Key:   NewString("thing"),
			Value: NewString("stuff"),
		}),
		Envs: NewArray([]TestObjEnv{
			NewObject(TestObjEnvProps{
				Key:   NewString("another"),
				Value: NewString("one"),
			}),
		}),
		Tags: NewArray([]String{
			// Both init styles
			{Val: "production"},
			NewString("server"),
		}),
		Services: NewRecord(map[string]Service{
			"service1": NewObject(ServiceProps{
				Name: NewString("service1"),
				Port: NewNumber(8000),
			}),
			"service2": NewObject(ServiceProps{
				Name: NewString("service2"),
				Port: NewNumber(8080),
			}),
		}),
		MoreTags: NewRecord(map[string]String{
			"type":      {Val: "thing"},
			"something": NewString("else"),
		}),
	})
	duration := time.Since(start)
	log.Printf("Inline typed created in %fs or %d microseconds", duration.Seconds(), duration.Microseconds())

	expect(t, data.Props.Id.Value(), "gabe")
	expect(t, data.Props.Image.Value(), "nginx:10")
	expect(t, data.Props.IntNumber.Value(), 123)
	expect(t, data.Props.IsProduction.Value(), true)
	expect(t, data.Props.IsDevelopment.Value(), false)
	expect(t, data.Props.IncludedStr.OrElse(NewString("fallback")).Value(), "included")
	expect(t, data.Props.ExcludedStr.OrElse(NewString("fallback")).Value(), "fallback")
	expect(t, data.Props.Env.Props.Key.Value(), "thing")
	expect(t, data.Props.Env.Props.Value.Value(), "stuff")
	expect(t, data.Props.Envs.Items[0].Props.Key.Value(), "another")
	expect(t, data.Props.Envs.Items[0].Props.Value.Value(), "one")
	expect(t, data.Props.Tags.Items[0].Value(), "production")
	expect(t, data.Props.Tags.Items[1].Value(), "server")
	svc := data.Props.Services.Items["service1"]
	expect(t, svc.Props.Name.Value(), "service1")
	expect(t, svc.Props.Port.Value(), 8000)
	svc = data.Props.Services.Items["service2"]
	expect(t, svc.Props.Name.Value(), "service2")
	expect(t, svc.Props.Port.Value(), 8080)
	s := data.Props.MoreTags.Items["type"]
	expect(t, s.Value(), "thing")
	s = data.Props.MoreTags.Items["something"]
	expect(t, s.Value(), "else")

	// TODO: Test Serialization
	// snaps.MatchJSON(t, data)
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
		t.Errorf("got %#v, wanted %#v", got, want)
	}
}
