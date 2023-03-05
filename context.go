package logy

import (
	"context"
	"sync/atomic"
)

const MappedContextKey = "logyMappedContext"

type Field struct {
	Key       string
	Value     any
	jsonValue atomic.Value
}

type MappedContext struct {
	values       []Field
	keyIndexMap  map[string]int
	encoder      *jsonEncoder
	jsonValue    atomic.Value
	contextIndex int
}

func NewMappedContext() *MappedContext {
	mc := &MappedContext{values: []Field{}, encoder: &jsonEncoder{}, keyIndexMap: map[string]int{}}
	mc.encoder.buf = newBuffer()
	mc.jsonValue.Store("{}")
	return mc
}

func (mc *MappedContext) Size() int {
	return len(mc.values)
}

func (mc *MappedContext) put(key string, value any) {
	mc.encoder.buf.Reset()
	if index, ok := mc.keyIndexMap[key]; ok {
		mc.values[index].Value = value
		mc.encoder.AddAny(key, value)
		mc.values[index].jsonValue.Store(mc.encoder.buf.String())
	} else {
		field := Field{key, value, atomic.Value{}}
		mc.keyIndexMap[key] = len(mc.values)
		mc.encoder.AddAny(key, value)
		field.jsonValue.Store(mc.encoder.buf.String())
		mc.values = append(mc.values, field)
	}

	mc.encoder.buf.Reset()
	mc.rewriteJson()
}

func (mc *MappedContext) Value(key string) any {
	if index, ok := mc.keyIndexMap[key]; ok {
		return mc.values[index].Value
	}

	return nil
}

func (mc *MappedContext) Values(callback func(key string, val any)) {
	for _, field := range mc.values {
		callback(field.Key, field.Value)
	}
}

func (mc *MappedContext) ValuesAsJson() string {
	return mc.jsonValue.Load().(string)
}

func (mc *MappedContext) clone() *MappedContext {
	c := *mc

	copyOfFields := make([]Field, len(mc.values))
	for index, field := range mc.values {
		copyOfFields[index] = field
	}
	c.values = copyOfFields

	copyOfKeyIndexMap := make(map[string]int, len(mc.keyIndexMap))
	for key, val := range mc.keyIndexMap {
		copyOfKeyIndexMap[key] = val
	}
	c.keyIndexMap = copyOfKeyIndexMap

	c.encoder = &jsonEncoder{buf: newBuffer()}
	return &c
}

func (mc *MappedContext) rewriteJson() {
	mc.encoder.buf.WriteByte('{')
	for _, field := range mc.values {
		mc.encoder.buf.WriteString(field.jsonValue.Load().(string))
	}
	mc.encoder.buf.WriteByte('}')
	mc.jsonValue.Store(mc.encoder.buf.String())
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

func WithValue(parent context.Context, key string, value any) context.Context {
	ctx := WithMappedContext(parent)
	val := ctx.Value(MappedContextKey)

	if val != nil {
		if mc, isMc := val.(*MappedContext); isMc {
			mc.put(key, value)
		}
	}

	return ctx
}

func ValueFrom(ctx context.Context, key string) any {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return nil
	}

	return mc.Value(key)
}
