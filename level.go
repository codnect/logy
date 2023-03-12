package logy

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"strings"
)

type syslogLevel int

const (
	syslogLevelEmergency syslogLevel = iota
	syslogLevelAlert
	syslogLevelCritical
	syslogLevelError
	syslogLevelWarning
	syslogLevelNotice
	syslogLevelInformational
	syslogLevelDebug
)

type Level int

const (
	LevelError Level = iota + 1
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (l Level) syslogLevel() syslogLevel {
	switch l {
	case LevelTrace:
		return syslogLevelNotice
	case LevelDebug:
		return syslogLevelDebug
	case LevelInfo:
		return syslogLevelInformational
	case LevelWarn:
		return syslogLevelWarning
	default:
		return syslogLevelError
	}
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	default:
		return "ERROR"
	}
}

func (l *Level) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('"')
	builder.WriteString(l.String())
	builder.WriteByte('"')
	return []byte(builder.String()), nil
}

func (l *Level) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}

func (l *Level) UnmarshalYAML(node *yaml.Node) error {
	switch node.Value {
	case "TRACE":
		*l = LevelTrace
	case "DEBUG":
		*l = LevelDebug
	case "INFO":
		*l = LevelInfo
	case "WARN":
		*l = LevelWarn
	case "ERROR":
		*l = LevelError
	}

	return nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch val {
	case "TRACE":
		*l = LevelTrace
	case "DEBUG":
		*l = LevelDebug
	case "INFO":
		*l = LevelInfo
	case "WARN":
		*l = LevelWarn
	case "ERROR":
		*l = LevelError
	}

	return nil
}
