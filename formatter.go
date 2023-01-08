package logy

import "fmt"

type Formatter interface {
	Format(record Record) string
}

type SimpleFormatter struct {
}

func NewSimpleFormatter() *SimpleFormatter {
	return &SimpleFormatter{}
}

func (f *SimpleFormatter) Format(record Record) string {
	return fmt.Sprintf("%s [%10s] %s %s : %s", record.Time.Format("2006-01-02 15:04:05.000"), record.LoggerName, record.Level.String(), record.Caller.Package(), record.Message)
}
