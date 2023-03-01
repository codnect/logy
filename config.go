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
	config   *Config
	configMu sync.RWMutex
)

type ConfigProperties map[string]any

type Target string

const (
	TargetStderr  Target = "stderr"
	TargetStdout  Target = "stdout"
	TargetDiscard Target = "discard"
)

type additionalFields struct {
	encoder *jsonEncoder
}

type JsonAdditionalFields map[string]any

type JsonConfig struct {
	ExcludeKeys      []string             `json:"exclude-keys" xml:"exclude-keys" yaml:"exclude-keys"`
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
	Handlers          []string `json:"handlers" xml:"handlers" yaml:"handlers"`
}

type Config struct {
	Level            Level                       `json:"level" xml:"level" yaml:"level"`
	IncludeCaller    bool                        `json:"include-caller" xml:"include-caller" yaml:"include-caller"`
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

	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err == nil {
		err = LoadConfig(cfg)
	}
}

func LoadConfig(cfg *Config) error {
	if cfg == nil {
		return errors.New("config cannot be nil")
	}

	if cfg.Level == 0 {
		cfg.Level = LevelInfo
	}

	enableConsole := false

	if cfg.Handlers == nil || len(cfg.Handlers) == 0 {
		cfg.Handlers = []string{"console"}
		enableConsole = true
		if cfg.File != nil && cfg.File.Enable {
			cfg.Handlers = append(cfg.Handlers, "file")
		}
	}

	err := initializePackageConfig(cfg)
	if err != nil {
		return err
	}

	err = initializeConsoleConfig(cfg, enableConsole)
	if err != nil {
		return err
	}

	err = initializeFileConfig(cfg)
	if err != nil {
		return err
	}

	defer configMu.Unlock()
	configMu.Lock()

	config = cfg

	err = configureHandlers(cfg)
	if err != nil {
		return err
	}

	return configureLoggers()
}

func initializePackageConfig(cfg *Config) error {
	if cfg.Package == nil {
		cfg.Package = map[string]*PackageConfig{}
	}

	for pkg, pkgCfg := range cfg.Package {
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

func initializeConsoleConfig(cfg *Config, enableConsole bool) error {

	if cfg.Console == nil {
		cfg.Console = &ConsoleConfig{
			Enable: true,
			Target: TargetStderr,
			Format: defaultLogFormat,
			Color:  true,
			Level:  LevelDebug,
			Json:   nil,
		}
	} else {
		if cfg.Console.Level == 0 {
			cfg.Console.Level = LevelDebug
		}

		if cfg.Console.Enable && strings.TrimSpace(cfg.Console.Format) == "" {
			return errors.New("console.format cannot be empty or blank")
		}

		if enableConsole {
			cfg.Console.Enable = true
		}

		if strings.TrimSpace(cfg.Console.Format) == "" {
			cfg.Console.Format = defaultLogFormat
		}
	}

	return nil
}

func initializeFileConfig(cfg *Config) error {

	if cfg.File == nil {
		cfg.File = &FileConfig{
			Name:   defaultLogFilePath,
			Enable: false,
			Path:   defaultLogFileName,
			Format: defaultLogFormat,
			Level:  LevelInfo,
			Json:   nil,
		}
	} else {
		if cfg.File.Level == 0 {
			cfg.File.Level = LevelInfo
		}

		if strings.TrimSpace(cfg.File.Format) == "" {
			cfg.File.Format = defaultLogFormat
		}

		if cfg.File.Level == 0 {
			cfg.File.Level = LevelInfo
		}

		if cfg.File.Name == "" {
			cfg.File.Name = defaultLogFileName
		}

		if cfg.File.Path == "" {
			cfg.File.Path = defaultLogFilePath
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
