package logy

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"
	"unicode/utf8"
)

var _jsonPool = sync.Pool{New: func() interface{} {
	return &jsonEncoder{}
}}

func getJSONEncoder() *jsonEncoder {
	return _jsonPool.Get().(*jsonEncoder)
}

func putJSONEncoder(enc *jsonEncoder) {
	enc.buf = nil
	enc.openNamespaces = 0
	_jsonPool.Put(enc)
}

type jsonEncoder struct {
	buf            *buffer
	openNamespaces int
}

/* object */

func (enc *jsonEncoder) AddAny(key string, val any) error {
	enc.addKey(key)
	return enc.AppendAny(val)
}

func (enc *jsonEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *jsonEncoder) AddStrings(key string, arr []string) {
	enc.addKey(key)
	enc.AppendStrings(arr)
}

func (enc *jsonEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *jsonEncoder) AddInt(k string, v int) { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInts(k string, v []int) {
	enc.addKey(k)
	enc.AppendInts(v)
}

func (enc *jsonEncoder) AddInt8(k string, v int8) { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt8s(k string, v []int8) {
	enc.addKey(k)
	enc.AppendInt8s(v)
}

func (enc *jsonEncoder) AddInt16(k string, v int16) { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt16s(k string, v []int16) {
	enc.addKey(k)
	enc.AppendInt16s(v)
}

func (enc *jsonEncoder) AddInt32(k string, v int32) { enc.AddInt64(k, int64(v)) }
func (enc *jsonEncoder) AddInt32s(k string, v []int32) {
	enc.addKey(k)
	enc.AppendInt32s(v)
}

func (enc *jsonEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *jsonEncoder) AddInt64s(k string, v []int64) {
	enc.addKey(k)
	enc.AppendInt64s(v)
}

func (enc *jsonEncoder) AddUint(k string, v uint) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUints(k string, v []uint) {
	enc.addKey(k)
	enc.AppendUints(v)
}

func (enc *jsonEncoder) AddUint8(k string, v uint8) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint8s(k string, v []uint8) {
	enc.addKey(k)
	enc.AppendUint8s(v)
}

func (enc *jsonEncoder) AddUint16(k string, v uint16) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint16s(k string, v []uint16) {
	enc.addKey(k)
	enc.AppendUint16s(v)
}

func (enc *jsonEncoder) AddUint32(k string, v uint32) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUint32s(k string, v []uint32) {
	enc.addKey(k)
	enc.AppendUint32s(v)
}

func (enc *jsonEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *jsonEncoder) AddUint64s(k string, v []uint64) {
	enc.addKey(k)
	enc.AppendUint64s(v)
}

func (enc *jsonEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *jsonEncoder) AddUintptrs(key string, val []uintptr) {
	enc.addKey(key)
	enc.AppendUintptrs(val)
}

func (enc *jsonEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *jsonEncoder) AddBools(key string, val []bool) {
	enc.addKey(key)
	enc.AppendBools(val)
}

func (enc *jsonEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
}

func (enc *jsonEncoder) AddFloat32s(key string, val []float32) {
	enc.addKey(key)
	enc.AppendFloat32s(val)
}

func (enc *jsonEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *jsonEncoder) AddFloat64s(key string, val []float64) {
	enc.addKey(key)
	enc.AppendFloat64s(val)
}

func (enc *jsonEncoder) AddError(key string, val error) {
	enc.addKey(key)
	enc.AppendError(val)
}

func (enc *jsonEncoder) AddErrors(key string, val []error) {
	enc.addKey(key)
	enc.AppendErrors(val)
}

func (enc *jsonEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *jsonEncoder) AddTimes(key string, val []time.Time) {
	enc.addKey(key)
	enc.AppendTimes(val)
}

func (enc *jsonEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *jsonEncoder) AddDurations(key string, val []time.Duration) {
	enc.addKey(key)
	enc.AppendDurations(val)
}

func (enc *jsonEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
}

func (enc *jsonEncoder) AddComplex64s(key string, val []complex64) {
	enc.addKey(key)
	enc.AppendComplex64s(val)
}

func (enc *jsonEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *jsonEncoder) AddComplex128s(key string, val []complex128) {
	enc.addKey(key)
	enc.AppendComplex128s(val)
}

func (enc *jsonEncoder) AddArray(key string, arr ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *jsonEncoder) AddObject(key string, obj ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

/* append */

func (enc *jsonEncoder) AppendAny(val any) error {
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
	case complex64:
		enc.AppendComplex64(typed)
	case complex128:
		enc.AppendComplex128(typed)
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
	case []complex64:
		enc.AppendComplex64s(typed)
	case []complex128:
		enc.AppendComplex128s(typed)
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
			}
		}
	}

	return nil
}

func (enc *jsonEncoder) AppendObject(obj ObjectMarshaler) error {
	old := enc.openNamespaces
	enc.openNamespaces = 0
	enc.addElementSeparator()
	enc.buf.WriteByte('{')
	err := obj.MarshalObject(enc)
	enc.buf.WriteByte('}')
	enc.CloseOpenNamespaces()
	enc.openNamespaces = old
	return err
}

func (enc *jsonEncoder) appendMap(rValue *reflect.Value) {
	enc.addElementSeparator()

	size := rValue.Len()
	enc.buf.WriteByte('{')
	i := 0

	for item := rValue.MapRange(); item.Next(); {
		enc.AppendAny(item.Key().Interface())
		enc.buf.WriteByte(':')
		enc.AppendAny(item.Value().Interface())

		if i != size-1 {
			enc.buf.WriteByte(',')
		}

		i++
	}
	enc.buf.WriteByte('}')
}

func (enc *jsonEncoder) appendSlice(rValue *reflect.Value) {
	enc.addElementSeparator()

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

func (enc *jsonEncoder) AppendStrings(arr []string) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendString(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendInts(arr []int) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt64(int64(item))
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendInt8s(arr []int8) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt8(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendInt16s(arr []int16) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt16(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendInt32s(arr []int32) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendInt64s(arr []int64) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendInt64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUints(arr []uint) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUint8s(arr []uint8) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint8(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUint16s(arr []uint16) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint16(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUint32s(arr []uint32) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUint64s(arr []uint64) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUint64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendUintptrs(arr []uintptr) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendUintptr(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendBools(arr []bool) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendBool(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendFloat32s(arr []float32) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendFloat32(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendFloat64s(arr []float64) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendFloat64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendErrors(arr []error) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendError(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendTimes(arr []time.Time) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendTime(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendComplex64s(arr []complex64) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendComplex64(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendComplex128s(arr []complex128) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendComplex128(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendDurations(arr []time.Duration) {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	for i, item := range arr {
		enc.AppendDuration(item)
		if i != len(arr)-1 {
			enc.buf.WriteByte(',')
		}
	}
	enc.buf.WriteByte(']')
}

func (enc *jsonEncoder) AppendArray(arr ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.WriteByte('[')
	err := arr.MarshalArray(enc)
	enc.buf.WriteByte(']')
	return err
}

/* append primitive */

func (enc *jsonEncoder) AppendString(val string) {
	enc.addElementSeparator()
	enc.buf.WriteByte('"')
	enc.safeAddString(val)
	enc.buf.WriteByte('"')
}

func (enc *jsonEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.WriteByte('"')
	enc.safeAddByteString(val)
	enc.buf.WriteByte('"')
}

func (enc *jsonEncoder) AppendInt(v int)     { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt8(v int8)   { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt16(v int16) { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt32(v int32) { enc.AppendInt64(int64(v)) }
func (enc *jsonEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.WriteInt(val)
}

func (enc *jsonEncoder) AppendUint(v uint)     { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint8(v uint8)   { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint16(v uint16) { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint32(v uint32) { enc.AppendUint64(uint64(v)) }
func (enc *jsonEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.WriteUint(val)
}

func (enc *jsonEncoder) AppendUintptr(v uintptr) { enc.AppendUint64(uint64(v)) }

func (enc *jsonEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.WriteBool(val)
}

func (enc *jsonEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.WriteString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.WriteString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.WriteString(`"-Inf"`)
	default:
		enc.buf.WriteFloat(val, bitSize)
	}
}
func (enc *jsonEncoder) AppendFloat64(v float64) { enc.appendFloat(v, 64) }
func (enc *jsonEncoder) AppendFloat32(v float32) { enc.appendFloat(float64(v), 32) }

func (enc *jsonEncoder) AppendError(val error) {
	enc.buf.WriteString(val.Error())
}

func (enc *jsonEncoder) AppendTime(t time.Time) {
	enc.addElementSeparator()
	enc.buf.WriteByte('"')
	*enc.buf = t.AppendFormat(*enc.buf, time.RFC3339)
	enc.buf.WriteByte('"')
}

func (enc *jsonEncoder) AppendDuration(val time.Duration) {
	enc.AppendInt64(int64(val))
}

func (enc *jsonEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()

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
func (enc *jsonEncoder) AppendComplex64(v complex64)   { enc.appendComplex(complex128(v), 32) }
func (enc *jsonEncoder) AppendComplex128(v complex128) { enc.appendComplex(complex128(v), 64) }

/* other */

func (enc *jsonEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.WriteByte('{')
	enc.openNamespaces++
}

func (enc *jsonEncoder) CloseNamespace() {
	enc.buf.WriteByte('}')
	enc.openNamespaces--
}

func (enc *jsonEncoder) CloseOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.WriteByte('}')
	}
	enc.openNamespaces = 0
}

func (enc *jsonEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.AppendString(key)
	enc.buf.WriteByte(':')
}

func (enc *jsonEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	default:
		enc.buf.WriteByte(',')
	}
}

func (enc *jsonEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.WriteString(s[i : i+size])
		i += size
	}
}

func (enc *jsonEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

func (enc *jsonEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.WriteByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.WriteByte('\\')
		enc.buf.WriteByte(b)
	case '\n':
		enc.buf.WriteByte('\\')
		enc.buf.WriteByte('n')
	case '\r':
		enc.buf.WriteByte('\\')
		enc.buf.WriteByte('r')
	case '\t':
		enc.buf.WriteByte('\\')
		enc.buf.WriteByte('t')
	default:
		enc.buf.WriteString(`\u00`)
		enc.buf.WriteByte(hex[b>>4])
		enc.buf.WriteByte(hex[b&0xF])
	}
	return true
}

func (enc *jsonEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.WriteString(`\ufffd`)
		return true
	}
	return false
}
