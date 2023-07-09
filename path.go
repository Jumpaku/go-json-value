package jsonvalue

import (
	"strconv"

	"github.com/Jumpaku/go-assert"
	"golang.org/x/exp/slices"
)

type Key string

func (e Key) String() string {
	return string(e)
}
func (e Key) Integer() int {
	elm, err := strconv.ParseInt(e.String(), 10, 64)
	assert.State(err == nil, "%v cannot be parsed to int: %w", e, err)

	return int(elm)
}

func KeyInt(i int) Key {
	return Key(strconv.FormatInt(int64(i), 10))
}

type Path []Key

func (k Path) Equals(other Path) bool {
	return slices.Equal(k, other)
}
func (k Path) Len() int {
	return len(k)
}
func (k Path) Append(key Key) Path {
	return Path(append(append([]Key{}, k...), key))
}

func (k Path) Slice(begin int, endExclusive int) Path {
	assert.Params(0 <= begin && begin <= k.Len(), "begin %v must be in [0, %d]", begin, k.Len())
	assert.Params(0 <= endExclusive && endExclusive <= k.Len(), "endExclusive %v must be in [0, %d]", endExclusive, k.Len())
	assert.Params(begin <= endExclusive, "begin %v and endExclusive %v must be begin <= endExclusive", begin, endExclusive)

	return Path(append([]Key{}, k[begin:endExclusive]...))
}

func (k Path) Get(index int) Key {
	assert.Params(0 <= index && index < len(k), "index must be in [0, %d)", len(k))

	return k[index]
}

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
