package jsonvalue

import (
	"strconv"

	"github.com/Jumpaku/go-assert"
	"golang.org/x/exp/slices"
)

// Key represents a key for a member in a JSON object or an index for an element in a JSON array.
type Key string

// KeyInt creates Key from JOSN array index.
func KeyInt(index int) Key {
	return Key(strconv.FormatInt(int64(index), 10))
}

// String returns key for JSON object in string.
func (k Key) String() string {
	return string(k)
}

// String returns index for JSON array in int.
func (k Key) Int() int {
	elm, err := strconv.ParseInt(k.String(), 10, 64)
	assert.State(err == nil, "%v cannot be parsed to int: %w", k, err)

	return int(elm)
}

// Path represents a sequence of the Keys.
// An empty path Path{} represents a root JSON value.
type Path []Key

// Equals returns true if two Paths contain the same keys in the same order; otherwise false.
func (p Path) Equals(other Path) bool {
	return slices.Equal(p, other)
}

// Len returns the number of keys.
func (p Path) Len() int {
	return len(p)
}

// Append creates a new Path appended the key to the end of the original Path.
func (p Path) Append(key Key) Path {
	return Path(append(append([]Key{}, p...), key))
}

// Slice creates a new Path sliced from begin to and of the original Path.
func (p Path) Slice(begin int, endExclusive int) Path {
	assert.Params(0 <= begin && begin <= p.Len(), "begin %v must be in [0, %d]", begin, p.Len())
	assert.Params(0 <= endExclusive && endExclusive <= p.Len(), "endExclusive %v must be in [0, %d]", endExclusive, p.Len())
	assert.Params(begin <= endExclusive, "begin %v and endExclusive %v must be begin <= endExclusive", begin, endExclusive)

	return Path(append([]Key{}, p[begin:endExclusive]...))
}

// Get returns a Key at the index.
func (p Path) Get(index int) Key {
	assert.Params(0 <= index && index < len(p), "index must be in [0, %d)", len(p))

	return p[index]
}

// Walk traverses a JSON value v and calls the visitor function for each the JSON values included in v.
// If a call of visitor returned an error, Walk immediately returns with the error.
func Walk(v Value, visitor func(path Path, val Value) error) error {
	return walkImpl(Path{}, v, visitor)
}

func walkImpl(parentKey Path, val Value, walkFunc func(key Path, val Value) error) error {
	if err := walkFunc(parentKey, val); err != nil {
		return err
	}
	switch val.Type() {
	case TypeObject:
		for _, key := range val.ObjectKeys() {
			val := val.ObjectGetElm(key)
			if err := walkImpl(parentKey.Append(Key(key)), val, walkFunc); err != nil {
				return err
			}
		}
	case TypeArray:
		for i := 0; i < val.ArrayLen(); i++ {
			val := val.ArrayGetElm(i)
			if err := walkImpl(parentKey.Append(Key(strconv.FormatInt(int64(i), 10))), val, walkFunc); err != nil {
				return err
			}
		}
	}
	return nil
}

// Find finds the JSON value specified by the Path in a JSON value v.
// If the JSON value associated with the Path exists, the found JSON value and true are returned; otherwise nil and false are returned.
func Find(v Value, path Path) (Value, bool) {
	var found Value
	_ = Walk(v, func(p Path, val Value) error {
		if p.Equals(path) {
			found = val
		}
		return nil
	})

	return found, found != nil
}
