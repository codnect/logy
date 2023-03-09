package logy

import (
	"context"
	"sync/atomic"
)

const ContextKey = "$logyMappedContext"

type Field struct {
	key       string
	value     any
	jsonValue string
}

func (f Field) Key() string {
	return f.key
}

func (f Field) Value() any {
	return f.value
}

/*
func (f Field) AsJson() string {
	return f.jsonValue
}*/

type Iterator struct {
	fields  []Field
	current int
}

func (i *Iterator) HasNext() bool {
	return i.current < len(i.fields)
}

func (i *Iterator) Next() (Field, bool) {
	if i.current >= len(i.fields) {
		return Field{}, false
	}

	field := i.fields[i.current]
	i.current++
	return field, true
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
		mc.values[index].value = value
		mc.encoder.AddAny(key, value)
		mc.values[index].jsonValue = mc.encoder.buf.String()
	} else {
		field := Field{key, value, ""}

		mc.keyIndexMap[key] = len(mc.values)
		mc.encoder.AddAny(key, value)
		field.jsonValue = mc.encoder.buf.String()

		mc.values = append(mc.values, field)
	}

	mc.encoder.buf.Reset()
	mc.rewriteJson()
}

func (mc *MappedContext) Value(key string) any {
	if index, ok := mc.keyIndexMap[key]; ok {
		return mc.values[index].Value()
	}

	return nil
}

func (mc *MappedContext) Values(callback func(key string, val any)) {
	for _, field := range mc.values {
		callback(field.Key(), field.Value())
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
	for index, field := range mc.values {
		mc.encoder.buf.WriteString(field.jsonValue)
		if index != len(mc.values)-1 {
			mc.encoder.buf.WriteByte(',')
		}
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

func PutValue(ctx context.Context, key string, val any) {
	ctxVal := ctx.Value(MappedContextKey)

	if ctxVal != nil {
		if mc, isMc := ctxVal.(*MappedContext); isMc {
			mc.put(key, val)
		}
	}
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

func Values(ctx context.Context) *Iterator {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return &Iterator{}
	}

	return &Iterator{
		fields:  mc.values,
		current: 0,
	}
}
