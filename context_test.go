package logy

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutValue(t *testing.T) {
	ctx := WithContextFields(context.Background())
	PutValue(ctx, "anyKey", "anyValue")
	PutValue(ctx, "anotherKey", "anotherValue")

	ctxFields := ContextFieldsFrom(ctx)

	expectedFields := []Field{
		{
			key:       "anyKey",
			value:     "anyValue",
			textValue: "anyValue",
			jsonValue: "{\"anyKey\":\"anyValue\"}",
		},
		{
			key:       "anotherKey",
			value:     "anotherValue",
			textValue: "anotherValue",
			jsonValue: "{\"anotherKey\":\"anotherValue\"}",
		},
	}

	assert.Equal(t, expectedFields, ctxFields.fields)
}

func TestPutValue_PanicsIfContextIsNil(t *testing.T) {
	assert.Panics(t, func() {
		PutValue(nil, "anyKey", "anyValue")
	})
}

func TestPutValue_PanicsIfKeyIsEmpty(t *testing.T) {
	assert.Panics(t, func() {
		PutValue(context.Background(), "", "anyValue")
	})
}

func TestWithValue(t *testing.T) {
	ctx := WithContextFields(context.Background())
	ctx = WithValue(ctx, "anyKey", "anyValue")
	ctx = WithValue(ctx, "anotherKey", "anotherValue")

	ctxFields := ContextFieldsFrom(ctx)

	expectedFields := []Field{
		{
			key:       "anyKey",
			value:     "anyValue",
			textValue: "anyValue",
			jsonValue: "{\"anyKey\":\"anyValue\"}",
		},
		{
			key:       "anotherKey",
			value:     "anotherValue",
			textValue: "anotherValue",
			jsonValue: "{\"anotherKey\":\"anotherValue\"}",
		},
	}

	assert.Equal(t, expectedFields, ctxFields.fields)
}

func TestIterator_HasNext(t *testing.T) {

	ctx := WithContextFields(context.Background())
	ctx = WithValue(ctx, "anyKey", "anyValue")
	ctx = WithValue(ctx, "anotherKey", "anotherValue")

	fields := 0
	iterator := Values(ctx)
	for {
		if !iterator.HasNext() {
			break
		}

		iterator.Next()
		fields++
	}

	assert.Equal(t, 2, fields)
}

func TestIterator_Next(t *testing.T) {
	ctx := WithContextFields(context.Background())
	ctx = WithValue(ctx, "anyKey", "anyValue")
	ctx = WithValue(ctx, "anotherKey", "anotherValue")

	expectedFields := []Field{
		{
			key:       "anyKey",
			value:     "anyValue",
			textValue: "anyValue",
			jsonValue: "{\"anyKey\":\"anyValue\"}",
		},
		{
			key:       "anotherKey",
			value:     "anotherValue",
			textValue: "anotherValue",
			jsonValue: "{\"anotherKey\":\"anotherValue\"}",
		},
	}

	fields := make([]Field, 0)
	iterator := Values(ctx)
	for {
		field, next := iterator.Next()
		if !next {
			break
		}
		fields = append(fields, field)
	}

	assert.Equal(t, expectedFields, fields)
}

func TestWithValue_ShouldOverrideFieldIfItAlreadyExists(t *testing.T) {
	ctx := WithValue(context.Background(), "anyKey", "anyValue")

	expectedFields := []Field{
		{
			key:       "anyKey",
			value:     "anyValue",
			textValue: "anyValue",
			jsonValue: "{\"anyKey\":\"anyValue\"}",
		},
	}

	fields := make([]Field, 0)
	iterator := Values(ctx)
	for {
		field, next := iterator.Next()
		if !next {
			break
		}
		fields = append(fields, field)
	}

	assert.Equal(t, expectedFields, fields)

	cloneContext := WithValue(ctx, "anyKey", "anotherValue")

	expectedFields = []Field{
		{
			key:       "anyKey",
			value:     "anotherValue",
			textValue: "anotherValue",
			jsonValue: "{\"anyKey\":\"anotherValue\"}",
		},
	}

	fields = make([]Field, 0)
	iterator = Values(cloneContext)
	for {
		field, next := iterator.Next()
		if !next {
			break
		}
		fields = append(fields, field)
	}

	assert.Equal(t, expectedFields, fields)
}

func TestContextFieldsFrom_ShouldReturnNilIfContextIsNil(t *testing.T) {
	assert.Nil(t, ContextFieldsFrom(nil))
}

func TestWithValue_ShouldPanicIfParentContextIsNil(t *testing.T) {
	assert.Panics(t, func() {
		WithValue(nil, "", "")
	})
}

func TestValues_ShouldPanicIfContextIsNil(t *testing.T) {
	assert.Panics(t, func() {
		Values(nil)
	})
}

func TestValues_ShouldReturnDefaultIteratorIfContextDoesNotIncludeContextFields(t *testing.T) {
	assert.Equal(t, defaultIterator, Values(context.Background()))
}
