package logy

import (
	"errors"
	"strings"
	"sync"
)

const (
	defaultLogFormat = "%d{2006-01-02 15:04:05.000} %l %p : %m%n"

	PropertyLevel   = "level"
	PropertyEnabled = "enabled"
)

var (
	cfg   *Config
	cfgMu sync.RWMutex
)

type ConfigProperties map[string]any

type OutputTarget string

const (
	TargetStderr  OutputTarget = "stderr"
	TargetStdout  OutputTarget = "stdout"
	TargetDiscard OutputTarget = "discard"
)

type JsonFieldType string

const (
	FieldTypeString  JsonFieldType = "string"
	FieldTypeInteger JsonFieldType = "int"
)

type JsonAdditionalField struct {
	Value any           `json:"value" xml:"value" yaml:"value"`
	Type  JsonFieldType `json:"type" xml:"type" yaml:"type"`
}

type JsonConfig struct {
	ContextKeys      []string                       `json:"context-keys" xml:"context-keys" yaml:"context-keys"`
	ExcludeKeys      []string                       `json:"exclude-keys" xml:"exclude-keys" yaml:"exclude-keys"`
	AdditionalFields map[string]JsonAdditionalField `json:"additional-fields" xml:"additional-fields" yaml:"additional-fields"`
}

type ConsoleConfig struct {
	Enabled bool         `json:"enabled" xml:"enabled" yaml:"enabled"`
	Target  OutputTarget `json:"target" xml:"target" yaml:"target"`
	Format  string       `json:"format" xml:"format" yaml:"format"`
	Color   bool         `json:"color" xml:"color" yaml:"color"`
	Level   Level        `json:"level" xml:"level" yaml:"level"`
	Json    *JsonConfig  `json:"json" xml:"json" yaml:"json"`
}

type FileConfig struct {
	Enabled bool        `json:"enabled" xml:"enabled" yaml:"enabled"`
	Path    string      `json:"path" xml:"path" yaml:"path"`
	Format  string      `json:"format" xml:"format" yaml:"format"`
	Level   Level       `json:"level" xml:"level" yaml:"level"`
	Json    *JsonConfig `json:"json" xml:"json" yaml:"json"`
}

type PackageConfig struct {
	Level             Level    `json:"level" xml:"level" yaml:"level"`
	UseParentHandlers bool     `json:"use-parent-handlers" xml:"use-parent-handlers" yaml:"use-parent-handlers"`
	Handlers          []string `json:"handlers" xml:"handlers" yaml:"handlers"`
}

type Config struct {
	Level            Level                       `json:"level" xml:"level" yaml:"level"`
	Handlers         []string                    `json:"handlers" xml:"handlers" yaml:"handlers"`
	Console          *ConsoleConfig              `json:"console" xml:"console" yaml:"console"`
	File             *FileConfig                 `json:"file" xml:"file" yaml:"file"`
	Package          map[string]*PackageConfig   `json:"package" xml:"package" yaml:"package"`
	ExternalHandlers map[string]ConfigProperties `json:"-" xml:"-" yaml:"-"`
}

func LoadConfig(config *Config) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}

	if config.Level == 0 {
		config.Level = LevelInfo
	}

	enableConsole := false

	if config.Handlers == nil || len(config.Handlers) == 0 {
		config.Handlers = []string{"console"}
		enableConsole = true
	}

	err := initializePackageConfig(config)
	if err != nil {
		return err
	}

	err = initializeConsoleConfig(config, enableConsole)
	if err != nil {
		return err
	}

	err = initializeFileConfig(config)
	if err != nil {
		return err
	}

	defer cfgMu.Unlock()
	cfgMu.Lock()

	cfg = config

	configureHandlers(cfg)
	return configureLoggers()
}

func initializePackageConfig(config *Config) error {
	if config.Package == nil {
		config.Package = map[string]*PackageConfig{}
	}

	for pkg, pkgCfg := range config.Package {
		if strings.TrimSpace(pkg) == "" {
			return errors.New("package cannot be empty or blank")
		}

		if pkgCfg.Level == 0 {
			pkgCfg.Level = config.Level
		}

		if pkgCfg.Handlers == nil && len(pkgCfg.Handlers) == 0 {
			pkgCfg.Handlers = config.Handlers
			pkgCfg.UseParentHandlers = true
		}
	}

	return nil
}

func initializeConsoleConfig(config *Config, enableConsole bool) error {

	if config.Console == nil {
		config.Console = &ConsoleConfig{
			Enabled: true,
			Target:  TargetStderr,
			Format:  defaultLogFormat,
			Color:   true,
			Level:   LevelDebug,
			Json:    nil,
		}
	} else {
		if config.Console.Level == 0 {
			config.Console.Level = LevelDebug
		}

		if config.Console.Enabled && strings.TrimSpace(config.Console.Format) == "" {
			return errors.New("console.format cannot be empty or blank")
		}

		if enableConsole {
			config.Console.Enabled = true
			config.Console.Format = defaultLogFormat
		}
	}

	return nil
}

func initializeFileConfig(config *Config) error {

	if config.File == nil {
		config.File = &FileConfig{
			Enabled: false,
			Path:    "logy.log",
			Format:  defaultLogFormat,
			Level:   LevelInfo,
			Json:    nil,
		}
	} else {
		if config.File.Level == 0 {
			config.File.Level = LevelInfo
		}

		if config.File.Enabled && strings.TrimSpace(config.File.Format) == "" {
			return errors.New("file.format cannot be empty or blank")
		}
	}

	return nil
}

func configureHandlers(config *Config) {
	defer handlerMu.Unlock()
	handlerMu.Lock()

	for name, handler := range handlers {
		if name == "console" {
			console, ok := handler.(*ConsoleHandler)

			if !ok {
				continue
			}

			console.onConfigure(config.Console)
			continue
		}

		if cfg, ok := config.ExternalHandlers[name]; ok {
			configurable, isConfigurable := handler.(ConfigurableHandler)

			if !isConfigurable {
				level, exists := cfg[PropertyLevel]

				if exists {
					switch level.(type) {
					case int:
						handler.SetLevel(Level(level.(int)))
					}
				}

				enabled, exists := cfg[PropertyEnabled]
				if exists {
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
}

func configureLoggers() error {
	defer loggerCacheMu.Unlock()
	loggerCacheMu.Lock()

	defer handlerMu.Unlock()
	handlerMu.Lock()

	rootLogger.onConfigure(cfg)
	return nil
}
