package logy

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTextEncoder_Namespaces(t *testing.T) {
	encoder := getTextEncoder(newBuffer())
	encoder.OpenNamespace("anyNamespace")
	encoder.AddString("anyStringKey", "anyStringValue")
	encoder.CloseNamespace()
	assert.Equal(t, "\"anyNamespace\":{\"anyStringKey\":\"anyStringValue\"}", encoder.buf.String())
}

func TestTextEncoder_AddObject(t *testing.T) {
	encoder := getTextEncoder(newBuffer())
	encoder.AddObject("anyObject", &testObject{Name: "anyName"})
	assert.Equal(t, "\"anyObject\":{\"name\":\"anyName\"}", encoder.buf.String())
}

func TestTextEncoder_AppendObject(t *testing.T) {
	encoder := getTextEncoder(newBuffer())
	encoder.AppendObject(&testObject{Name: "anyName"})
	assert.Equal(t, "{\"name\":\"anyName\"}", encoder.buf.String())
}

func TestTextEncoder_AddArray(t *testing.T) {
	encoder := getTextEncoder(newBuffer())
	encoder.AddArray("anyArray", testArray{41, 11, 52})
	assert.Equal(t, "\"anyArray\":[41,11,52]", encoder.buf.String())
}

func TestTextEncoder_AppendArray(t *testing.T) {
	encoder := getTextEncoder(newBuffer())
	encoder.AppendArray(testArray{41, 11, 52})
	assert.Equal(t, "[41,11,52]", encoder.buf.String())
}

func TestTextEncoder_AddAny(t *testing.T) {
	anyDate := time.Now()
	anotherDate := time.Now()

	testCases := []struct {
		Key      string
		Value    any
		Expected any
	}{
		{
			Key:      "anyKey",
			Value:    "anyStringValue",
			Expected: "\"anyKey\":\"anyStringValue\"",
		},
		{
			Key:      "anyKey",
			Value:    int8(41),
			Expected: "\"anyKey\":41",
		},
		{
			Key:      "anyKey",
			Value:    int16(11),
			Expected: "\"anyKey\":11",
		},
		{
			Key:      "anyKey",
			Value:    int32(75),
			Expected: "\"anyKey\":75",
		},
		{
			Key:      "anyKey",
			Value:    int64(156),
			Expected: "\"anyKey\":156",
		},
		{
			Key:      "anyKey",
			Value:    617,
			Expected: "\"anyKey\":617",
		},
		{
			Key:      "anyKey",
			Value:    uint8(41),
			Expected: "\"anyKey\":41",
		},
		{
			Key:      "anyKey",
			Value:    uint16(11),
			Expected: "\"anyKey\":11",
		},
		{
			Key:      "anyKey",
			Value:    uint32(75),
			Expected: "\"anyKey\":75",
		},
		{
			Key:      "anyKey",
			Value:    uint64(156),
			Expected: "\"anyKey\":156",
		},
		{
			Key:      "anyKey",
			Value:    uint(617),
			Expected: "\"anyKey\":617",
		},
		{
			Key:      "anyKey",
			Value:    true,
			Expected: "\"anyKey\":true",
		},
		{
			Key:      "anyKey",
			Value:    false,
			Expected: "\"anyKey\":false",
		},
		{
			Key:      "anyKey",
			Value:    float32(25.7),
			Expected: "\"anyKey\":25.7",
		},
		{
			Key:      "anyKey",
			Value:    72.8,
			Expected: "\"anyKey\":72.8",
		},
		{
			Key:      "anyKey",
			Value:    anyDate,
			Expected: fmt.Sprintf("\"anyKey\":\"%s\"", anyDate.Format(time.RFC3339)),
		},
		{
			Key:      "anyKey",
			Value:    time.Hour * 5,
			Expected: "\"anyKey\":18000000000000",
		},
		{
			Key:      "anyKey",
			Value:    []string{"anyStringValue", "anotherStringValue"},
			Expected: "\"anyKey\":[\"anyStringValue\",\"anotherStringValue\"]",
		},
		{
			Key:      "anyKey",
			Value:    []int8{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []int16{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []int32{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []int64{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []int{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uint8{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uint16{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uint32{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uint64{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uint{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []uintptr{41, 100, 3},
			Expected: "\"anyKey\":[41,100,3]",
		},
		{
			Key:      "anyKey",
			Value:    []bool{true, false},
			Expected: "\"anyKey\":[true,false]",
		},
		{
			Key:      "anyKey",
			Value:    []float32{41.5, 100.6, 3.1},
			Expected: "\"anyKey\":[41.5,100.6,3.1]",
		},
		{
			Key:      "anyKey",
			Value:    []float64{41.5, 100.6, 3.1},
			Expected: "\"anyKey\":[41.5,100.6,3.1]",
		},
		{
			Key:      "anyKey",
			Value:    []complex64{41.5, 100.6, 3.1},
			Expected: "\"anyKey\":[\"41.5+0i\",\"100.6+0i\",\"3.1+0i\"]",
		},
		{
			Key:      "anyKey",
			Value:    []complex128{41.5, 100.6, 3.1},
			Expected: "\"anyKey\":[\"41.5+0i\",\"100.6+0i\",\"3.1+0i\"]",
		},
		{
			Key:      "anyKey",
			Value:    []time.Time{anyDate, anotherDate},
			Expected: fmt.Sprintf("\"anyKey\":[\"%s\",\"%s\"]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339)),
		},
		{
			Key:      "anyKey",
			Value:    []time.Duration{time.Second * 32, time.Hour * 5},
			Expected: "\"anyKey\":[32000000000,18000000000000]",
		},
		{
			Key:      "anyKey",
			Value:    []error{errors.New("anyError1"), errors.New("anyError2")},
			Expected: "\"anyKey\":[\"anyError1\",\"anyError2\"]",
		},
		{
			Key:      "anyKey",
			Value:    &testObject{Name: "anyName"},
			Expected: "\"anyKey\":{\"name\":\"anyName\"}",
		},
		{
			Key:      "anyKey",
			Value:    testArray{41, 11, 52},
			Expected: "\"anyKey\":[41,11,52]",
		},
		{
			Key:      "anyKey",
			Value:    &testStringer{},
			Expected: "\"anyKey\":\"anyValue\"",
		},
		{
			Key: "anyKey",
			Value: map[string]any{
				"anyMapKey": map[string]any{
					"anySubMapKey": "anySubMapValue",
				},
			},
			Expected: "\"anyKey\":{\"anyMapKey\":{\"anySubMapKey\":\"anySubMapValue\"}}",
		},
		{
			Key:      "anyKey",
			Value:    []any{"anyStringValue", 37},
			Expected: "\"anyKey\":[\"anyStringValue\",37]",
		},
		{
			Key:      "anyJsonObject",
			Value:    testJsonMarshaler{Name: "anyName"},
			Expected: "\"anyJsonObject\":{\"anyKey\":\"anyName\"}",
		},
	}

	encoder := getTextEncoder(newBuffer())

	for _, testCase := range testCases {
		encoder.buf.Reset()
		encoder.AddAny(testCase.Key, testCase.Value)
		assert.Equal(t, testCase.Expected, encoder.buf.String())
	}

}

func TestTextEncoder_AppendAny(t *testing.T) {
	anyDate := time.Now()
	anotherDate := time.Now()

	testCases := []struct {
		Value    any
		Expected any
	}{
		{
			Value:    "anyStringValue",
			Expected: "anyStringValue",
		},
		{
			Value:    int8(41),
			Expected: "41",
		},
		{
			Value:    int16(11),
			Expected: "11",
		},
		{
			Value:    int32(75),
			Expected: "75",
		},
		{
			Value:    int64(156),
			Expected: "156",
		},
		{
			Value:    617,
			Expected: "617",
		},
		{
			Value:    uint8(41),
			Expected: "41",
		},
		{
			Value:    uint16(11),
			Expected: "11",
		},
		{
			Value:    uint32(75),
			Expected: "75",
		},
		{
			Value:    uint64(156),
			Expected: "156",
		},
		{
			Value:    uint(617),
			Expected: "617",
		},
		{
			Value:    true,
			Expected: "true",
		},
		{
			Value:    false,
			Expected: "false",
		},
		{
			Value:    float32(25.7),
			Expected: "25.7",
		},
		{
			Value:    72.8,
			Expected: "72.8",
		},
		{
			Value:    anyDate,
			Expected: fmt.Sprintf("%s", anyDate.Format(time.RFC3339)),
		},
		{
			Value:    time.Hour * 5,
			Expected: "18000000000000",
		},
		{
			Value:    []string{"anyStringValue", "anotherStringValue"},
			Expected: "[anyStringValue,anotherStringValue]",
		},
		{
			Value:    []int8{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []int16{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []int32{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []int64{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []int{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uint8{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uint16{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uint32{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uint64{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uint{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []uintptr{41, 100, 3},
			Expected: "[41,100,3]",
		},
		{
			Value:    []bool{true, false},
			Expected: "[true,false]",
		},
		{
			Value:    []float32{41.5, 100.6, 3.1},
			Expected: "[41.5,100.6,3.1]",
		},
		{
			Value:    []float64{41.5, 100.6, 3.1},
			Expected: "[41.5,100.6,3.1]",
		},
		{
			Value:    []complex64{41.5, 100.6, 3.1},
			Expected: "[41.5+0i,100.6+0i,3.1+0i]",
		},
		{
			Value:    []complex128{41.5, 100.6, 3.1},
			Expected: "[41.5+0i,100.6+0i,3.1+0i]",
		},
		{
			Value:    []time.Time{anyDate, anotherDate},
			Expected: fmt.Sprintf("[%s,%s]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339)),
		},
		{
			Value:    []time.Duration{time.Second * 32, time.Hour * 5},
			Expected: "[32000000000,18000000000000]",
		},
		{
			Value:    []error{errors.New("anyError1"), errors.New("anyError2")},
			Expected: "[anyError1,anyError2]",
		},
		{
			Value:    &testObject{Name: "anyName"},
			Expected: "{\"name\":\"anyName\"}",
		},
		{
			Value:    testArray{41, 11, 52},
			Expected: "[41,11,52]",
		},
		{
			Value:    &testStringer{},
			Expected: "anyValue",
		},
		{
			Value: map[string]any{
				"anyMapKey": map[string]any{
					"anySubMapKey": "anySubMapValue",
				},
			},
			Expected: "{anyMapKey={anySubMapKey=anySubMapValue}}",
		},
		{
			Value:    []any{"anyStringValue", 37},
			Expected: "[anyStringValue,37]",
		},
		{
			Value:    testJsonMarshaler{Name: "anyName"},
			Expected: "{\"anyKey\":\"anyName\"}",
		},
	}

	encoder := getTextEncoder(newBuffer())

	for _, testCase := range testCases {
		encoder.buf.Reset()
		encoder.AppendAny(testCase.Value)
		assert.Equal(t, testCase.Expected, encoder.buf.String())
	}

}

func TestTextEncoder_AddByteString(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddByteString("anyKey", []byte("anyStringValue"))
	assert.Equal(t, "\"anyKey\":\"anyStringValue\"", encoder.buf.String())

	encoder.AddByteString("anotherKey", []byte("anotherStringValue"))
	assert.Equal(t, "\"anyKey\":\"anyStringValue\",\"anotherKey\":\"anotherStringValue\"", encoder.buf.String())
}

func TestTextEncoder_AddUTF8ByteString(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddByteString("anyKey", []byte("ŞşÇçğĞüÜ\"\\\b\f\n\r\t"))
	assert.Equal(t, "\"anyKey\":\"ŞşÇçğĞüÜ\\\"\\\\\\b\\f\\n\\r\\t\"", encoder.buf.String())
}

func TestTextEncoder_AddString(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddString("anyKey", "anyStringValue")
	assert.Equal(t, "\"anyKey\":\"anyStringValue\"", encoder.buf.String())

	encoder.AddString("anotherKey", "anotherStringValue")
	assert.Equal(t, "\"anyKey\":\"anyStringValue\",\"anotherKey\":\"anotherStringValue\"", encoder.buf.String())
}

func TestTextEncoder_AddUTF8String(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddString("anyKey", "ŞşÇçğĞüÜ\"\\\b\f\n\r\t")
	assert.Equal(t, "\"anyKey\":\"ŞşÇçğĞüÜ\\\"\\\\\\b\\f\\n\\r\\t\"", encoder.buf.String())
}

func TestTextEncoder_AppendString(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendString("anyStringValue")
	assert.Equal(t, "anyStringValue", encoder.buf.String())

	encoder.AppendString("anotherStringValue")
	assert.Equal(t, "anyStringValueanotherStringValue", encoder.buf.String())
}

func TestTextEncoder_AddStrings(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddStrings("anyKey", []string{"anyStringValue"})
	assert.Equal(t, "\"anyKey\":[\"anyStringValue\"]", encoder.buf.String())

	encoder.AddStrings("anotherKey", []string{"anotherStringValue"})
	assert.Equal(t, "\"anyKey\":[\"anyStringValue\"],\"anotherKey\":[\"anotherStringValue\"]", encoder.buf.String())
}

func TestTextEncoder_AppendStrings(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendStrings([]string{"anyStringValue"})
	assert.Equal(t, "[anyStringValue]", encoder.buf.String())

	encoder.AppendStrings([]string{"anotherStringValue"})
	assert.Equal(t, "[anyStringValue][anotherStringValue]", encoder.buf.String())
}

func TestTextEncoder_AddInt8(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt8("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddInt8("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendInt8(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt8(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendInt8(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddInt8s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt8s("anyKey", []int8{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInt8s("anotherKey", []int8{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AddInt78s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt8s("anyKey", []int8{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInt8s("anotherKey", []int8{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AddInt16(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt16("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddInt16("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendInt16(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt16(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendInt16(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddInt16s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt16s("anyKey", []int16{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInt16s("anotherKey", []int16{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendInt16s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt16s([]int16{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendInt16s([]int16{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddInt32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt32("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddInt32("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AddInt32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt32s("anyKey", []int32{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInt32s("anotherKey", []int32{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendInt32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt32(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendInt32(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AppendInt32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt32s([]int32{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendInt32s([]int32{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddInt64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt64("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddInt64("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendInt64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt64(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendInt64(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddInt64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt64s("anyKey", []int64{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInt64s("anotherKey", []int64{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendInt64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt64s([]int64{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendInt64s([]int64{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddInt(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInt("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddInt("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendInt(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInt(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendInt(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddInts(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddInts("anyKey", []int{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddInts("anotherKey", []int{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendInts(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendInts([]int{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendInts([]int{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddUint8(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint8("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUint8("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUint8(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint8(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUint8(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUint8s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint8s("anyKey", []uint8{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUint8s("anotherKey", []uint8{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUint8s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint8s([]uint8{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUint8s([]uint8{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddUint16(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint16("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUint16("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUint16(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint16(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUint16(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUint16s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint16s("anyKey", []uint16{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUint16s("anotherKey", []uint16{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUint16s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint16s([]uint16{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUint16s([]uint16{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddUint32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint32("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUint32("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUint32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint32(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUint32(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUint32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint32s("anyKey", []uint32{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUint32s("anotherKey", []uint32{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUint32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint32s([]uint32{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUint32s([]uint32{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddUint64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint64("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUint64("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUint64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint64(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUint64(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUint64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint64s("anyKey", []uint64{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUint64s("anotherKey", []uint64{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUint64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint64s([]uint64{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUint64s([]uint64{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddUint(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUint("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUint("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUint(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUint(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUint(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUints(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUints("anyKey", []uint{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUints("anotherKey", []uint{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUints(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUints([]uint{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUints([]uint{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddBool(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddBool("anyKey", true)
	assert.Equal(t, "\"anyKey\":true", encoder.buf.String())

	encoder.AddBool("anotherKey", false)
	assert.Equal(t, "\"anyKey\":true,\"anotherKey\":false", encoder.buf.String())
}

func TestTextEncoder_AppendBool(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendBool(true)
	assert.Equal(t, "true", encoder.buf.String())

	encoder.AppendBool(false)
	assert.Equal(t, "truefalse", encoder.buf.String())
}

func TestTextEncoder_AddBools(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddBools("anyKey", []bool{true, false})
	assert.Equal(t, "\"anyKey\":[true,false]", encoder.buf.String())

	encoder.AddBools("anotherKey", []bool{false, false})
	assert.Equal(t, "\"anyKey\":[true,false],\"anotherKey\":[false,false]", encoder.buf.String())
}

func TestTextEncoder_AppendBools(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendBools([]bool{true, false})
	assert.Equal(t, "[true,false]", encoder.buf.String())

	encoder.AppendBools([]bool{false, false})
	assert.Equal(t, "[true,false][false,false]", encoder.buf.String())
}

func TestTextEncoder_AddFloat32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddFloat32("anyKey", 41.5)
	assert.Equal(t, "\"anyKey\":41.5", encoder.buf.String())

	encoder.AddFloat32("anotherKey", 11.7)
	assert.Equal(t, "\"anyKey\":41.5,\"anotherKey\":11.7", encoder.buf.String())
}

func TestTextEncoder_AppendFloat32(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendFloat32(41.5)
	assert.Equal(t, "41.5", encoder.buf.String())

	encoder.AppendFloat32(11.7)
	assert.Equal(t, "41.511.7", encoder.buf.String())
}

func TestTextEncoder_AddFloat32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddFloat32s("anyKey", []float32{41.5, 100.6, 3.1})
	assert.Equal(t, "\"anyKey\":[41.5,100.6,3.1]", encoder.buf.String())

	encoder.AddFloat32s("anotherKey", []float32{11.8, 34.12})
	assert.Equal(t, "\"anyKey\":[41.5,100.6,3.1],\"anotherKey\":[11.8,34.12]", encoder.buf.String())
}

func TestTextEncoder_AppendFloat32s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendFloat32s([]float32{41.5, 100.6, 3.1})
	assert.Equal(t, "[41.5,100.6,3.1]", encoder.buf.String())

	encoder.AppendFloat32s([]float32{11.8, 34.12})
	assert.Equal(t, "[41.5,100.6,3.1][11.8,34.12]", encoder.buf.String())
}

func TestTextEncoder_AddFloat64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddFloat64("anyKey", 41.5)
	assert.Equal(t, "\"anyKey\":41.5", encoder.buf.String())

	encoder.AddFloat64("anotherKey", 11.7)
	assert.Equal(t, "\"anyKey\":41.5,\"anotherKey\":11.7", encoder.buf.String())
}

func TestTextEncoder_AppendFloat64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendFloat64(41.5)
	assert.Equal(t, "41.5", encoder.buf.String())

	encoder.AppendFloat64(11.7)
	assert.Equal(t, "41.511.7", encoder.buf.String())
}

func TestTextEncoder_AddFloat64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddFloat64s("anyKey", []float64{41.5, 100.6, 3.1})
	assert.Equal(t, "\"anyKey\":[41.5,100.6,3.1]", encoder.buf.String())

	encoder.AddFloat64s("anotherKey", []float64{11.8, 34.12})
	assert.Equal(t, "\"anyKey\":[41.5,100.6,3.1],\"anotherKey\":[11.8,34.12]", encoder.buf.String())
}

func TestTextEncoder_AppendFloat64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendFloat64s([]float64{41.5, 100.6, 3.1})
	assert.Equal(t, "[41.5,100.6,3.1]", encoder.buf.String())

	encoder.AppendFloat64s([]float64{11.8, 34.12})
	assert.Equal(t, "[41.5,100.6,3.1][11.8,34.12]", encoder.buf.String())
}

func TestTextEncoder_AddError(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddError("anyKey", errors.New("anyError"))
	assert.Equal(t, "\"anyKey\":\"anyError\"", encoder.buf.String())

	encoder.AddError("anyKey", errors.New("anotherError"))
	assert.Equal(t, "\"anyKey\":\"anyError\",\"anyKey\":\"anotherError\"", encoder.buf.String())
}

func TestTextEncoder_AppendError(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendError(errors.New("anyError"))
	assert.Equal(t, "anyError", encoder.buf.String())

	encoder.AppendError(errors.New("anotherError"))
	assert.Equal(t, "anyErroranotherError", encoder.buf.String())
}

func TestTextEncoder_AddErrors(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddErrors("anyKey", []error{errors.New("anyError1"), errors.New("anyError2")})
	assert.Equal(t, "\"anyKey\":[\"anyError1\",\"anyError2\"]", encoder.buf.String())

	encoder.AddErrors("anotherKey", []error{errors.New("anyError3"), errors.New("anyError4")})
	assert.Equal(t, "\"anyKey\":[\"anyError1\",\"anyError2\"],\"anotherKey\":[\"anyError3\",\"anyError4\"]", encoder.buf.String())
}

func TestTextEncoder_AppendErrors(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendErrors([]error{errors.New("anyError1"), errors.New("anyError2")})
	assert.Equal(t, "[anyError1,anyError2]", encoder.buf.String())

	encoder.AppendErrors([]error{errors.New("anyError3"), errors.New("anyError4")})
	assert.Equal(t, "[anyError1,anyError2][anyError3,anyError4]", encoder.buf.String())
}

func TestTextEncoder_AddTime(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	date := time.Now()
	anotherDate := time.Now()

	encoder.AddTime("anyKey", date)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":\"%s\"", date.Format(time.RFC3339)), encoder.buf.String())

	encoder.AddTime("anotherKey", anotherDate)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":\"%s\",\"anotherKey\":\"%s\"", date.Format(time.RFC3339), anotherDate.Format(time.RFC3339)), encoder.buf.String())
}

func TestTextEncoder_AppendTime(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	date := time.Now()
	anotherDate := time.Now()

	encoder.AppendTime(date)
	assert.Equal(t, fmt.Sprintf("%s", date.Format(time.RFC3339)), encoder.buf.String())

	encoder.AppendTime(anotherDate)
	assert.Equal(t, fmt.Sprintf("%s%s", date.Format(time.RFC3339), anotherDate.Format(time.RFC3339)), encoder.buf.String())
}

func TestTextEncoder_AddTimes(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	anyDate := time.Now()
	anotherDate := time.Now()

	encoder.AddTimes("anyKey", []time.Time{anyDate, anotherDate})
	assert.Equal(t, fmt.Sprintf("\"anyKey\":[\"%s\",\"%s\"]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339)), encoder.buf.String())

	current := time.Now()

	encoder.AddTimes("anotherKey", []time.Time{current})
	assert.Equal(t, fmt.Sprintf("\"anyKey\":[\"%s\",\"%s\"],\"anotherKey\":[\"%s\"]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339), current.Format(time.RFC3339)), encoder.buf.String())
}

func TestTextEncoder_AppendTimes(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	anyDate := time.Now()
	anotherDate := time.Now()

	encoder.AppendTimes([]time.Time{anyDate, anotherDate})
	assert.Equal(t, fmt.Sprintf("[%s,%s]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339)), encoder.buf.String())

	current := time.Now()

	encoder.AppendTimes([]time.Time{current})
	assert.Equal(t, fmt.Sprintf("[%s,%s][%s]", anyDate.Format(time.RFC3339), anotherDate.Format(time.RFC3339), current.Format(time.RFC3339)), encoder.buf.String())
}

func TestTextEncoder_AddTimeLayout(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	date := time.Now()
	anotherDate := time.Now()

	encoder.AddTimeLayout("anyKey", date, time.Stamp)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":\"%s\"", date.Format(time.Stamp)), encoder.buf.String())

	encoder.AddTimeLayout("anotherKey", anotherDate, time.Stamp)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":\"%s\",\"anotherKey\":\"%s\"", date.Format(time.Stamp), anotherDate.Format(time.Stamp)), encoder.buf.String())
}

func TestTextEncoder_AppendTimeLayout(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	date := time.Now()
	anotherDate := time.Now()

	encoder.AppendTimeLayout(date, time.Stamp)
	assert.Equal(t, fmt.Sprintf("%s", date.Format(time.Stamp)), encoder.buf.String())

	encoder.AppendTimeLayout(anotherDate, time.Stamp)
	assert.Equal(t, fmt.Sprintf("%s%s", date.Format(time.Stamp), anotherDate.Format(time.Stamp)), encoder.buf.String())
}

func TestTextEncoder_AddTimesLayout(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	anyDate := time.Now()
	anotherDate := time.Now()

	encoder.AddTimesLayout("anyKey", []time.Time{anyDate, anotherDate}, time.Stamp)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":[\"%s\",\"%s\"]", anyDate.Format(time.Stamp), anotherDate.Format(time.Stamp)), encoder.buf.String())

	current := time.Now()

	encoder.AddTimesLayout("anotherKey", []time.Time{current}, time.Stamp)
	assert.Equal(t, fmt.Sprintf("\"anyKey\":[\"%s\",\"%s\"],\"anotherKey\":[\"%s\"]", anyDate.Format(time.Stamp), anotherDate.Format(time.Stamp), current.Format(time.Stamp)), encoder.buf.String())
}

func TestTextEncoder_AppendTimesLayout(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	anyDate := time.Now()
	anotherDate := time.Now()

	encoder.AppendTimesLayout([]time.Time{anyDate, anotherDate}, time.Stamp)
	assert.Equal(t, fmt.Sprintf("[%s,%s]", anyDate.Format(time.Stamp), anotherDate.Format(time.Stamp)), encoder.buf.String())

	current := time.Now()

	encoder.AppendTimesLayout([]time.Time{current}, time.Stamp)
	assert.Equal(t, fmt.Sprintf("[%s,%s][%s]", anyDate.Format(time.Stamp), anotherDate.Format(time.Stamp), current.Format(time.Stamp)), encoder.buf.String())

}

func TestTextEncoder_AddDuration(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddDuration("anyKey", time.Second*35)
	assert.Equal(t, "\"anyKey\":35000000000", encoder.buf.String())

	encoder.AddDuration("anotherKey", time.Second*441)
	assert.Equal(t, "\"anyKey\":35000000000,\"anotherKey\":441000000000", encoder.buf.String())
}

func TestTextEncoder_AppendDuration(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendDuration(time.Second * 35)
	assert.Equal(t, "35000000000", encoder.buf.String())

	encoder.AppendDuration(time.Second * 441)
	assert.Equal(t, "35000000000441000000000", encoder.buf.String())
}

func TestTextEncoder_AddDurations(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddDurations("anyKey", []time.Duration{time.Second * 32, time.Hour * 5})
	assert.Equal(t, "\"anyKey\":[32000000000,18000000000000]", encoder.buf.String())

	encoder.AddDurations("anotherKey", []time.Duration{time.Minute * 8})
	assert.Equal(t, "\"anyKey\":[32000000000,18000000000000],\"anotherKey\":[480000000000]", encoder.buf.String())
}

func TestTextEncoder_AppendDurations(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendDurations([]time.Duration{time.Second * 32, time.Hour * 5})
	assert.Equal(t, "[32000000000,18000000000000]", encoder.buf.String())

	encoder.AppendDurations([]time.Duration{time.Minute * 8})
	assert.Equal(t, "[32000000000,18000000000000][480000000000]", encoder.buf.String())
}

func TestTextEncoder_AddUintptr(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUintptr("anyKey", 41)
	assert.Equal(t, "\"anyKey\":41", encoder.buf.String())

	encoder.AddUintptr("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":41,\"anotherKey\":11", encoder.buf.String())
}

func TestTextEncoder_AppendUintptr(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUintptr(41)
	assert.Equal(t, "41", encoder.buf.String())

	encoder.AppendUintptr(11)
	assert.Equal(t, "4111", encoder.buf.String())
}

func TestTextEncoder_AddUintptrs(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddUintptrs("anyKey", []uintptr{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[41,100,3]", encoder.buf.String())

	encoder.AddUintptrs("anotherKey", []uintptr{11, 34})
	assert.Equal(t, "\"anyKey\":[41,100,3],\"anotherKey\":[11,34]", encoder.buf.String())
}

func TestTextEncoder_AppendUintptrs(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendUintptrs([]uintptr{41, 100, 3})
	assert.Equal(t, "[41,100,3]", encoder.buf.String())

	encoder.AppendUintptrs([]uintptr{11, 34})
	assert.Equal(t, "[41,100,3][11,34]", encoder.buf.String())
}

func TestTextEncoder_AddComplex64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddComplex64("anyKey", 41)
	assert.Equal(t, "\"anyKey\":\"41+0i\"", encoder.buf.String())

	encoder.AddComplex64("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":\"41+0i\",\"anotherKey\":\"11+0i\"", encoder.buf.String())
}

func TestTextEncoder_AppendComplex64(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendComplex64(41)
	assert.Equal(t, "41+0i", encoder.buf.String())

	encoder.AppendComplex64(11)
	assert.Equal(t, "41+0i11+0i", encoder.buf.String())
}

func TestTextEncoder_AddComplex64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddComplex64s("anyKey", []complex64{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[\"41+0i\",\"100+0i\",\"3+0i\"]", encoder.buf.String())

	encoder.AddComplex64s("anotherKey", []complex64{11, 34})
	assert.Equal(t, "\"anyKey\":[\"41+0i\",\"100+0i\",\"3+0i\"],\"anotherKey\":[\"11+0i\",\"34+0i\"]", encoder.buf.String())
}

func TestTextEncoder_AppendComplex64s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendComplex64s([]complex64{41, 100, 3})
	assert.Equal(t, "[41+0i,100+0i,3+0i]", encoder.buf.String())

	encoder.AppendComplex64s([]complex64{11, 34})
	assert.Equal(t, "[41+0i,100+0i,3+0i][11+0i,34+0i]", encoder.buf.String())
}

func TestTextEncoder_AddComplex128(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddComplex128("anyKey", 41)
	assert.Equal(t, "\"anyKey\":\"41+0i\"", encoder.buf.String())

	encoder.AddComplex128("anotherKey", 11)
	assert.Equal(t, "\"anyKey\":\"41+0i\",\"anotherKey\":\"11+0i\"", encoder.buf.String())
}

func TestTextEncoder_AppendComplex128(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendComplex128(41)
	assert.Equal(t, "41+0i", encoder.buf.String())

	encoder.AppendComplex128(11)
	assert.Equal(t, "41+0i11+0i", encoder.buf.String())
}

func TestTextEncoder_AddComplex128s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AddComplex128s("anyKey", []complex128{41, 100, 3})
	assert.Equal(t, "\"anyKey\":[\"41+0i\",\"100+0i\",\"3+0i\"]", encoder.buf.String())

	encoder.AddComplex128s("anotherKey", []complex128{11, 34})
	assert.Equal(t, "\"anyKey\":[\"41+0i\",\"100+0i\",\"3+0i\"],\"anotherKey\":[\"11+0i\",\"34+0i\"]", encoder.buf.String())
}

func TestTextEncoder_AppendComplex128s(t *testing.T) {
	encoder := getTextEncoder(newBuffer())

	encoder.AppendComplex128s([]complex128{41, 100, 3})
	assert.Equal(t, "[41+0i,100+0i,3+0i]", encoder.buf.String())

	encoder.AppendComplex128s([]complex128{11, 34})
	assert.Equal(t, "[41+0i,100+0i,3+0i][11+0i,34+0i]", encoder.buf.String())
}
