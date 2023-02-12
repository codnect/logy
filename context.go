package logy

import "context"

const MappedContextKey = "logyMappedContext"

type MappedContext struct {
	values map[string]any
}

func NewMappedContext() *MappedContext {
	return &MappedContext{values: map[string]any{}}
}

func (mc *MappedContext) Keys() []string {
	keys := make([]string, len(mc.values))
	i := 0

	for key := range mc.values {
		keys[i] = key
		i++
	}

	return keys
}

func (mc *MappedContext) Put(key string, value any) {
	mc.values[key] = value
}

func (mc *MappedContext) Delete(key string) {
	delete(mc.values, key)
}

func (mc *MappedContext) Value(key string) any {
	return mc.values[key]
}

func MappedContextFrom(ctx context.Context) *MappedContext {
	return ctx.Value(MappedContextKey).(*MappedContext)
}

func WithMappedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, MappedContextKey, NewMappedContext())
}

func ToMappedContext(ctx context.Context, key string, value any) {
	mc := ctx.Value(MappedContextKey).(*MappedContext)
	mc.Put(key, value)
}

func FromMappedContext(ctx context.Context, key string) any {
	return ctx.Value(MappedContextKey).(*MappedContext)
}
