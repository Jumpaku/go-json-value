# go-json-value
A library providing Go interface to dynamically and flexibly manipulate JSON-structured values

## Overview

The go-json-value package is a Go library that provides an interface for dynamically and flexibly manipulating JSON-structured data. It offers functionality for creating, modifying, and accessing JSON-structured data with interface in Go. Additionally, it includes features for working with keys, paths, and traversing JSON structures.

## Usage

### Installation

Use `go get` to install the package:

```sh
go get -u github.com/Jumpaku/go-json-value
```

### Import

Import the jsonvalue package into your Go program:

```go
import "github.com/Jumpaku/go-json-value"
```

### Key API

Types for JSON values:
```go
// Type represents type of JSON values.
type Type int

const (
	// TypeNull represents the type for the JSON null value.
	TypeNull Type = iota
	// TypeString represents the type for the JSON string values.
	TypeString
	// TypeNumber represents the type for the JSON number values.
	TypeNumber
	// TypeBoolean represents the type for the JSON boolean values.
	TypeBoolean
	// TypeArray represents the type for the JSON array values.
	TypeArray
	// TypeObject represents the type for the JSON object values.
	TypeObject
)
```

Interface modeling JSON-structured data.
```go
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
```

Functions for creation of JSON values:
```go
// Null returns a JSON value representing null.
func Null() Value

// Boolean returns a JSON boolean value.
func Boolean(b bool) Value

// String returns a JSON string value of s.
func String(s string) Value

// Number returns a JSON number value of n.
func Number[V ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | json.Number](n V) Value

// Object returns a JSON object value containing specified properties.
func Object(p ...Props) Value
// Props representing properties of JSON object.
type Props map[string]Value

// Array returns a JSON array value containing specified values.
func Array(vs ...Value) Value
```

Functions for visiting each value included in a JSON value:
```go
// Walk traverses a JSON value v and calls the visitor function for each the JSON values included in v.
// If a call of visitor returned an error, Walk immediately returns with the error.
func Walk(v Value, visitor func(path Path, val Value) error) error

// Find finds the JSON value specified by the Path in a JSON value v.
// If the JSON value associated with the Path exists, the found JSON value and true are returned; otherwise nil and false are returned.
func Find(v Value, path Path) (Value, bool)
```
