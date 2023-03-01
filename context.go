package logy

import (
	"context"
	"sync/atomic"
)

const MappedContextKey = "logyMappedContext"

type Filter func(key string, val any) bool

type Field struct {
	Key   string
	Value any
}

type MappedContext struct {
	values      []Field
	keyIndexMap map[string]int
	size        int
	encoder     *jsonEncoder
	jsonValue   atomic.Value
}

func NewMappedContext() *MappedContext {
	mc := &MappedContext{values: []Field{}, keyIndexMap: map[string]int{}, encoder: &jsonEncoder{}}
	mc.encoder.buf = newBuffer()
	return mc
}

func (mc *MappedContext) Fields() []Field {
	return mc.values
}

func (mc *MappedContext) put(key string, value any) {
	if index, ok := mc.keyIndexMap[key]; ok {
		mc.values[index].Value = value
	} else {
		mc.values = append(mc.values, Field{key, value})
		mc.keyIndexMap[key] = mc.size
		mc.encoder.AddAny(key, value)
		mc.size++
	}

	mc.encoder.buf.Reset()
	mc.encoder.buf.WriteByte('{')
	for _, field := range mc.Fields() {
		mc.encoder.AddAny(field.Key, field.Value)
	}
	mc.encoder.buf.WriteByte('}')
	mc.jsonValue.Store(mc.encoder.buf.String())
	mc.encoder.buf.Reset()
}

func (mc *MappedContext) value(key string) any {
	if index, ok := mc.keyIndexMap[key]; ok {
		return mc.values[index].Value
	}

	return nil
}

func (mc *MappedContext) clone() *MappedContext {
	c := *mc
	c.encoder = &jsonEncoder{buf: newBuffer()}
	return &c
}

func (mc *MappedContext) ValuesAsText() string {
	return ""
}

func (mc *MappedContext) ValuesAsJSON(filter Filter) string {
	return mc.jsonValue.Load().(string)
}

func MappedContextFrom(ctx context.Context) *MappedContext {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(MappedContextKey)

	if val != nil {
		if mc, isMc := val.(*MappedContext); isMc {
			return mc
		}

		return nil
	}

	return nil
}

func WithMappedContext(ctx context.Context) context.Context {
	mc := MappedContextFrom(ctx)

	if mc != nil {
		return context.WithValue(ctx, MappedContextKey, mc.clone())
	}

	return context.WithValue(ctx, MappedContextKey, NewMappedContext())
}

func ToMappedContext(ctx context.Context, key string, value any) {
	mc := MappedContextFrom(ctx)

	if mc != nil {
		mc.put(key, value)
	}
}

func FromMappedContext(ctx context.Context, key string) any {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return nil
	}

	return mc.value(key)
}
