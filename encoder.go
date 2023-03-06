package logy

import (
	"time"
)

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

type ArrayEncoder interface {
	AppendAny(any) error
	AppendObject(ObjectMarshaler) error
	AppendArray(ArrayMarshaler) error

	AppendString(string)
	AppendStrings([]string)

	AppendByteString([]byte)

	AppendInt(int)
	AppendInts([]int)

	AppendInt8(int8)
	AppendInt8s([]int8)

	AppendInt16(int16)
	AppendInt16s([]int16)

	AppendInt32(int32)
	AppendInt32s([]int32)

	AppendInt64(int64)
	AppendInt64s([]int64)

	AppendUint(uint)
	AppendUints([]uint)

	AppendUint8(uint8)
	AppendUint8s([]uint8)

	AppendUint16(uint16)
	AppendUint16s([]uint16)

	AppendUint32(uint32)
	AppendUint32s([]uint32)

	AppendUint64(uint64)
	AppendUint64s([]uint64)

	AppendUintptr(uintptr)
	AppendUintptrs([]uintptr)

	AppendBool(bool)
	AppendBools([]bool)

	AppendFloat32(float32)
	AppendFloat32s([]float32)

	AppendFloat64(float64)
	AppendFloat64s([]float64)

	AppendError(error)
	AppendErrors([]error)

	AppendComplex64(complex64)
	AppendComplex64s([]complex64)
	AppendComplex128(complex128)
	AppendComplex128s([]complex128)

	AppendTime(time.Time)
	AppendTimes([]time.Time)
	AppendDuration(time.Duration)
	AppendDurations([]time.Duration)
}

type ObjectEncoder interface {
	AddAny(key string, val any) error
	AddObject(key string, marshaler ObjectMarshaler) error
	AddArray(key string, marshaler ArrayMarshaler) error

	AddString(key, value string)
	AddStrings(key string, arr []string)

	AddByteString(key string, arr []byte)

	AddInt(key string, value int)
	AddInts(key string, arr []int)

	AddInt8(key string, value int8)
	AddInt8s(key string, arr []int8)

	AddInt16(key string, value int16)
	AddInt16s(key string, arr []int16)

	AddInt32(key string, value int32)
	AddInt32s(key string, arr []int32)

	AddInt64(key string, value int64)
	AddInt64s(key string, arr []int64)

	AddUint(key string, value uint)
	AddUints(k string, v []uint)

	AddUint8(key string, value uint8)
	AddUint8s(k string, v []uint8)

	AddUint16(key string, value uint16)
	AddUint16s(k string, v []uint16)

	AddUint32(key string, value uint32)
	AddUint32s(k string, v []uint32)

	AddUint64(key string, value uint64)
	AddUint64s(k string, v []uint64)

	AddUintptr(key string, value uintptr)
	AddUintptrs(key string, arr []uintptr)

	AddBool(key string, value bool)
	AddBools(key string, value []bool)

	AddFloat32(key string, value float32)
	AddFloat32s(key string, val []float32)

	AddFloat64(key string, value float64)
	AddFloat64s(key string, val []float64)

	AddError(key string, val error)
	AddErrors(key string, val []error)

	AddTime(key string, value time.Time)
	AddTimes(key string, val []time.Time)

	AddDuration(key string, value time.Duration)
	AddDurations(key string, val []time.Duration)

	AddComplex64(key string, value complex64)
	AddComplex64s(key string, val []complex64)

	AddComplex128(key string, value complex128)
	AddComplex128s(key string, val []complex128)

	OpenNamespace(key string)
	CloseNamespace()
	CloseOpenNamespaces()
}
