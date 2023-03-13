package logy

import (
	"context"
)

const ContextKey = "$logyMappedContext"

type Field struct {
	key       string
	value     any
	textValue string
	jsonValue string
}

func (f Field) Key() string {
	return f.key
}

func (f Field) Value() any {
	return f.value
}

func (f Field) ValueAsText() string {
	return f.textValue
}

func (f Field) AsJson() string {
	return f.jsonValue
}

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
	values      []Field
	keyIndexMap map[string]int
	textEncoder *textEncoder
	jsonEncoder *jsonEncoder
}

func NewMappedContext() *MappedContext {
	mc := &MappedContext{
		values:      []Field{},
		textEncoder: getTextEncoder(newBuffer()),
		jsonEncoder: getJSONEncoder(newBuffer()),
		keyIndexMap: map[string]int{},
	}

	return mc
}

func (mc *MappedContext) Size() int {
	return len(mc.values)
}

func (mc *MappedContext) resetEncoders() {
	mc.textEncoder.buf.Reset()
	mc.jsonEncoder.buf.Reset()
}

func (mc *MappedContext) putText(index int, value any) {
	mc.textEncoder.AppendAny(value)
	mc.values[index].textValue = mc.textEncoder.buf.String()
}

func (mc *MappedContext) putJson(index int, key string, value any) {
	mc.jsonEncoder.buf.WriteByte('{')
	mc.jsonEncoder.AddAny(key, value)
	mc.jsonEncoder.buf.WriteByte('}')
	mc.values[index].jsonValue = mc.jsonEncoder.buf.String()
}

func (mc *MappedContext) Put(key string, value any) {
	mc.clone()
	mc.resetEncoders()

	if index, ok := mc.keyIndexMap[key]; ok {
		mc.values[index].value = value
		mc.putText(index, value)
		mc.putJson(index, key, value)
	} else {
		field := Field{
			key,
			value,
			"",
			"",
		}

		index := len(mc.values)
		mc.keyIndexMap[key] = index
		mc.values = append(mc.values, field)

		mc.putText(index, value)
		mc.putJson(index, key, value)
	}

	mc.resetEncoders()
}

func (mc *MappedContext) Field(key string) (Field, bool) {
	if index, ok := mc.keyIndexMap[key]; ok {
		return mc.values[index], true
	}

	return Field{}, false
}

func (mc *MappedContext) Iterator() *Iterator {
	return &Iterator{
		fields:  mc.values,
		current: 0,
	}
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

	c.textEncoder = getTextEncoder(newBuffer())
	c.jsonEncoder = getJSONEncoder(newBuffer())
	return &c
}

func MappedContextFrom(ctx context.Context) *MappedContext {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(ContextKey)

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
		return context.WithValue(ctx, ContextKey, mc.clone())
	}

	return context.WithValue(ctx, ContextKey, NewMappedContext())
}

func PutValue(ctx context.Context, key string, val any) {
	ctxVal := ctx.Value(ContextKey)

	if ctxVal != nil {
		if mc, isMc := ctxVal.(*MappedContext); isMc {
			mc.Put(key, val)
		}
	}
}

func WithValue(parent context.Context, key string, value any) context.Context {
	ctx := WithMappedContext(parent)
	val := ctx.Value(ContextKey)

	if val != nil {
		if mc, isMc := val.(*MappedContext); isMc {
			mc.Put(key, value)
		}
	}

	return ctx
}

func ValueFrom(ctx context.Context, key string) (any, bool) {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return nil, false
	}

	field, ok := mc.Field(key)
	return field.Value(), ok
}

func IteratorFrom(ctx context.Context) *Iterator {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return &Iterator{}
	}

	return &Iterator{
		fields:  mc.values,
		current: 0,
	}
}
