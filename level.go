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

var (
	levelValues = map[string]Level{
		"TRACE": LevelTrace,
		"DEBUG": LevelDebug,
		"INFO":  LevelInfo,
		"WARN":  LevelWarn,
		"ERROR": LevelError,
		"ALL":   LevelAll,
		"OFF":   LevelOff,
	}
)

const (
	LevelOff Level = iota + 1
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
	LevelAll
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
	case LevelError:
		return "ERROR"
	case LevelOff:
		return "OFF"
	default:
		return "ALL"
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
	if val, ok := levelValues[strings.ToUpper(node.Value)]; ok {
		*l = val
	} else {
		*l = LevelAll
	}

	return nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	var level string
	if err := json.Unmarshal(data, &level); err != nil {
		return err
	}

	if val, ok := levelValues[strings.ToUpper(level)]; ok {
		*l = val
	} else {
		*l = LevelAll
	}

	return nil
}
