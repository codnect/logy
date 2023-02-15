package logy

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

const (
	defaultLogFormat = "%d %p %c : %m%s%n"

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

type JsonAdditionalField struct {
	Value any `json:"value" xml:"value" yaml:"value"`
}

type JsonConfig struct {
	ExcludeKeys      []string                       `json:"exclude-keys" xml:"exclude-keys" yaml:"exclude-keys"`
	AdditionalFields map[string]JsonAdditionalField `json:"additional-fields" xml:"additional-fields" yaml:"additional-fields"`
}

type ConsoleConfig struct {
	Enable bool         `json:"enable" xml:"enable" yaml:"enable"`
	Target OutputTarget `json:"target" xml:"target" yaml:"target"`
	Format string       `json:"format" xml:"format" yaml:"format"`
	Color  bool         `json:"color" xml:"color" yaml:"color"`
	Level  Level        `json:"level" xml:"level" yaml:"level"`
	Json   *JsonConfig  `json:"json" xml:"json" yaml:"json"`
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

func loadConfigFromEnv() {
	cfgMap := map[string]any{}

	env := os.Environ()
	for _, variable := range env {
		kv := strings.SplitN(variable, "=", 2)
		if strings.HasPrefix(kv[0], "logy.") {
			key := strings.TrimSpace(kv[0])
			key = key[5:]
			cfgMap[key] = kv[1]
		}
	}

	flattenMap := flatMap(cfgMap)
	data, _ := json.Marshal(flattenMap)

	config := &Config{}
	err := json.Unmarshal(data, config)
	if err == nil {
		err = LoadConfig(config)
	}
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
		if config.File != nil && config.File.Enable {
			config.Handlers = append(config.Handlers, "file")
		}
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
			Enable: true,
			Target: TargetStderr,
			Format: defaultLogFormat,
			Color:  true,
			Level:  LevelDebug,
			Json:   nil,
		}
	} else {
		if config.Console.Level == 0 {
			config.Console.Level = LevelDebug
		}

		if config.Console.Enable && strings.TrimSpace(config.Console.Format) == "" {
			return errors.New("console.format cannot be empty or blank")
		}

		if enableConsole {
			config.Console.Enable = true
		}

		if strings.TrimSpace(config.Console.Format) == "" {
			config.Console.Format = defaultLogFormat
		}
	}

	return nil
}

func initializeFileConfig(config *Config) error {

	if config.File == nil {
		config.File = &FileConfig{
			Name:   "logy.log",
			Enable: false,
			Path:   ".",
			Format: defaultLogFormat,
			Level:  LevelInfo,
			Json:   nil,
		}
	} else {
		if config.File.Level == 0 {
			config.File.Level = LevelInfo
		}

		if strings.TrimSpace(config.File.Format) == "" {
			config.File.Format = defaultLogFormat
		}

		if config.File.Level == 0 {
			config.File.Level = LevelInfo
		}

		if config.File.Name == "" {
			config.File.Name = "logy.log"
		}

		if config.File.Path == "" {
			config.File.Path = "."
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

		if name == "file" {
			console, ok := handler.(*FileHandler)

			if !ok {
				continue
			}

			console.onConfigure(config.File)
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
