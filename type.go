package jsonvalue

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

// String provides a representation in string.
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
