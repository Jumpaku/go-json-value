package jsonvalue_test

import (
	"fmt"
	"reflect"
	"testing"

	jsonvalue "github.com/Jumpaku/go-json-value"
)

func isNil(value any) bool {
	rv := reflect.ValueOf(value)

	if !rv.IsValid() {
		return true
	}

	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	}

	return false
}

func IsNotNil(t *testing.T, actual any) {
	t.Helper()

	if isNil(actual) {
		t.Errorf("ASSERT IS NOT NIL\n  actual: %v:%T", actual, actual)
	}
}

func TestKey_String(t *testing.T) {
	v := jsonvalue.Key("abc").String()
	equal(t, v, "abc")
}

func TestKey_Integer(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := jsonvalue.Key("123").Integer()
		equal(t, v, 123)
	})
}

func TestPath_Equals(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		p := jsonvalue.Path([]jsonvalue.Key{"abc", "123"})
		equal(t, p.Equals(jsonvalue.Path([]jsonvalue.Key{"abc", "123"})), true)
	})
	t.Run("not equal", func(t *testing.T) {
		p := jsonvalue.Path([]jsonvalue.Key{"abc", "123"})
		equal(t, p.Equals(jsonvalue.Path([]jsonvalue.Key{"abc", "123", "xyz"})), false)
	})
}

func TestPath_Get(t *testing.T) {
	p := jsonvalue.Path([]jsonvalue.Key{"abc", "123"})
	k0 := p.Get(0)
	equal(t, k0, jsonvalue.Key("abc"))
	k1 := p.Get(1)
	equal(t, k1, jsonvalue.Key("123"))
}

func TestPath_Len(t *testing.T) {
	p := jsonvalue.Path([]jsonvalue.Key{"abc", "123"})
	equal(t, p.Len(), 2)
}

func TestPath_Append(t *testing.T) {
	p := jsonvalue.Path([]jsonvalue.Key{"abc", "123"}).Append("xyz")
	equal(t, p.Len(), 3)
	k1 := p.Get(2)
	equal(t, k1, jsonvalue.Key("xyz"))
}

func TestWalk(t *testing.T) {
	t.Run(`error`, func(t *testing.T) {
		v := jsonvalue.Null()
		err := jsonvalue.Walk(v, func(path jsonvalue.Path, val jsonvalue.Value) error {
			return fmt.Errorf("")
		})
		IsNotNil(t, err)
	})
	t.Run(`null`, func(t *testing.T) {
		v := jsonvalue.Null()
		p := []jsonvalue.Path{}
		_ = jsonvalue.Walk(v, func(path jsonvalue.Path, val jsonvalue.Value) error {
			p = append(p, path)
			return nil
		})
		equal(t, len(p), 1)
	})
	t.Run(`object`, func(t *testing.T) {
		v := jsonvalue.Object(map[string]jsonvalue.Value{
			"a": jsonvalue.Null(),
			"b": jsonvalue.Object(map[string]jsonvalue.Value{
				"x": jsonvalue.Null(),
				"y": jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				"z": jsonvalue.Array(
					jsonvalue.Null(),
				),
			}),
			"c": jsonvalue.Array(
				jsonvalue.Null(),
				jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				jsonvalue.Array(
					jsonvalue.Null(),
				),
			),
		})
		p := []jsonvalue.Path{}
		_ = jsonvalue.Walk(v, func(path jsonvalue.Path, val jsonvalue.Value) error {
			p = append(p, path)
			return nil
		})
		equal(t, len(p), 14)
	})
	t.Run(`array`, func(t *testing.T) {
		v := jsonvalue.Array(
			jsonvalue.Null(),
			jsonvalue.Object(map[string]jsonvalue.Value{
				"x": jsonvalue.Null(),
				"y": jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				"z": jsonvalue.Array(
					jsonvalue.Null(),
				),
			}),
			jsonvalue.Array(
				jsonvalue.Null(),
				jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				jsonvalue.Array(
					jsonvalue.Null(),
				),
			),
		)
		p := []jsonvalue.Path{}
		_ = jsonvalue.Walk(v, func(path jsonvalue.Path, val jsonvalue.Value) error {
			p = append(p, path)
			return nil
		})
		equal(t, len(p), 14)
	})
}
func TestFind(t *testing.T) {
	t.Run(`not found`, func(t *testing.T) {
		t.Run(`null`, func(t *testing.T) {
			v := jsonvalue.Null()
			_, ok := jsonvalue.Find(v, jsonvalue.Path{"xxx"})
			equal(t, ok, false)
		})
		t.Run(`object`, func(t *testing.T) {
			v := jsonvalue.Object(map[string]jsonvalue.Value{
				"a": jsonvalue.Null(),
				"b": jsonvalue.Object(map[string]jsonvalue.Value{
					"x": jsonvalue.Null(),
					"y": jsonvalue.Object(map[string]jsonvalue.Value{
						"w": jsonvalue.Null(),
					}),
					"z": jsonvalue.Array(
						jsonvalue.Null(),
					),
				}),
				"c": jsonvalue.Array(
					jsonvalue.Null(),
					jsonvalue.Object(map[string]jsonvalue.Value{
						"w": jsonvalue.Null(),
					}),
					jsonvalue.Array(
						jsonvalue.Null(),
					),
				),
			})
			_, ok := jsonvalue.Find(v, jsonvalue.Path{"xxx"})
			equal(t, ok, false)
		})
		t.Run(`array`, func(t *testing.T) {
			v := jsonvalue.Array(
				jsonvalue.Null(),
				jsonvalue.Object(map[string]jsonvalue.Value{
					"x": jsonvalue.Null(),
					"y": jsonvalue.Object(map[string]jsonvalue.Value{
						"w": jsonvalue.Null(),
					}),
					"z": jsonvalue.Array(
						jsonvalue.Null(),
					),
				}),
				jsonvalue.Array(
					jsonvalue.Null(),
					jsonvalue.Object(map[string]jsonvalue.Value{
						"w": jsonvalue.Null(),
					}),
					jsonvalue.Array(
						jsonvalue.Null(),
					),
				),
			)
			_, ok := jsonvalue.Find(v, jsonvalue.Path{"xxx"})
			equal(t, ok, false)
		})
	})
	t.Run(`null`, func(t *testing.T) {
		v := jsonvalue.Null()
		a, ok := jsonvalue.Find(v, jsonvalue.Path{})
		equal(t, ok, true)
		equal(t, a.Type(), jsonvalue.TypeNull)
	})
	t.Run(`object`, func(t *testing.T) {
		v := jsonvalue.Object(map[string]jsonvalue.Value{
			"a": jsonvalue.Null(),
			"b": jsonvalue.Object(map[string]jsonvalue.Value{
				"x": jsonvalue.Null(),
				"y": jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				"z": jsonvalue.Array(
					jsonvalue.Null(),
				),
			}),
			"c": jsonvalue.Array(
				jsonvalue.Null(),
				jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				jsonvalue.Array(
					jsonvalue.Null(),
				),
			),
		})
		t.Run(".", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".a", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"a"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".b", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".b.x", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b", "x"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".b.y", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b", "y"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".b.y.w", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b", "y", "w"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".b.z", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b", "z"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".b.z.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"b", "z", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".c", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".c.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".c.1", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c", "1"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".c.1.w", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c", "1", "w"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".c.2", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c", "2"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".c.2.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"c", "2", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
	})
	t.Run(`array`, func(t *testing.T) {
		v := jsonvalue.Array(
			jsonvalue.Null(),
			jsonvalue.Object(map[string]jsonvalue.Value{
				"x": jsonvalue.Null(),
				"y": jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				"z": jsonvalue.Array(
					jsonvalue.Null(),
				),
			}),
			jsonvalue.Array(
				jsonvalue.Null(),
				jsonvalue.Object(map[string]jsonvalue.Value{
					"w": jsonvalue.Null(),
				}),
				jsonvalue.Array(
					jsonvalue.Null(),
				),
			),
		)
		t.Run(".", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".1", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".1.x", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1", "x"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".1.y", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1", "y"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".1.y.w", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1", "y", "w"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".1.z", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1", "z"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".1.z.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"1", "z", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".2", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".2.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".2.1", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2", "1"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeObject)
		})
		t.Run(".2.1.w", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2", "1", "w"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
		t.Run(".2.2", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2", "2"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeArray)
		})
		t.Run(".2.2.0", func(t *testing.T) {
			a, ok := jsonvalue.Find(v, jsonvalue.Path{"2", "2", "0"})
			equal(t, ok, true)
			equal(t, a.Type(), jsonvalue.TypeNull)
		})
	})
}
