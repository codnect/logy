package logy

import (
	"time"
)

// An Attribute is a key-value pair.
type Attribute struct {
	Key   string
	Value Value
}

func newAttribute(key string, value Value) Attribute {
	return Attribute{
		key,
		value,
	}
}

func Attr[T any](key string, value T) Attribute {
	switch typed := any(value).(type) {
	case string:
		return newAttribute(key, newValue(typed, StringKind))
	case int:
		return newAttribute(key, newValue(int64(typed), Int64Kind))
	case uint:
		return newAttribute(key, newValue(uint64(typed), Uint64Kind))
	case int64:
		return newAttribute(key, newValue(typed, Int64Kind))
	case uint64:
		return newAttribute(key, newValue(typed, Uint64Kind))
	case bool:
		return newAttribute(key, newValue(typed, BoolKind))
	case time.Time:
		return newAttribute(key, newValue(typed, TimeKind))
	case time.Duration:
		return newAttribute(key, newValue(typed, DurationKind))
	case int8:
		return newAttribute(key, newValue(int64(typed), Int64Kind))
	case int16:
		return newAttribute(key, newValue(int64(typed), Int64Kind))
	case int32:
		return newAttribute(key, newValue(int64(typed), Int64Kind))
	case uint8:
		return newAttribute(key, newValue(uint64(typed), Uint64Kind))
	case uint16:
		return newAttribute(key, newValue(uint64(typed), Uint64Kind))
	case uint32:
		return newAttribute(key, newValue(uint64(typed), Uint64Kind))
	case uintptr:
		return newAttribute(key, newValue(uint64(typed), Uint64Kind))
	case float32:
		return newAttribute(key, newValue(float64(typed), Float64Kind))
	case float64:
		return newAttribute(key, newValue(typed, Float64Kind))
	default:
		return newAttribute(key, newValue(typed, AnyKind))
	}
}

func GroupAttrs(key string, attrs ...Attribute) Attribute {
	return newAttribute(key, newValue(attrs, GroupKind))
}
