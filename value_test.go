package jsonvalue_test

import (
	"encoding/json"
	"math"
	"testing"

	jsonvalue "github.com/Jumpaku/go-json-value"
	"golang.org/x/exp/slices"
)

func anyEqual(actual any, expect any) bool {
	return actual == expect
}
func equal[T comparable](t *testing.T, actual T, expect T) {
	t.Helper()

	if !anyEqual(actual, expect) {
		t.Errorf("ASSERT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}

func closeTo(t *testing.T, got any, want any, tolerance float64) bool {
	t.Helper()

	var g float64
	var w float64
	switch got := got.(type) {
	case float32:
		g = float64(got)
	case float64:
		g = float64(got)
	default:
		t.Fatalf(`got value is not float`)
	}
	switch want := want.(type) {
	case float32:
		w = float64(want)
	case float64:
		w = float64(want)
	default:
		t.Fatalf(`want value is not float`)
	}

	return math.Abs(g-w) <= tolerance
}

func TestType(t *testing.T) {
	type testCase struct {
		sut  jsonvalue.Value
		want jsonvalue.Type
	}
	testCases := []testCase{
		{
			sut:  jsonvalue.Null(),
			want: jsonvalue.TypeNull,
		},
		{
			sut:  jsonvalue.Number(123),
			want: jsonvalue.TypeNumber,
		},
		{
			sut:  jsonvalue.String("abc"),
			want: jsonvalue.TypeString,
		},
		{
			sut:  jsonvalue.Boolean(true),
			want: jsonvalue.TypeBoolean,
		},
		{
			sut:  jsonvalue.Array(),
			want: jsonvalue.TypeArray,
		},
		{
			sut:  jsonvalue.Object(),
			want: jsonvalue.TypeObject,
		},
	}

	for i, testCase := range testCases {
		got := testCase.sut.Type()
		if got != testCase.want {
			t.Errorf("case=%d: got != want\n  got  = %v\n  want = %v", i, got, testCase.want)
		}
	}
}

func TestNumberGet(t *testing.T) {
	type testCase struct {
		sut  jsonvalue.Value
		want any
	}
	testCases := []testCase{
		{
			sut:  jsonvalue.Number(json.Number("-123.45")),
			want: float64(-123.45),
		},
		{
			sut:  jsonvalue.Number(float32(-123.5)),
			want: float64(-123.5),
		},
		{
			sut:  jsonvalue.Number(float64(-123.45)),
			want: float64(-123.45),
		},
		{
			sut:  jsonvalue.Number(json.Number("123")),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(int(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(int8(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(int16(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(int32(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(int64(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(uint(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(uint8(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(uint16(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(uint32(123)),
			want: int64(123),
		},
		{
			sut:  jsonvalue.Number(uint64(123)),
			want: int64(123),
		},
	}
	for i, testCase := range testCases {
		got := testCase.sut.NumberGet()
		switch want := testCase.want.(type) {
		case float64:
			got, err := got.Float64()
			if err != nil {
				t.Errorf("case=%d: err = %#v", i, err)
			}
			if !closeTo(t, got, want, 1e-10) {
				t.Errorf("case=%d: got is not close to want\n  got  = %#v\n  want = %#v", i, got, want)
			}
		case int64:
			got, err := got.Int64()
			if err != nil {
				t.Errorf("case=%d: err = %#v", i, err)
			}
			if got != want {
				t.Errorf("case=%d: got != want\n  got  = %#v\n  want = %#v", i, got, want)
			}
		default:
			t.Fatal()
		}
	}
}

func TestStringGet(t *testing.T) {
	want := "abc"
	sut := jsonvalue.String("abc")
	got := sut.StringGet()
	equal(t, got, want)
}

func TestBooleanGet(t *testing.T) {
	t.Run(`true`, func(t *testing.T) {
		want := true
		sut := jsonvalue.Boolean(true)
		got := sut.BooleanGet()
		equal(t, got, want)
	})

	t.Run(`false`, func(t *testing.T) {
		want := false
		sut := jsonvalue.Boolean(false)
		got := sut.BooleanGet()
		equal(t, got, want)
	})
}

func checkSliceContains[T comparable](t *testing.T, a []T, val T) {
	t.Helper()
	if !slices.Contains(a, val) {
		t.Errorf(`"slice does not have %#v", val`, val)
	}
}
func TestObjectKeys(t *testing.T) {
	o := jsonvalue.Object(jsonvalue.Props{
		"a": jsonvalue.Null(),
		"b": jsonvalue.Number(123),
		"c": jsonvalue.String("abc"),
		"d": jsonvalue.Boolean(true),
		"e": jsonvalue.Object(),
		"f": jsonvalue.Array(),
	})
	aKeys := o.ObjectKeys()

	equal(t, len(aKeys), 6)
	checkSliceContains(t, aKeys, "a")
	checkSliceContains(t, aKeys, "b")
	checkSliceContains(t, aKeys, "c")
	checkSliceContains(t, aKeys, "d")
	checkSliceContains(t, aKeys, "e")
	checkSliceContains(t, aKeys, "f")
}

func TestObjectGetElm(t *testing.T) {
	o := jsonvalue.Object(jsonvalue.Props{
		"a": jsonvalue.Null(),
		"b": jsonvalue.Number(123),
		"c": jsonvalue.String("abc"),
		"d": jsonvalue.Boolean(true),
		"e": jsonvalue.Object(),
		"f": jsonvalue.Array(),
	})

	t.Run("get null", func(t *testing.T) {
		v := o.ObjectGetElm("a")
		equal(t, v.Type(), jsonvalue.TypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v := o.ObjectGetElm("b")
		equal(t, v.Type(), jsonvalue.TypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v := o.ObjectGetElm("c")
		equal(t, v.Type(), jsonvalue.TypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v := o.ObjectGetElm("d")
		equal(t, v.Type(), jsonvalue.TypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v := o.ObjectGetElm("e")
		equal(t, v.Type(), jsonvalue.TypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v := o.ObjectGetElm("f")
		equal(t, v.Type(), jsonvalue.TypeArray)
	})
}

func TestObjectSetElm(t *testing.T) {
	t.Run("set null", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("a", jsonvalue.Null())
		v := o.ObjectGetElm("a")
		equal(t, v.Type(), jsonvalue.TypeNull)
	})

	t.Run("set number", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("b", jsonvalue.Number(123))
		v := o.ObjectGetElm("b")
		equal(t, v.Type(), jsonvalue.TypeNumber)
	})

	t.Run("set string", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("c", jsonvalue.String("abc"))
		v := o.ObjectGetElm("c")
		equal(t, v.Type(), jsonvalue.TypeString)
	})

	t.Run("set boolean", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("d", jsonvalue.Boolean(true))
		v := o.ObjectGetElm("d")
		equal(t, v.Type(), jsonvalue.TypeBoolean)
	})

	t.Run("set object", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("e", jsonvalue.Object())
		v := o.ObjectGetElm("e")
		equal(t, v.Type(), jsonvalue.TypeObject)
	})

	t.Run("set array", func(t *testing.T) {
		o := jsonvalue.Object()
		o.ObjectSetElm("f", jsonvalue.Array())
		v := o.ObjectGetElm("f")
		equal(t, v.Type(), jsonvalue.TypeArray)
	})
}

func TestObjectDelElm(t *testing.T) {
	t.Run("delete null", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"a": jsonvalue.Null()})
		o.ObjectDelElm("a")
		ok := o.ObjectHasElm("a")
		equal(t, ok, false)
	})

	t.Run("delete number", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"b": jsonvalue.Number(123)})
		o.ObjectDelElm("b")
		ok := o.ObjectHasElm("b")
		equal(t, ok, false)
	})

	t.Run("delete string", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"c": jsonvalue.String("abc")})
		o.ObjectDelElm("c")
		ok := o.ObjectHasElm("c")
		equal(t, ok, false)
	})

	t.Run("delete boolean", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"d": jsonvalue.Boolean(true)})
		o.ObjectDelElm("d")
		ok := o.ObjectHasElm("d")
		equal(t, ok, false)
	})

	t.Run("delete object", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"e": jsonvalue.Object()})
		o.ObjectDelElm("e")
		ok := o.ObjectHasElm("e")
		equal(t, ok, false)
	})

	t.Run("delete array", func(t *testing.T) {
		o := jsonvalue.Object(jsonvalue.Props{"f": jsonvalue.Array()})
		o.ObjectDelElm("f")
		ok := o.ObjectHasElm("f")
		equal(t, ok, false)
	})
}

func TestObjectLen(t *testing.T) {
	o := jsonvalue.Object(jsonvalue.Props{
		"a": jsonvalue.Null(),
		"b": jsonvalue.Number(123),
		"c": jsonvalue.String("abc"),
		"d": jsonvalue.Boolean(true),
		"e": jsonvalue.Object(),
		"f": jsonvalue.Array(),
	})

	equal(t, o.ObjectLen(), 6)
}

func TestArrayGetElm(t *testing.T) {
	o := jsonvalue.Array(
		jsonvalue.Null(),
		jsonvalue.Number(123),
		jsonvalue.String("abc"),
		jsonvalue.Boolean(true),
		jsonvalue.Object(),
		jsonvalue.Array(),
	)

	t.Run("get null", func(t *testing.T) {
		v := o.ArrayGetElm(0)
		equal(t, v.Type(), jsonvalue.TypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v := o.ArrayGetElm(1)
		equal(t, v.Type(), jsonvalue.TypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v := o.ArrayGetElm(2)
		equal(t, v.Type(), jsonvalue.TypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v := o.ArrayGetElm(3)
		equal(t, v.Type(), jsonvalue.TypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v := o.ArrayGetElm(4)
		equal(t, v.Type(), jsonvalue.TypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v := o.ArrayGetElm(5)
		equal(t, v.Type(), jsonvalue.TypeArray)
	})
}

func newExampleArray(n int, v jsonvalue.Value) jsonvalue.Value {
	var a = make([]jsonvalue.Value, n)
	for i := 0; i < n; i++ {
		a[i] = v
	}
	return jsonvalue.Array(a...)
}
func TestArraySetElm(t *testing.T) {
	t.Run("get null", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Array())
		o.ArraySetElm(0, jsonvalue.Null())
		v := o.ArrayGetElm(0)
		equal(t, v.Type(), jsonvalue.TypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Null())
		o.ArraySetElm(1, jsonvalue.Number(123))
		v := o.ArrayGetElm(1)
		equal(t, v.Type(), jsonvalue.TypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Null())
		o.ArraySetElm(2, jsonvalue.String("abc"))
		v := o.ArrayGetElm(2)
		equal(t, v.Type(), jsonvalue.TypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Null())
		o.ArraySetElm(3, jsonvalue.Boolean(true))
		v := o.ArrayGetElm(3)
		equal(t, v.Type(), jsonvalue.TypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Null())
		o.ArraySetElm(4, jsonvalue.Object())
		v := o.ArrayGetElm(4)
		equal(t, v.Type(), jsonvalue.TypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		o := newExampleArray(6, jsonvalue.Null())
		o.ArraySetElm(5, jsonvalue.Array())
		v := o.ArrayGetElm(5)
		equal(t, v.Type(), jsonvalue.TypeArray)
	})
}

func TestArrayLen(t *testing.T) {
	o := jsonvalue.Array(
		jsonvalue.Null(),
		jsonvalue.Number(123),
		jsonvalue.String("abc"),
		jsonvalue.Boolean(true),
		jsonvalue.Object(),
		jsonvalue.Array(),
	)

	equal(t, o.ArrayLen(), 6)
}

func TestArrayAddElm(t *testing.T) {
	o := jsonvalue.Array()

	o.ArrayAddElm(jsonvalue.Null())
	o.ArrayAddElm(jsonvalue.Number(123))
	o.ArrayAddElm(jsonvalue.String("abc"))
	o.ArrayAddElm(jsonvalue.Boolean(true))
	o.ArrayAddElm(jsonvalue.Object())
	o.ArrayAddElm(jsonvalue.Array())

	t.Run("get null", func(t *testing.T) {
		v := o.ArrayGetElm(0)
		equal(t, v.Type(), jsonvalue.TypeNull)
	})

	t.Run("get number", func(t *testing.T) {
		v := o.ArrayGetElm(1)
		equal(t, v.Type(), jsonvalue.TypeNumber)
	})

	t.Run("get string", func(t *testing.T) {
		v := o.ArrayGetElm(2)
		equal(t, v.Type(), jsonvalue.TypeString)
	})

	t.Run("get boolean", func(t *testing.T) {
		v := o.ArrayGetElm(3)
		equal(t, v.Type(), jsonvalue.TypeBoolean)
	})

	t.Run("get object", func(t *testing.T) {
		v := o.ArrayGetElm(4)
		equal(t, v.Type(), jsonvalue.TypeObject)
	})

	t.Run("get array", func(t *testing.T) {
		v := o.ArrayGetElm(5)
		equal(t, v.Type(), jsonvalue.TypeArray)
	})
}

func TestArraySlice(t *testing.T) {
	o := jsonvalue.Array(
		jsonvalue.Null(),
		jsonvalue.Number(123),
		jsonvalue.String("abc"),
		jsonvalue.Boolean(true),
		jsonvalue.Object(),
		jsonvalue.Array(),
	)

	t.Run("whole", func(t *testing.T) {
		a := o.ArraySlice(0, o.ArrayLen())
		equal(t, a.ArrayLen(), 6)
	})

	t.Run("empty", func(t *testing.T) {
		a := o.ArraySlice(3, 3)
		equal(t, a.ArrayLen(), 0)
	})

	t.Run("sub-array", func(t *testing.T) {
		a := o.ArraySlice(2, 4)
		equal(t, a.ArrayLen(), 2)
	})
}
