package logy

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"
	"reflect"
	"testing"
)

type testConfigFileReader struct {
	mock.Mock
}

func (r *testConfigFileReader) ReadFile(name string) (data []byte, err error) {
	var (
		ok bool
	)

	args := r.Called(name)
	if len(args) == 2 {
		data, ok = args[0].([]byte)
		if !ok {
			data = nil
		}

		err, ok = args[1].(error)
		if !ok {
			err = nil
		}
	}

	return
}

var (
	mockConfigReader = &testConfigFileReader{}
)

func init() {
	configReader = mockConfigReader
}

func TestLoadConfig(t *testing.T) {
	LoadConfig(&Config{})
}

func TestLoadConfigFromYaml(t *testing.T) {
	testConfig := &Config{
		Level:         LevelAll,
		IncludeCaller: false,
		Handlers:      Handlers{"testHandler", "console", "file"},
		Console: &ConsoleConfig{
			Enabled: true,
			Target:  TargetStderr,
			Format:  "%d %s%e%n",
			Color:   true,
			Level:   LevelAll,
			Json: &JsonConfig{
				Enabled: true,
				ExcludedKeys: ExcludedKeys{
					"timestamp", "anyKey",
				},
				KeyOverrides: KeyOverrides{
					"timestamp": "@timestamp",
				},
				AdditionalFields: AdditionalFields{
					"application": "any_application",
				},
			},
		},
		File: &FileConfig{
			Name:    "anyFileName",
			Enabled: false,
			Path:    "anyFilePath",
			Format:  "%d %s%e%n",
			Level:   LevelWarn,
			Json: &JsonConfig{
				Enabled: true,
				ExcludedKeys: ExcludedKeys{
					"timestamp", "anyKey",
				},
				KeyOverrides: KeyOverrides{
					"timestamp": "@timestamp",
				},
				AdditionalFields: AdditionalFields{
					"application": "any_application",
				},
			},
		},
		Syslog: &SyslogConfig{
			Enabled:          false,
			Endpoint:         "anyEndpoint",
			AppName:          "anyAppName",
			Hostname:         "anyHostname",
			Facility:         FacilityLogAlert,
			LogType:          RFC5424,
			Protocol:         ProtocolUDP,
			Format:           "%d %s%e%n",
			Level:            LevelInfo,
			BlockOnReconnect: false,
		},
		Package: map[string]*PackageConfig{
			"anyPackage": {
				Level:             LevelAll,
				UseParentHandlers: true,
				Handlers:          Handlers{"file"},
			},
		},
		ExternalHandlers: map[string]ConfigProperties{
			"external_handler": {
				"anyPropertyKey": "anyPropertyValue",
			},
		},
	}

	configData, _ := yaml.Marshal(testConfig)
	testConfigMap := make(map[string]interface{})

	yaml.Unmarshal(configData, &testConfigMap)

	logyMap := make(map[string]interface{})
	testConfigMap["logy"] = logyMap
	logyMap["external_handler"] = ConfigProperties{
		"anyPropertyKey": "anyPropertyValue",
	}

	for key, value := range testConfigMap {
		if key == "logy" {
			continue
		}

		if logyField, ok := testConfigMap["logy"]; ok {
			subMap := logyField.(map[string]interface{})
			subMap[key] = value
		}

		delete(testConfigMap, key)
	}

	configData, _ = yaml.Marshal(testConfigMap)

	mockConfigReader.On("ReadFile", "anyFileName").Return(configData, nil)
	err := LoadConfigFromYaml("anyFileName")

	assert.Nil(t, err)
	same := reflect.DeepEqual(testConfig, config)
	assert.True(t, same)
}
