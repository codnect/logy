package logy

type ObjectMarshaler interface {
	MarshalObject(encoder ObjectEncoder) error
}

type ArrayMarshaler interface {
	MarshalArray(encoder ArrayEncoder) error
}
