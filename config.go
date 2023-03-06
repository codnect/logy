package logy

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

const (
	defaultLogFormat   = "%d %p %c : %m%s%n"
	defaultLogFileName = "logy.log"
	defaultLogFilePath = "."

	PropertyLevel   = "level"
	PropertyEnabled = "enabled"
)

var (
	config = &Config{
		Level:         LevelInfo,
		IncludeCaller: false,
		Handlers:      []string{"console"},
		Console: &ConsoleConfig{
			Enable: true,
			Target: TargetStderr,
			Format: defaultLogFormat,
			Color:  true,
			Level:  LevelDebug,
			Json:   nil,
		},
		File: &FileConfig{
			Name:   defaultLogFileName,
			Enable: false,
			Path:   defaultLogFilePath,
			Format: defaultLogFormat,
			Level:  LevelDebug,
			Json:   nil,
		},
		Package:          map[string]*PackageConfig{},
		ExternalHandlers: map[string]ConfigProperties{},
	}
	reservedPropertiesKey = []string{"level", "include-caller", "handlers", "console", "file", "syslog", "package"}
	configMu              sync.RWMutex
)

type ConfigProperties map[string]any

type Target string

const (
	TargetStderr  Target = "stderr"
	TargetStdout  Target = "stdout"
	TargetDiscard Target = "discard"
)

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

type ExcludedKeys []string

func (ek *ExcludedKeys) UnmarshalJSON(data []byte) error {
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
		*ek = ExcludedKeys{}
		for _, item := range val.([]any) {
			*ek = append(*ek, strings.TrimSpace(item.(string)))
		}
		return nil
	}

	items := strings.Split(val.(string), ",")
	if len(items) == 0 {
		return nil
	}

	*ek = ExcludedKeys{}
	for _, item := range items {
		*ek = append(*ek, strings.TrimSpace(item))
	}

	return nil
}

type JsonAdditionalFields map[string]any

type JsonConfig struct {
	ExcludeKeys      ExcludedKeys         `json:"exclude-keys" xml:"exclude-keys" yaml:"exclude-keys"`
	AdditionalFields JsonAdditionalFields `json:"additional-fields" xml:"additional-fields" yaml:"additional-fields"`
}

type ConsoleConfig struct {
	Enable bool        `json:"enable" xml:"enable" yaml:"enable"`
	Target Target      `json:"target" xml:"target" yaml:"target"`
	Format string      `json:"format" xml:"format" yaml:"format"`
	Color  bool        `json:"color" xml:"color" yaml:"color"`
	Level  Level       `json:"level" xml:"level" yaml:"level"`
	Json   *JsonConfig `json:"json" xml:"json" yaml:"json"`
}

type FileConfig struct {
	Name   string      `json:"name" xml:"name" yaml:"name"`
	Enable bool        `json:"enable" xml:"enable" yaml:"enable"`
	Path   string      `json:"path" xml:"path" yaml:"path"`
	Format string      `json:"format" xml:"format" yaml:"format"`
	Level  Level       `json:"level" xml:"level" yaml:"level"`
	Json   *JsonConfig `json:"json" xml:"json" yaml:"json"`
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

	if cfg.File != nil && cfg.File.Enable {
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

		if cfg.File.Level == 0 {
			cfg.File.Level = config.File.Level
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

func configureHandlers(config *Config) error {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	for name, handler := range handlers {
		if name == "console" {
			console, ok := handler.(*ConsoleHandler)

			if !ok {
				continue
			}

			err := console.onConfigure(config.Console)
			if err != nil {
				return err
			}

			continue
		}

		if name == "file" {
			console, ok := handler.(*FileHandler)

			if !ok {
				continue
			}

			err := console.onConfigure(config.File)
			if err != nil {
				return err
			}

			continue
		}

		if cfg, ok := config.ExternalHandlers[name]; ok {
			configurable, isConfigurable := handler.(ConfigurableHandler)

			if !isConfigurable {
				level, ok := cfg[PropertyLevel]

				if ok {
					switch level.(type) {
					case int:
						handler.SetLevel(Level(level.(int)))
					}
				}

				var enabled any
				enabled, ok = cfg[PropertyEnabled]
				if ok {
					switch enabled.(type) {
					case bool:
						handler.SetEnabled(enabled.(bool))
					}
				}

				continue
			}

			configurable.OnConfigure(cfg)
		}
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
