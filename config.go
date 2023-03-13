package logy

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

const (
	DefaultTextFormat = "%d{2006-01-02 15:04:05.000000} %p %c : %m%s%n"

	DefaultLogFileName = "logy.log"
	DefaultLogFilePath = "."

	DefaultSyslogFormat   = "%d %p %m%s%n"
	DefaultSyslogEndpoint = "localhost:514"

	PropertyLevel   = "level"
	PropertyEnabled = "enabled"
)

var (
	config = &Config{
		Level:         LevelInfo,
		IncludeCaller: false,
		Handlers:      []string{"console"},
		Console: &ConsoleConfig{
			Enabled: true,
			Target:  TargetStderr,
			Format:  DefaultTextFormat,
			Color:   true,
			Level:   LevelDebug,
			Json:    nil,
		},
		File: &FileConfig{
			Name:    DefaultLogFileName,
			Enabled: false,
			Path:    DefaultLogFilePath,
			Format:  DefaultTextFormat,
			Level:   LevelDebug,
			Json:    nil,
		},
		Syslog: &SyslogConfig{
			Enabled:  false,
			Endpoint: DefaultSyslogEndpoint,
			AppName:  os.Args[0],
			Hostname: "",
			Facility: FacilityUserLevel,
			LogType:  RFC5424,
			Protocol: ProtocolTCP,
			Format:   DefaultSyslogFormat,
			Level:    LevelTrace,
		},
		Package:          map[string]*PackageConfig{},
		ExternalHandlers: map[string]ConfigProperties{},
	}
	reservedPropertiesKey = []string{"level", "include-caller", "handlers", "console", "file", "syslog", "package"}
	configMu              sync.RWMutex
)

type ConfigProperties map[string]any

type Target int

const (
	TargetStderr Target = iota + 1
	TargetStdout
	TargetDiscard
)

var (
	targetValues = map[string]Target{
		"stderr":  TargetStderr,
		"stdout":  TargetStdout,
		"discard": TargetDiscard,
	}
)

func (t *Target) String() string {
	switch *t {
	case TargetStderr:
		return "stderr"
	case TargetStdout:
		return "stdout"
	default:
		return "discard"
	}
}

func (t *Target) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('"')
	builder.WriteString(t.String())
	builder.WriteByte('"')
	return []byte(builder.String()), nil
}

func (t *Target) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *Target) UnmarshalJSON(data []byte) error {
	var target string
	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}

	if val, ok := targetValues[strings.ToLower(target)]; ok {
		*t = val
	} else {
		*t = TargetDiscard
	}

	return nil
}

func (t *Target) UnmarshalYAML(node *yaml.Node) error {
	if val, ok := targetValues[strings.ToLower(node.Value)]; ok {
		*t = val
	} else {
		*t = TargetDiscard
	}

	return nil
}

type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

func (p *Protocol) String() string {
	return string(*p)
}

func (p *Protocol) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('"')
	builder.WriteString(p.String())
	builder.WriteByte('"')
	return []byte(builder.String()), nil
}

func (p *Protocol) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}

func (p *Protocol) UnmarshalJSON(data []byte) error {
	var protocol string
	if err := json.Unmarshal(data, &protocol); err != nil {
		return err
	}

	switch strings.ToLower(protocol) {
	case "udp":
		*p = ProtocolUDP
	default:
		*p = ProtocolTCP
	}

	return nil
}

func (p *Protocol) UnmarshalYAML(node *yaml.Node) error {
	switch strings.ToLower(node.Value) {
	case "udp":
		*p = ProtocolUDP
	default:
		*p = ProtocolTCP
	}

	return nil
}

type SysLogType int

const (
	RFC5424 SysLogType = iota + 1
	RFC3164
)

var (
	syslogTypes = map[string]SysLogType{
		"rfc3164": RFC3164,
		"rfc5424": RFC5424,
	}
)

func (t *SysLogType) String() string {
	switch *t {
	case RFC3164:
		return "rfc3164"
	case RFC5424:
		return "rfc5424"
	default:
		return ""
	}
}

func (t *SysLogType) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteByte('"')
	builder.WriteString(t.String())
	builder.WriteByte('"')
	return []byte(builder.String()), nil
}

func (t *SysLogType) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *SysLogType) UnmarshalJSON(data []byte) error {
	var typ string
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	if val, ok := syslogTypes[strings.ToLower(typ)]; ok {
		*t = val
	} else {
		*t = RFC5424
	}

	return nil
}

func (t *SysLogType) UnmarshalYAML(node *yaml.Node) error {
	if val, ok := syslogTypes[strings.ToLower(node.Value)]; ok {
		*t = val
	} else {
		*t = RFC5424
	}

	return nil
}

type Facility int

const (
	FacilityKernel Facility = iota + 1
	FacilityUserLevel
	FacilityMailSystem
	FacilitySystemDaemons
	FacilitySecurity
	FacilitySyslogd
	FacilityLinePrinter
	FacilityNetworkNews
	FacilityUUCP
	FacilityClockDaemon
	FacilitySecurity2
	FacilityFTPDaemon
	FacilityNTP
	FacilityLogAudit
	FacilityLogAlert
	FacilityClockDaemon2
	FacilityLocalUse0
	FacilityLocalUse1
	FacilityLocalUse2
	FacilityLocalUse3
	FacilityLocalUse4
	FacilityLocalUse5
	FacilityLocalUse6
	FacilityLocalUse7
)

var (
	facilityValues = map[string]Facility{
		"kernel":         FacilityKernel,
		"user-level":     FacilityUserLevel,
		"mail-system":    FacilityMailSystem,
		"system-daemons": FacilitySystemDaemons,
		"security":       FacilitySecurity,
		"syslogd":        FacilitySyslogd,
		"line-printer":   FacilityLinePrinter,
		"network-news":   FacilityNetworkNews,
		"uucp":           FacilityUUCP,
		"clock-daemon":   FacilityClockDaemon,
		"security2":      FacilitySecurity2,
		"ftp-daemon":     FacilityFTPDaemon,
		"ntp":            FacilityNTP,
		"log-audit":      FacilityLogAudit,
		"log-alert":      FacilityLogAlert,
		"local-use-0":    FacilityLocalUse0,
		"local-use-1":    FacilityLocalUse1,
		"local-use-2":    FacilityLocalUse2,
		"local-use-3":    FacilityLocalUse3,
		"local-use-4":    FacilityLocalUse4,
		"local-use-5":    FacilityLocalUse5,
		"local-use-6":    FacilityLocalUse6,
		"local-use-7":    FacilityLocalUse7,
	}
)

func (f *Facility) UnmarshalJSON(data []byte) error {
	var facility string
	if err := json.Unmarshal(data, &facility); err != nil {
		return err
	}

	if val, ok := facilityValues[strings.ToLower(facility)]; ok {
		*f = val
	} else {
		*f = FacilityUserLevel
	}

	return nil
}

func (f *Facility) UnmarshalYAML(node *yaml.Node) error {
	if val, ok := facilityValues[strings.ToLower(node.Value)]; ok {
		*f = val
	} else {
		*f = FacilityUserLevel
	}

	return nil
}

type Handlers []string

func (h *Handlers) UnmarshalJSON(data []byte) error {
	var (
		val     any
		isArray bool
	)

	if data[0] == '[' || data[len(data)-1] == ']' {
		val = make([]string, 0)
		isArray = true
	} else {
		val = ""
	}

	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	if isArray {
		*h = Handlers{}
		for _, item := range val.([]any) {
			*h = append(*h, strings.TrimSpace(item.(string)))
		}
		return nil
	}

	items := strings.Split(val.(string), ",")
	if len(items) == 0 {
		return nil
	}

	*h = Handlers{}
	for _, item := range items {
		*h = append(*h, strings.TrimSpace(item))
	}

	return nil
}

func (h *Handlers) UnmarshalYAML(node *yaml.Node) error {
	if len(node.Content) != 0 {
		*h = Handlers{}
		for _, item := range node.Content {
			*h = append(*h, strings.TrimSpace(item.Value))
		}
	}

	return nil
}

type KeyOverrides map[string]string
type JsonAdditionalFields map[string]any

type JsonConfig struct {
	KeyOverrides     KeyOverrides         `json:"key-overrides" xml:"key-overrides" yaml:"key-overrides"`
	AdditionalFields JsonAdditionalFields `json:"additional-fields" xml:"additional-fields" yaml:"additional-fields"`
}

type ConsoleConfig struct {
	Enabled bool        `json:"enabled" xml:"enabled" yaml:"enabled"`
	Target  Target      `json:"target" xml:"target" yaml:"target"`
	Format  string      `json:"format" xml:"format" yaml:"format"`
	Color   bool        `json:"color" xml:"color" yaml:"color"`
	Level   Level       `json:"level" xml:"level" yaml:"level"`
	Json    *JsonConfig `json:"json" xml:"json" yaml:"json"`
}

type FileConfig struct {
	Name    string      `json:"name" xml:"name" yaml:"name"`
	Enabled bool        `json:"enabled" xml:"enabled" yaml:"enabled"`
	Path    string      `json:"path" xml:"path" yaml:"path"`
	Format  string      `json:"format" xml:"format" yaml:"format"`
	Level   Level       `json:"level" xml:"level" yaml:"level"`
	Json    *JsonConfig `json:"json" xml:"json" yaml:"json"`
}

type SyslogConfig struct {
	Enabled          bool       `json:"enabled" xml:"enabled" yaml:"enabled"`
	Endpoint         string     `json:"endpoint" xml:"endpoint" yaml:"endpoint"`
	AppName          string     `json:"app-name" xml:"app-name" yaml:"app-name"`
	Hostname         string     `json:"hostname" xml:"hostname" yaml:"hostname"`
	Facility         Facility   `json:"facility" xml:"facility" yaml:"facility"`
	LogType          SysLogType `json:"log-type" xml:"log-type" yaml:"log-type"`
	Protocol         Protocol   `json:"protocol" xml:"protocol" yaml:"protocol"`
	Format           string     `json:"format" xml:"format" yaml:"format"`
	Level            Level      `json:"level" xml:"level" yaml:"level"`
	BlockOnReconnect bool       `json:"block-on-reconnect" xml:"block-on-reconnect" yaml:"block-on-reconnect"`
}

type PackageConfig struct {
	Level             Level    `json:"level" xml:"level" yaml:"level"`
	UseParentHandlers bool     `json:"use-parent-handlers" xml:"use-parent-handlers" yaml:"use-parent-handlers"`
	Handlers          Handlers `json:"handlers" xml:"handlers" yaml:"handlers"`
}

type Config struct {
	Level            Level                       `json:"level" xml:"level" yaml:"level"`
	IncludeCaller    bool                        `json:"include-caller" xml:"include-caller" yaml:"include-caller"`
	Handlers         Handlers                    `json:"handlers" xml:"handlers" yaml:"handlers"`
	Console          *ConsoleConfig              `json:"console" xml:"console" yaml:"console"`
	File             *FileConfig                 `json:"file" xml:"file" yaml:"file"`
	Syslog           *SyslogConfig               `json:"syslog" xml:"syslog" yaml:"syslog"`
	Package          map[string]*PackageConfig   `json:"package" xml:"package" yaml:"package"`
	ExternalHandlers map[string]ConfigProperties `json:"-" xml:"-" yaml:"-"`
}

func isReservedPropertyKey(key string) bool {
	for _, reserved := range reservedPropertiesKey {
		if reserved == key {
			return true
		}
	}

	return false
}

func loadConfigFromEnv() {
	cfgMap := map[string]any{}

	env := os.Environ()
	for _, variable := range env {
		kv := strings.SplitN(variable, "=", 2)
		if strings.HasPrefix(kv[0], "logy.") {
			key := strings.TrimSpace(kv[0])
			key = key[5:]

			parts := strings.Split(key, ".")
			if len(parts) == 1 {
				cfgMap[key] = kv[1]
			}

			properties := cfgMap
			for index, part := range parts {
				if m, exists := properties[part]; exists {
					properties = m.(map[string]any)
				} else if index != len(parts)-1 {
					temp := map[string]any{}
					properties[part] = temp
					properties = temp
				}

				if index == len(parts)-1 {
					properties[part] = kv[1]
				}
			}
		}
	}

	data, _ := json.Marshal(cfgMap)

	cfg := &Config{
		Package:          map[string]*PackageConfig{},
		ExternalHandlers: map[string]ConfigProperties{},
	}

	_ = json.Unmarshal(data, cfg)

	for key, val := range cfgMap {
		if isReservedPropertyKey(key) {
			continue
		}

		cfg.ExternalHandlers[key] = val.(map[string]any)
	}

	_ = LoadConfig(cfg)
}

func LoadConfigFromYaml(name string) error {
	yamlFile, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return err
	}

	cfgMap := map[string]any{}

	for key, value := range data {
		if key == "logy" {
			switch typed := value.(type) {
			case map[string]any:
				for propertyName, propertyValue := range typed {
					cfgMap[propertyName] = propertyValue
				}
			}
		}
	}

	cfg := &Config{
		Package:          map[string]*PackageConfig{},
		ExternalHandlers: map[string]ConfigProperties{},
	}

	byteData, _ := yaml.Marshal(cfgMap)

	err = yaml.Unmarshal(byteData, &cfg)
	if err != nil {
		return err
	}

	for key, val := range cfgMap {
		if isReservedPropertyKey(key) {
			continue
		}

		cfg.ExternalHandlers[key] = val.(map[string]any)
	}

	return LoadConfig(cfg)
}

func LoadConfig(cfg *Config) error {
	defer configMu.Unlock()
	configMu.Lock()

	if cfg == nil {
		return errors.New("config cannot be nil")
	}

	if cfg.Level == 0 {
		cfg.Level = config.Level
	}

	if cfg.Handlers == nil {
		cfg.Handlers = config.Handlers
	}

	if cfg.File != nil && cfg.File.Enabled {
		fileHandlerExists := false
		for _, handler := range cfg.Handlers {
			if handler == "file" {
				fileHandlerExists = true
			}
		}

		if !fileHandlerExists {
			cfg.Handlers = append(cfg.Handlers, "file")
		}
	}

	err := initializePackageConfig(cfg)
	if err != nil {
		return err
	}

	err = initializeConsoleConfig(cfg)
	if err != nil {
		return err
	}

	err = initializeFileConfig(cfg)
	if err != nil {
		return err
	}

	err = initializeSyslogConfig(cfg)
	if err != nil {
		return err
	}

	config = cfg

	err = configureHandlers(cfg)
	if err != nil {
		return err
	}

	return configureLoggers()
}

func initializePackageConfig(cfg *Config) error {
	if cfg.Package == nil {
		cfg.Package = config.Package
	}

	for pkg, pkgCfg := range cfg.Package {
		if strings.TrimSpace(pkg) == "" {
			return errors.New("package cannot be empty or blank")
		}

		if pkgCfg.Level == 0 {
			rootLevel := cfg.Level

			if rootLevel == 0 {
				rootLevel = config.Level
			}

			pkgCfg.Level = rootLevel
		}

		if pkgCfg.Handlers == nil && len(pkgCfg.Handlers) == 0 {
			rootHandlers := cfg.Handlers
			if rootHandlers == nil {
				rootHandlers = config.Handlers
			}

			pkgCfg.Handlers = rootHandlers
			pkgCfg.UseParentHandlers = true
		}
	}

	return nil
}

func initializeConsoleConfig(cfg *Config) error {

	if cfg.Console == nil {
		cfg.Console = config.Console
	} else {
		if cfg.Console.Level == 0 {
			cfg.Console.Level = config.Console.Level
		}

		if cfg.Console.Target == 0 {
			cfg.Console.Target = config.Console.Target
		}

		if strings.TrimSpace(cfg.Console.Format) == "" {
			cfg.Console.Format = config.Console.Format
		}

		if cfg.Console.Json == nil {
			cfg.Console.Json = config.Console.Json
		}
	}

	return nil
}

func initializeFileConfig(cfg *Config) error {

	if cfg.File == nil {
		cfg.File = config.File
	} else {
		if cfg.File.Level == 0 {
			cfg.File.Level = config.File.Level
		}

		if strings.TrimSpace(cfg.File.Format) == "" {
			cfg.File.Format = config.File.Format
		}

		if cfg.File.Name == "" {
			cfg.File.Name = config.File.Name
		}

		if cfg.File.Path == "" {
			cfg.File.Path = config.File.Path
		}

		if cfg.File.Json == nil {
			cfg.File.Json = config.File.Json
		}
	}

	return nil
}

func initializeSyslogConfig(cfg *Config) error {

	if cfg.Syslog == nil {
		cfg.Syslog = config.Syslog
	} else {
		if cfg.Syslog.Level == 0 {
			cfg.Syslog.Level = config.Syslog.Level
		}

		if strings.TrimSpace(cfg.Syslog.Format) == "" {
			cfg.Syslog.Format = config.Syslog.Format
		}

		if cfg.Syslog.AppName == "" {
			cfg.Syslog.AppName = config.Syslog.AppName
		}

		if cfg.Syslog.Hostname == "" {
			cfg.Syslog.Hostname = config.Syslog.Hostname
		}

		if cfg.Syslog.Protocol == "" {
			cfg.Syslog.Protocol = config.Syslog.Protocol
		}

		if cfg.Syslog.LogType == 0 {
			cfg.Syslog.LogType = config.Syslog.LogType
		}

		if cfg.Syslog.Facility == 0 {
			cfg.Syslog.Facility = config.Syslog.Facility
		}

		if cfg.Syslog.Endpoint == "" {
			cfg.Syslog.Endpoint = config.Syslog.Endpoint
		}
	}

	return nil
}

func configureHandlers(config *Config) error {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	for _, handler := range handlers {
		configurable, isConfigurable := handler.(ConfigurableHandler)

		if !isConfigurable {
			continue
		}

		_ = configurable.OnConfigure(*config)
	}

	return nil
}

func configureLoggers() error {
	defer loggerCacheMu.Unlock()
	loggerCacheMu.Lock()

	defer handlerMu.Unlock()
	handlerMu.Lock()

	rootLogger.onConfigure(config)
	return nil
}
