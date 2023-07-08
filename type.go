package jsonvalue

// Type is
type Type int

const (
	TypeNull Type = iota
	TypeString
	TypeNumber
	TypeBoolean
	TypeArray
	TypeObject
)

func (t Type) String() string {
	switch t {
	case TypeNull:
		return `null`
	case TypeString:
		return `string`
	case TypeNumber:
		return `number`
	case TypeBoolean:
		return `boolean`
	case TypeArray:
		return `array`
	case TypeObject:
		return `object`
	default:
		panic("invalid JsonType")
	}
}
