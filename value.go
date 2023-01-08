package logy

import (
	"fmt"
	"time"
)

type Kind int

const (
	AnyKind Kind = iota
	BoolKind
	DurationKind
	Float64Kind
	Int64Kind
	StringKind
	TimeKind
	Uint64Kind
	GroupKind
)

var kindStrings = []string{
	"Any",
	"Bool",
	"Duration",
	"Float64",
	"Int64",
	"String",
	"Time",
	"Uint64",
	"GroupKind",
}

func (k Kind) String() string {
	if k >= 0 && int(k) < len(kindStrings) {
		return kindStrings[k]
	}
	return "<unknown logy.Kind>"
}

type Value struct {
	data any
	kind Kind
}

func newValue(data any, kind Kind) Value {
	return Value{
		data: data,
		kind: kind,
	}
}

func (v Value) Kind() Kind {
	return v.kind
}

func (v Value) String() string {
	if v.kind != StringKind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, StringKind))
	}

	return v.data.(string)
}

func (v Value) Int64() int64 {
	if v.kind != Int64Kind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, Int64Kind))
	}

	return v.data.(int64)
}

func (v Value) Uint64() uint64 {
	if v.kind != Uint64Kind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, Uint64Kind))
	}

	return v.data.(uint64)
}

func (v Value) Bool() bool {
	if v.kind != BoolKind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, BoolKind))
	}

	return v.data.(bool)
}

func (v Value) Duration() time.Duration {
	if v.kind != DurationKind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, DurationKind))
	}

	return v.data.(time.Duration)
}

func (v Value) Time() time.Time {
	if v.kind != TimeKind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, TimeKind))
	}

	return v.data.(time.Time)
}

func (v Value) Float64() float64 {
	if v.kind != Float64Kind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, Float64Kind))
	}

	return v.data.(float64)
}

func (v Value) Group() []Attribute {
	if v.kind != GroupKind {
		panic(fmt.Sprintf("Value kind is %s, not %s", v.kind, GroupKind))
	}

	return v.data.([]Attribute)
}
