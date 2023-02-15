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

func (mc *MappedContext) Clone() *MappedContext {
	c := *mc
	return &c
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
		return context.WithValue(ctx, MappedContextKey, mc.Clone())
	}

	return context.WithValue(ctx, MappedContextKey, NewMappedContext())
}

func ToMappedContext(ctx context.Context, key string, value any) {
	mc := MappedContextFrom(ctx)

	if mc != nil {
		mc.Put(key, value)
	}
}

func FromMappedContext(ctx context.Context, key string) any {
	mc := MappedContextFrom(ctx)

	if mc == nil {
		return nil
	}

	return mc.Value(key)
}
