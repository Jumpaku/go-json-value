package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Jumpaku/go-assert"
)

// Value models a JSON-structured data.
type Value interface {
	json.Marshaler
	json.Unmarshaler
	// Type returns JSON type.
	Type() Type
	// Assign assigns a JSON value to this object.
	Assign(v Value)
	// Clone deeply copies itself.
	Clone() Value
	// NumberGet returns this JSON value as a number.
	NumberGet() json.Number
	// StringGet returns this JSON value as a string.
	StringGet() string
	// BooleanGet returns this JSON value as a boolean.
	BooleanGet() bool
	// ObjectKeys returns keys of this JSON value as a object.
	ObjectKeys() []string
	// ObjectHasElm returns whether this JSON value as a object has the key.
	ObjectHasElm(key string) bool
	// ObjectHasElm returns a JSON value associated the key.
	ObjectGetElm(key string) Value
	// ObjectHasElm associates a JSON value by the key.
	ObjectSetElm(key string, v Value)
	// ObjectDelElm deletes the key and the associated JSON value.
	ObjectDelElm(key string)
	// ObjectLen returns the number of keys.
	ObjectLen() int
	// ArrayGetElm returns a JSON value indexed.
	ArrayGetElm(index int) Value
	// ArraySetElm sets a JSON value at the index.
	ArraySetElm(index int, v Value)
	// ArrayAddElm adds JSON values to the back.
	ArrayAddElm(vs ...Value)
	// ArrayLen returns the number of elements.
	ArrayLen() int
	// ArraySlice returns a sliced JSON array.
	ArraySlice(begin int, endExclusive int) Value
}

// Props representing properties of JSON object.
type Props map[string]Value

// value implements Value
type value struct {
	typ        Type
	numberVal  json.Number
	booleanVal bool
	stringVal  string
	objectVal  Props
	arrayVal   []Value
}

// Null returns a JSON value representing null.
func Null() Value {
	return &value{typ: TypeNull}
}

// Boolean returns a JSON boolean value.
func Boolean(b bool) Value {
	return &value{typ: TypeBoolean, booleanVal: b}
}

// String returns a JSON string value of s.
func String(s string) Value {
	return &value{typ: TypeString, stringVal: s}
}

// Number returns a JSON number value of n.
func Number[V ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | json.Number](n V) Value {
	var v json.Number
	var a any = n
	switch a := a.(type) {
	case int:
		v = json.Number(strconv.FormatInt(int64(a), 10))
	case int8:
		v = json.Number(strconv.FormatInt(int64(a), 10))
	case int16:
		v = json.Number(strconv.FormatInt(int64(a), 10))
	case int32:
		v = json.Number(strconv.FormatInt(int64(a), 10))
	case int64:
		v = json.Number(strconv.FormatInt(int64(a), 10))
	case uint:
		v = json.Number(strconv.FormatUint(uint64(a), 10))
	case uint8:
		v = json.Number(strconv.FormatUint(uint64(a), 10))
	case uint16:
		v = json.Number(strconv.FormatUint(uint64(a), 10))
	case uint32:
		v = json.Number(strconv.FormatUint(uint64(a), 10))
	case uint64:
		v = json.Number(strconv.FormatUint(uint64(a), 10))
	case float32:
		v = json.Number(strconv.FormatFloat(float64(a), 'f', 16, 64))
	case float64:
		v = json.Number(strconv.FormatFloat(float64(a), 'f', 16, 64))
	case json.Number:
		v = a
	}

	return &value{typ: TypeNumber, numberVal: json.Number(v)}
}

// Object returns a JSON object value containing specified properties.
func Object(p ...Props) Value {
	o := Props{}
	for _, m := range p {
		for k, v := range m {
			assert.Params(v != nil, "Value must not be nil")
			o[k] = v
		}
	}

	return &value{typ: TypeObject, objectVal: o}
}

// Array returns a JSON array value containing specified values.
func Array(vs ...Value) Value {
	a := make([]Value, len(vs))
	for i, v := range vs {
		assert.Params(v != nil, "Value must not be nil")
		a[i] = v
	}

	return &value{typ: TypeArray, arrayVal: a}
}

func (v *value) MarshalJSON() ([]byte, error) {
	switch v.Type() {
	case TypeNull:
		return json.Marshal(nil)
	case TypeBoolean:
		return json.Marshal(v.booleanVal)
	case TypeNumber:
		return json.Marshal(v.numberVal)
	case TypeString:
		return json.Marshal(v.stringVal)
	case TypeArray:
		return json.Marshal(v.arrayVal)
	case TypeObject:
		return json.Marshal(v.objectVal)
	default:
		return assert.Unexpected2[[]byte, error](`invalid JsonType: %v`, v.Type())
	}
}

func fromGo(a any) Value {
	switch a := a.(type) {
	case nil:
		return Null()
	case json.Number:
		return Number(a)
	case string:
		return String(a)
	case bool:
		return Boolean(a)
	case []any:
		arr := Array()
		for _, a := range a {
			arr.ArrayAddElm(fromGo(a))
		}
		return arr
	case map[string]any:
		obj := Object()
		for k, a := range a {
			obj.ObjectSetElm(k, fromGo(a))
		}
		return obj
	default:
		return assert.Unexpected1[Value]("unexpected value that cannot be converted to Value: %#v", a)
	}
}
func (v *value) UnmarshalJSON(b []byte) error {
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.UseNumber()

	var a any
	if err := decoder.Decode(&a); err != nil {
		return fmt.Errorf(`fail to unmarshal value to Value: %w`, err)
	}

	v.Assign(fromGo(a))

	return nil
}

func (v *value) Type() Type {
	return v.typ
}
func (v *value) Assign(other Value) {
	v.typ = other.Type()
	switch other.Type() {
	default:
		assert.Unexpected("unexpected Type: %v", other.Type())
	case TypeArray:
		l := other.ArrayLen()
		v.arrayVal = make([]Value, l)
		for i := 0; i < l; i++ {
			v.arrayVal[i] = other.ArrayGetElm(i)
		}
	case TypeObject:
		v.objectVal = map[string]Value{}
		keys := other.ObjectKeys()
		for _, k := range keys {
			v.objectVal[k] = other.ObjectGetElm(k)
		}
	case TypeBoolean:
		v.booleanVal = other.BooleanGet()
	case TypeNumber:
		v.numberVal = other.NumberGet()
	case TypeString:
		v.stringVal = other.StringGet()
	}
}

func (v *value) Clone() Value {
	switch v.Type() {
	case TypeArray:
		clone := Array()
		for i := 0; i < v.ArrayLen(); i++ {
			e := v.ArrayGetElm(i)
			clone.ArrayAddElm(e.Clone())
		}
		return clone
	case TypeObject:
		clone := Object()
		for _, k := range v.ObjectKeys() {
			e := v.ObjectGetElm(k)
			clone.ObjectSetElm(k, e.Clone())
		}
		return clone
	case TypeBoolean:
		return Boolean(v.BooleanGet())
	case TypeNumber:
		return Number(v.NumberGet())
	case TypeString:
		return String(v.StringGet())
	case TypeNull:
		return Null()
	default:
		return assert.Unexpected1[Value](`invalid JsonType: %v`, v.Type())
	}
}

func (v *value) NumberGet() json.Number {
	assert.Params(v.Type() == TypeNumber, "Value must be JSON number")

	return v.numberVal
}
func (v *value) StringGet() string {
	assert.Params(v.Type() == TypeString, "Value must be JSON string")

	return v.stringVal
}
func (v *value) BooleanGet() bool {
	assert.Params(v.Type() == TypeBoolean, "Value must be JSON boolean")

	return v.booleanVal
}
func (v *value) ObjectKeys() []string {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")

	keys := []string{}
	for key := range v.objectVal {
		keys = append(keys, key)
	}

	return keys
}
func (v *value) ObjectHasElm(key string) bool {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")

	_, ok := v.objectVal[key]

	return ok
}
func (v *value) ObjectGetElm(key string) Value {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")
	assert.Params(v.ObjectHasElm(key), "Value object must have key: %v", key)

	return v.objectVal[key]
}
func (v *value) ObjectSetElm(key string, val Value) {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")
	assert.Params(val != nil, "Value must be not nil")

	v.objectVal[key] = val
}
func (v *value) ObjectDelElm(key string) {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")

	delete(v.objectVal, key)
}
func (v *value) ObjectLen() int {
	assert.Params(v.Type() == TypeObject, "Value must be JSON object")

	return len(v.objectVal)
}
func (v *value) ArrayGetElm(index int) Value {
	assert.Params(v.Type() == TypeArray, "Value must be JSON array")
	assert.Params(0 <= index && index < v.ArrayLen(), "index must be in [0, %d)", v.ArrayLen())

	return v.arrayVal[index]
}
func (v *value) ArraySetElm(index int, val Value) {
	assert.Params(v.Type() == TypeArray, "Value must be JSON array")
	assert.Params(0 <= index && index < v.ArrayLen(), "index must be in [0, %d)", v.ArrayLen())

	v.arrayVal[index] = val
}
func (v *value) ArrayLen() int {
	assert.Params(v.Type() == TypeArray, "Value must be JSON array")

	return len(v.arrayVal)
}
func (v *value) ArrayAddElm(vals ...Value) {
	assert.Params(v.Type() == TypeArray, "Value must be JSON array")

	v.arrayVal = append(v.arrayVal, vals...)
}
func (v *value) ArraySlice(begin int, endExclusive int) Value {
	assert.Params(v.Type() == TypeArray, "Value must be JSON array")
	assert.Params(0 <= begin && begin <= v.ArrayLen(), "begin %v must be in [0, %d]", begin, v.ArrayLen())
	assert.Params(0 <= endExclusive && endExclusive <= v.ArrayLen(), "endExclusive %v must be in [0, %d]", endExclusive, v.ArrayLen())
	assert.Params(begin <= endExclusive, "begin %v and endExclusive %v must be begin <= endExclusive", begin, endExclusive)

	return Array(v.arrayVal[begin:endExclusive]...)
}
