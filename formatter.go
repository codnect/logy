package slog

import "fmt"

type Formatter interface {
	Format(record *Record) string
}

type SimpleFormatter struct {
}

func NewSimpleFormatter() *SimpleFormatter {
	return &SimpleFormatter{}
}

func (f *SimpleFormatter) Format(record *Record) string {
	return fmt.Sprintf("%s %s", record.Level.String(), record.Message)
}
