package logy

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"
)

var _textPool = sync.Pool{New: func() interface{} {
	return &textEncoder{jsonEncoder: &jsonEncoder{}}
}}

func getTextEncoder() *textEncoder {
	return _textPool.Get().(*textEncoder)
}

func putTextEncoder(enc *textEncoder) {
	enc.buf = nil
	_textPool.Put(enc)
}

type textEncoder struct {
	*jsonEncoder
}

func (enc *textEncoder) AppendAny(val any) error {
	switch typed := val.(type) {
	case string:
		enc.AppendString(typed)
	case int:
		enc.AppendInt(typed)
	case int8:
		enc.AppendInt8(typed)
	case int16:
		enc.AppendInt16(typed)
	case int32:
		enc.AppendInt32(typed)
	case int64:
		enc.AppendInt64(typed)
	case uint:
		enc.AppendUint(typed)
	case uint8:
		enc.AppendUint8(typed)
	case uint16:
		enc.AppendUint16(typed)
	case uint32:
		enc.AppendUint32(typed)
	case uint64:
		enc.AppendUint64(typed)
	case uintptr:
		enc.AppendUintptr(typed)
	case bool:
		enc.AppendBool(typed)
	case float32:
		enc.AppendFloat32(typed)
	case float64:
		enc.AppendFloat64(typed)
	case error:
		enc.AppendError(typed)
	case time.Time:
		enc.AppendTime(typed)
	case time.Duration:
		enc.AppendDuration(typed)
	case []string:
		enc.AppendStrings(typed)
	case []int:
		enc.AppendInts(typed)
	case []int8:
		enc.AppendInt8s(typed)
	case []int16:
		enc.AppendInt16s(typed)
	case []int32:
		enc.AppendInt32s(typed)
	case []int64:
		enc.AppendInt64s(typed)
	case []uint:
		enc.AppendUints(typed)
	case []uint8:
		enc.AppendUint8s(typed)
	case []uint16:
		enc.AppendUint16s(typed)
	case []uint32:
		enc.AppendUint32s(typed)
	case []uint64:
		enc.AppendUint64s(typed)
	case []uintptr:
		enc.AppendUintptrs(typed)
	case []bool:
		enc.AppendBools(typed)
	case []float32:
		enc.AppendFloat32s(typed)
	case []float64:
		enc.AppendFloat64s(typed)
	case []error:
		enc.AppendErrors(typed)
	case []time.Time:
		enc.AppendTimes(typed)
	case []time.Duration:
		enc.AppendDurations(typed)
	default:
		if marshaler, ok := typed.(ObjectMarshaler); ok {
			return enc.AppendObject(marshaler)
		} else if marshaler, ok := typed.(ArrayMarshaler); ok {
			return enc.AppendArray(marshaler)
		} else if stringer, ok := typed.(fmt.Stringer); ok {
			enc.AppendString(stringer.String())
		} else {
			rValue := reflect.ValueOf(typed)
			if rValue.Kind() == reflect.Pointer {
				return enc.AppendAny(rValue.Interface())
			} else if rValue.Kind() == reflect.Map {
				enc.appendMap(&rValue)
			} else if rValue.Kind() == reflect.Array || rValue.Kind() == reflect.Slice {
				enc.appendSlice(&rValue)
			} else {
				enc.AppendString(rValue.String())
			}
		}
	}

	return nil
}

func (enc *textEncoder) AppendObject(obj ObjectMarshaler) error {
	return enc.jsonEncoder.AppendObject(obj)
}

func (enc *textEncoder) appendMap(rValue *reflect.Value) {
	size := rValue.Len()
	enc.buf.WriteByte('{')
	i := 0

	for item := rValue.MapRange(); item.Next(); {
		enc.AppendAny(item.Key().Interface())
		enc.buf.WriteByte('=')
		enc.AppendAny(item.Value().Interface())
		if i != size-1 {
			enc.buf.WriteByte(',')
		}
		i++
	}
	enc.buf.WriteByte('}')
}

func (enc *textEncoder) appendSlice(rValue *reflect.Value) {
	size := rValue.Len()
	enc.buf.WriteByte('[')

	for i := 0; i < size; i++ {
		item := rValue.Index(i)

		enc.AppendAny(item.Interface())

		if i != size-1 {
			enc.buf.WriteByte(',')
		}
	}

	enc.buf.WriteByte(']')
}

/* append arrays */

func (enc *textEncoder) AppendStrings(arr []string) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendString(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendInts(arr []int) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt64(int64(item))
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendInt8s(arr []int8) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt8(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendInt16s(arr []int16) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt16(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendInt32s(arr []int32) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendInt64s(arr []int64) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUints(arr []uint) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUint8s(arr []uint8) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint8(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUint16s(arr []uint16) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint16(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUint32s(arr []uint32) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUint64s(arr []uint64) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendUintptrs(arr []uintptr) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUintptr(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendBools(arr []bool) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendBool(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendFloat32s(arr []float32) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendFloat32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendFloat64s(arr []float64) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendFloat64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendErrors(arr []error) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendError(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendTimes(arr []time.Time) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendTime(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendDurations(arr []time.Duration) {
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendDuration(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *textEncoder) AppendArray(arr ArrayMarshaler) error {
	enc.buf.WriteByte('[')
	err := arr.MarshalArray(enc)
	enc.buf.WriteByte(']')
	return err
}

/* append primitive */

func (enc *textEncoder) AppendString(val string) {
	enc.buf.WriteString(val)
}

func (enc *textEncoder) AppendByteString(val []byte) {
	_, _ = enc.buf.Write(val)
}

func (enc *textEncoder) AppendInt(v int)     { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt8(v int8)   { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt16(v int16) { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt32(v int32) { enc.AppendInt64(int64(v)) }
func (enc *textEncoder) AppendInt64(val int64) {
	enc.buf.WriteInt(val)
}

func (enc *textEncoder) AppendUint(v uint)     { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint8(v uint8)   { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint16(v uint16) { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint32(v uint32) { enc.AppendUint64(uint64(v)) }
func (enc *textEncoder) AppendUint64(val uint64) {
	enc.buf.WriteUint(val)
}

func (enc *textEncoder) AppendUintptr(v uintptr) { enc.AppendUint64(uint64(v)) }

func (enc *textEncoder) AppendBool(val bool) {
	enc.buf.WriteBool(val)
}

func (enc *textEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.WriteString("NaN")
	case math.IsInf(val, 1):
		enc.buf.WriteString("+Inf")
	case math.IsInf(val, -1):
		enc.buf.WriteString("-Inf")
	default:
		enc.buf.WriteFloat(val, bitSize)
	}
}
func (enc *textEncoder) AppendFloat64(v float64) { enc.appendFloat(v, 64) }
func (enc *textEncoder) AppendFloat32(v float32) { enc.appendFloat(float64(v), 32) }

func (enc *textEncoder) AppendError(val error) {
	enc.buf.WriteString(val.Error())
}

func (enc *textEncoder) AppendTime(t time.Time) {
	year, month, day := t.Date()
	enc.buf.WritePosIntWidth(year, 2)
	enc.buf.WriteByte('-')

	enc.buf.WritePosIntWidth(int(month), 2)
	enc.buf.WriteByte('-')

	enc.buf.WritePosIntWidth(day, 2)
	enc.buf.WriteByte(' ')

	hour, min, sec := t.Clock()
	enc.buf.WritePosIntWidth(hour, 2)
	enc.buf.WriteByte(':')
	enc.buf.WritePosIntWidth(min, 2)
	enc.buf.WriteByte(':')
	enc.buf.WritePosIntWidth(sec, 2)

	enc.buf.WriteByte('.')
	enc.buf.WritePosIntWidth(t.Nanosecond()/1e3, 6)
}

func (enc *textEncoder) AppendDuration(val time.Duration) {
	enc.AppendInt64(int64(val))
}

func (enc *textEncoder) appendComplex(val complex128, precision int) {

	r, i := float64(real(val)), float64(imag(val))
	enc.buf.WriteByte('"')
	enc.buf.WriteFloat(r, precision)
	if i >= 0 {
		enc.buf.WriteByte('+')
	}
	enc.buf.WriteFloat(i, precision)
	enc.buf.WriteByte('i')
	enc.buf.WriteByte('"')
}
func (enc *textEncoder) AppendComplex64(v complex64)   { enc.appendComplex(complex128(v), 32) }
func (enc *textEncoder) AppendComplex128(v complex128) { enc.appendComplex(complex128(v), 64) }

/* other */
