package jsonvalue_test

import (
	"testing"

	jsonvalue "github.com/Jumpaku/go-json-value"
)

func TestType_String(t *testing.T) {
	t.Run(`null`, func(t *testing.T) {
		equal(t, jsonvalue.TypeNull.String(), "null")
	})
	t.Run(`string`, func(t *testing.T) {
		equal(t, jsonvalue.TypeString.String(), "string")
	})
	t.Run(`number`, func(t *testing.T) {
		equal(t, jsonvalue.TypeNumber.String(), "number")
	})
	t.Run(`boolean`, func(t *testing.T) {
		equal(t, jsonvalue.TypeBoolean.String(), "boolean")
	})
	t.Run(`object`, func(t *testing.T) {
		equal(t, jsonvalue.TypeObject.String(), "object")
	})
	t.Run(`array`, func(t *testing.T) {
		equal(t, jsonvalue.TypeArray.String(), "array")
	})
}
