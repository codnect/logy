package logy

import "context"

const ContextFieldsKey = "logyContextFields"

var (
	defaultIterator = &Iterator{}
	defaultField    = Field{}
)

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
	fields  *[]Field
	current int
}

func (i *Iterator) HasNext() bool {
	return i.current < len(*i.fields)
}

func (i *Iterator) Next() (Field, bool) {
	if i.current >= len(*i.fields) {
		return Field{}, false
	}

	field := (*i.fields)[i.current]
	i.current++
	return field, true
}

type ContextFields struct {
	fields []Field
}

func NewContextFields() *ContextFields {
	contextFields := &ContextFields{}
	return contextFields
}

func withField(field Field) *ContextFields {
	contextFields := NewContextFields()
	contextFields.fields = append(contextFields.fields, field)
	return contextFields
}

func (cf *ContextFields) Field(name string) (Field, bool) {
	for _, field := range cf.fields {
		if field.key == name {
			return field, true
		}
	}

	return defaultField, false
}

func (cf *ContextFields) IsEmpty() bool {
	return len(cf.fields) == 0
}

func (cf *ContextFields) Iterator() *Iterator {
	if len(cf.fields) == 0 {
		return defaultIterator
	}

	return &Iterator{
		fields:  &cf.fields,
		current: 0,
	}
}

func (cf *ContextFields) put(another Field) {
	added := false

	for i, field := range cf.fields {
		if field.key == another.key {
			cf.fields[i] = another
			added = true
		} else {
			cf.fields[i] = field
		}
	}

	if !added {
		cf.fields = append(cf.fields, another)
	}
}

func (cf *ContextFields) cloneWith(another Field) *ContextFields {
	clone := *cf
	added := false

	copyOfFields := make([]Field, len(cf.fields))

	for i, field := range cf.fields {
		if field.key == another.key {
			copyOfFields[i] = another
			added = true
		} else {
			copyOfFields[i] = field
		}
	}

	if !added {
		copyOfFields = append(copyOfFields, another)
	}

	clone.fields = copyOfFields
	return &clone
}

func WithContextFields(parent context.Context) context.Context {
	return context.WithValue(parent, ContextFieldsKey, NewContextFields())
}

func ContextFieldsFrom(ctx context.Context) *ContextFields {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(ContextFieldsKey)
	if val != nil {
		switch typed := val.(type) {
		case *ContextFields:
			return typed
		}
	}

	return nil
}

func PutValue(ctx context.Context, key string, value any) {
	if ctx == nil {
		panic("nil ctx")
	}

	if key == "" {
		panic("empty key")
	}

	val := ctx.Value(ContextFieldsKey)

	field := Field{
		key:   key,
		value: value,
	}

	if val != nil {
		switch typed := val.(type) {
		case *ContextFields:
			buf := newBuffer()
			defer buf.Free()

			jsonEncoder := getJSONEncoder(buf)
			jsonEncoder.buf.WriteByte('{')
			jsonEncoder.AddAny(key, value)
			jsonEncoder.buf.WriteByte('}')
			field.jsonValue = jsonEncoder.buf.String()

			putJSONEncoder(jsonEncoder)
			buf.Reset()

			textEncoder := getTextEncoder(buf)
			textEncoder.AppendAny(value)
			field.textValue = textEncoder.buf.String()

			putTextEncoder(textEncoder)
			buf.Reset()

			typed.put(field)
		}
	}
}

func WithValue(parent context.Context, key string, value any) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	val := parent.Value(ContextFieldsKey)

	field := Field{
		key:   key,
		value: value,
	}

	buf := newBuffer()
	defer buf.Free()

	jsonEncoder := getJSONEncoder(buf)
	jsonEncoder.buf.WriteByte('{')
	jsonEncoder.AddAny(key, value)
	jsonEncoder.buf.WriteByte('}')
	field.jsonValue = jsonEncoder.buf.String()

	putJSONEncoder(jsonEncoder)
	buf.Reset()

	textEncoder := getTextEncoder(buf)
	textEncoder.AppendAny(value)
	field.textValue = textEncoder.buf.String()

	putTextEncoder(textEncoder)
	buf.Reset()

	if val != nil {
		switch typed := val.(type) {
		case *ContextFields:
			return context.WithValue(parent, ContextFieldsKey, typed.cloneWith(field))
		}
	}

	return context.WithValue(parent, ContextFieldsKey, withField(field))
}

func Values(ctx context.Context) *Iterator {
	if ctx == nil {
		panic("nil ctx")
	}

	val := ctx.Value(ContextFieldsKey)

	if val != nil {
		switch typed := val.(type) {
		case *ContextFields:
			return typed.Iterator()
		}
	}

	return defaultIterator
}
