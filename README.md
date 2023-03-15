![Logy logo](https://user-images.githubusercontent.com/5354910/224152840-2c913efa-f7c3-41ea-b0cc-a215a7ec02cf.png)

# Logy

The Logy package provides a fast, configurable and simple logger.

## Installation

```bash
go get -u github.com/procyon-projects/logy
```

## Loggers
A Logger instance is used to log messages for an application. Loggers are named, 
using a hierarchical dot and slash separated namespace.

For example, the logger named `github.com/procyon-projects` is a parent of the logger named `github.com/procyon-projects/logy`.
Similarly, `net` is a parent of `net/http` and an ancestor of `net/http/cookiejar`.

Logger names can be arbitrary strings, however we recommend that they are based on the package name or struct name of the logged component.

Logging messages will be forwarded to handlers attached to the loggers.

Loggers are obtained by using one of the following approaches. They will either create a new Logger or return an existing Logger.
Depending on your application code, you can use one of them. 

```go
package test

import (
	"github.com/procyon-projects/logy"
	"net/http"
)

var (
    // The name of the logger will be the name of the package the logy.New() function was called from
    x = logy.New()
    // The name of the logger will be `github.com`. Its parent logger will be `github`
    y = logy.Named("github.com")
    // The name of the logger will be `net/http.Client`.
    z = logy.Of[http.Client]
)
```

Invoking the `logy.Named()` function with the same name or the `logy.New()`function in the same package will return always the exact same Logger.

```go
package test

import (
	"github.com/procyon-projects/logy"
)

var (
    // x and y loggers will be the same
    x = logy.Named("foo")
    y = logy.Named("foo")
    // z and q loggers will be the same
    z = logy.New()
    q = logy.New()
)
```

## Logging Levels

Logy provides many logging levels. Below is the complete list.

* **ERROR**
* **WARN**
* **INFO**
* **DEBUG**
* **TRACE**

## Log Handlers

A log handler is a logging component that sends log messages to a writer. Logy
includes the following log handlers:

### Console Log Handler

The console log handler is enabled by default. It outputs all log messages to the console of your application.
(typically to the system's `stdout`)

### File Log Handler

The file log handler is disabled by default. It outputs all log messages to a file.
`Note that there is no support log file rotation.`

### Syslog Log Handler

The syslog log handler is disabled by default. It send all log messages to a syslog server (by default,
the syslog server runs on the same host as the application)

### External Log Handler
Customized log handlers that implements the `Handler` interface can be used, but you need to register them by calling `logy.RegisterHandler()`
function. After registration, it's ready for receiving log messages and output them.

```go
type Handler interface {
    Handle(record Record) error
    SetLevel(level Level)
    Level() Level
    SetEnabled(enabled bool)
    IsEnabled() bool
    IsLoggable(record Record) bool
    Writer() io.Writer
}
```

Here is an example of how you can register your custom handlers:

```go
// CustomHandler implements the Handler interface
type CustomHandler struct {
	
}

func newCustomHandler() *CustomHandler {
    return &CustomHandler{...}
}
...

func init() {
    logy.RegisterHandler("handlerName", newCustomHandler())
}
```

## Logging Configuration

In order to configure the logging, you can use the following approaches:
* Environment Variables
* Programmatically

You can load the yaml logging configuration files as shown below.
```go
func init() {
    err := logy.LoadConfigFromYaml("logy.config.yaml")
	
    if err != nil {
        panic(err)
    }
}
```

As an alternative, you can configure the logging by invoking `logy.LoadConfig()` function.
```go
func init() {
    err := logy.LoadConfig(&logy.Config{
            Level:    logy.LevelTrace,
            Handlers: logy.Handlers{"console", "file"},
            Console: &logy.ConsoleConfig{
            Level:   logy.LevelTrace,
            Enabled: true,
            // this will be ignored because console json logging is enabled
            Format: "%d{2006-01-02 15:04:05.000} %l [%x{traceId},%x{spanId}] %p : %s%e%n",
            Color:   true,
            Json:    &logy.JsonConfig{
                Enabled: true,
                KeyOverrides: logy.KeyOverrides{
                    "timestamp": "@timestamp",
                },
                AdditionalFields: logy.JsonAdditionalFields{
                    "application-name": "my-logy-app",
                },
            },
        },
        File: &logy.FileConfig{
            Enabled: true,
            Name: "file_trace.log",
            Path: "/var",
            // this will be ignored because file json logging is enabled
            Format: "d{2006-01-02 15:04:05} %p %s%e%n",
            Json:    &logy.JsonConfig{
                Enabled: true,
                KeyOverrides: logy.KeyOverrides{
                    "timestamp": "@timestamp",
                },
                AdditionalFields: logy.JsonAdditionalFields{
                    "application-name": "my-logy-app",
                },
            },
        },
    })

    if err != nil {
        panic(err)
    }
}
```

### Logging Format

By default, Logy uses a pattern-based logging format.

You can customize the format for each log handler using a dedicated configuration property.
For the console handler, the property is `logy.console.format`.

The following table shows the logging format string symbols that you can use to configure the format of the log
messages.

Supported logging format symbols:

| Symbol              |              Summary               |                                                Description                                                 |
|:--------------------|:----------------------------------:|:----------------------------------------------------------------------------------------------------------:|
| `%%`                |                `%`                 |                                           A simple `%%`character                                           |
| `%c`                |            Logger name             |                                              The logger name                                               |
| `%C`                |            Package name            |                                              The package name                                              |
| `%d{layout}`        |                Date                |                                     Date with the given layout string                                      |
| `%e`                |               Error                |                                           The error stack trace                                            |
| `%F`                |            Source file             |                                            The source file name                                            |
| `%i`                |             Process ID             |                                          The current process PID                                           |
| `%l`                |          Source location           |                          The source location(file name, line number, method name)                          |
| `%L`                |            Source line             |                                           The source line number                                           |
| `%m`                |            Full Message            |                                   The log message including error trace                                    |
| `%M`                |           Source method            |                                           The source method name                                           |
| `%n`                |              Newline               |                                         The line separator string                                          |
| `%N`                |            Process name            |                                      The name of the current process                                       |
| `%p`                |               Level                |                                      The logging level of the message                                      |
| `%s`                |           Simple message           |                                    The log message without error trace                                     |
| `%X{property-name}` |        Mapped Context Value        |                        The value from Mapped Context  `property-key=property-value`                        |
| `%x{property-name}` |  Mapped Context Value without key  |                   The value without key from Mapped Context  in format `property-value`                    |
| `%X`                |       Mapped Context Values        | All the values from Mapped Context in format `property-key1=property-value1,property-key2=property-value2` |
| `%x`                | Mapped Context Values without keys |        All the values without keys from Mapped Context in format `property-value1,property-value2`         |


### Console Handler Properties
You can configure the console handler with the following configuration properties:

| Property                                              |                                         Description                                         |                                         Type |                   Default                    |
|:------------------------------------------------------|:-------------------------------------------------------------------------------------------:|---------------------------------------------:|:--------------------------------------------:|
| `logy.console.enabled`                                |                                 Enable the console logging                                  |                                         bool |                    `true`                    |
| `logy.console.target`                                 |                              Override keys with custom values                               |        Target(`stdout`, `stderr`, `discard`) |                   `stdout`                   |
| `logy.console.format`                                 | The console log format. Note that this value will be ignored if json is enabled for console |                                       string | `d{2006-01-02 15:04:05.000000} %p %c : %m%n` |
| `logy.console.color`                                  |                Enable color coded output if the target terminal supports it                 |                                         bool |                    `true`                    |
| `logy.console.level`                                  |                                    The console log level                                    | Level(`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`) |                   `TRACE`                    |
| `logy.console.json.enabled`                           |                             Enable the JSON console formatting                              |                                         bool |                   `false`                    |
| `logy.console.json.key-overrides`.`property-name`     |                              Override keys with custom values                               |                            map[string]string |                                              |
| `logy.console.json.additional-fields`.`property-name` |                                   Additional field values                                   |                               map[string]any |                                              |

### File Handler Properties
You can configure the file handler with the following configuration properties:

| Property                                           |                                      Description                                      |                                         Type |                   Default                    |
|:---------------------------------------------------|:-------------------------------------------------------------------------------------:|---------------------------------------------:|:--------------------------------------------:|
| `logy.file.enabled`                                |                                Enable the file logging                                |                                         bool |                   `false`                    |
| `logy.file.format`                                 | The file log format. Note that this value will be ignored if json is enabled for file |                                       string | `d{2006-01-02 15:04:05.000000} %p %c : %m%n` |
| `logy.file.name`                                   |                  The name of the file in which logs will be written                   |                                       string |                  `logy.log`                  |
| `logy.file.path`                                   |                  The path of the file in which logs will be written                   |                                       string |              Working directory               |
| `logy.file.level`                                  |                     The level of logs to be written into the file                     | Level(`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`) |                   `TRACE`                    |
| `logy.file.json.enabled`                           |                            Enable the JSON file formatting                            |                                         bool |                   `false`                    |
| `logy.file.json.key-overrides`.`property-name`     |                           Override keys with custom values                            |                            map[string]string |                                              |
| `logy.file.json.additional-fields`.`property-name` |                                Additional field values                                |                               map[string]any |                                              |

### Syslog Handler Properties
You can configure the syslog handler with the following configuration properties:

| Property                         |                                            Description                                            |                                                                                                                                                                                                                                                                                                                          Type |                   Default                    |
|:---------------------------------|:-------------------------------------------------------------------------------------------------:|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:--------------------------------------------:|
| `logy.syslog.enabled`            |                                     Enable the syslog logging                                     |                                                                                                                                                                                                                                                                                                                          bool |                   `false`                    |
| `logy.syslog.endpoint`           |                           The IP address and port of the syslog server                            |                                                                                                                                                                                                                                                                                                                     host:port |               `localhost:514`                |
| `logy.syslog.app-name`           |                 The app name used when formatting the message in `RFC5424` format                 |                                                                                                                                                                                                                                                                                                                        string |                                              |
| `logy.syslog.hostname`           |                       The name of the host the messages are being sent from                       |                                                                                                                                                                                                                                                                                                                        string |                                              |
| `logy.syslog.facility`           | The facility used when calculating the priority of the message in `RFC5424` and  `RFC3164` format | Facility(`kernel`,`user-level`,`mail-system`,`system-daemons`,`security`,`syslogd`,`line-printer`,`network-news`,`uucp`,`clock-daemon`,`security2`,`ftp-daemon`,`ntp`,`log-audit`,`log-alert`,`clock-daemon2`,`local-use-0`,`local-use-1`,`local-use-2`,`local-use-3`,`local-use-4`,`local-use-5`,`local-use-6`,`local-use-7` |                 `user-level`                 |
| `logy.syslog.log-type`           |                     The message format type used when formatting the message                      |                                                                                                                                                                                                                                                                                               SysLogType(`rfc5424`,`rfc3164`) |                  `rfc5424`                   |
| `logy.syslog.protocol`           |                         The protocol used to connect to the syslog server                         |                                                                                                                                                                                                                                                                                                         Protocol(`tcp`,`udp`) |                    `tcp`                     |
| `logy.syslog.block-on-reconnect` |             Enable or disable blocking when attempting to reconnect the syslog server             |                                                                                                                                                                                                                                                                                                                          bool |                   `false`                    |
| `logy.syslog.format`             |                                      The log message format                                       |                                                                                                                                                                                                                                                                                                                        string | `d{2006-01-02 15:04:05.000000} %p %c : %m%n` |
| `logy.syslog.level`              |                        The level of the logs to be logged by syslog logger                        |                                                                                                                                                                                                                                                                                  Level(`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`) |                   `TRACE`                    |


## Examples YAML Logging Configuration

*Console Logging Configuration*

```yaml
logy:
  level: INFO

  console:
    enabled: true
    # Send output to stderr
    target: stderr
    format: "%d{2006-01-02 15:04:05.000} %l [%x{traceId},%x{spanId}] %p : %s%e%n"
    # Disable color coded output
    color: false
    level: DEBUG
```

*Console JSON Logging Configuration*

```yaml
logy:
  level: INFO
  console:
    enabled: true
    # Send output to stderr
    target: stderr
    format: "%d{2006-01-02 15:04:05.000} %l [%x{traceId},%x{spanId}] %p : %s%e%n"
    level: DEBUG
    json: 
      enabled: true
      key-overrides:
        timestamp: "@timestamp"
      additional-fields:
        application-name: "my-logy-app"
```
Note that console log will only contain `INFO` or higher order logs because we set the root logger level to `INFO`.


*File Logging Configuration*

```yaml
logy:
  level: INFO
  file:
    enabled: true
    level: TRACE
    format: "d{2006-01-02 15:04:05} %p %s%e%n"
    # Send output to a file_trace.log under the /var directory
    name: file_trace.log
    path: /var
```

*File JSON Logging Configuration*

```yaml
logy:
  level: INFO
  file:
    enabled: true
    level: TRACE
    # Send output to a file_trace.log under the /var directory
    name: file_trace.log
    path: /var
    json: 
      enabled: true
      key-overrides:
        timestamp: "@timestamp"
      additional-fields:
        application-name: "my-logy-app"
```
Note that file log will only contain `INFO` or higher order logs because we set the root logger level to `INFO`.

## Performance

Log a message without context fields

| Package                 |    Time     | Objects Allocated |
|:------------------------|:-----------:|:-----------------:|
| :star: logy             | 27.99 ns/op |    0 allocs/op    |
| :star: logy(formatting) | 883.8 ns/op |    7 allocs/op    |
| :zap: exp/slog          | 38.08 ns/op |    0 allocs/op    |
| zerolog                 | 37.49 ns/op |    0 allocs/op    |
| zerolog(formatting)     | 3030 ns/op  |   108 allocs/op   |
| zap                     | 98.30 ns/op |    0 allocs/op    |
| zap sugar               | 110.9 ns/op |    1 allocs/op    |
| zap sugar (formatting)  | 3369 ns/op  |   108 allocs/op   |
| go-kit                  | 248.5 ns/op |    9 allocs/op    |
| log15                   | 2490 ns/op  |   20 allocs/op    |
| apex/log                | 1139 ns/op  |    6 allocs/op    |
| logrus                  | 1831 ns/op  |   23 allocs/op    |

Log a message with a logger that already has 10 fields of context:

| Package                 |     Time     | Objects Allocated |
|:------------------------|:------------:|:-----------------:|
| :star: logy             | 61.43 ns/op  |    0 allocs/op    |
| :star: logy(formatting) | 1208.0 ns/op |    7 allocs/op    |
| :zap: exp/slog          | 266.3 ns/op  |    0 allocs/op    |
| zerolog                 | 44.84 ns/op  |    0 allocs/op    |
| zerolog(formatting)     | 3103.0 ns/op |   108 allocs/op   |
| zap                     | 92.50 ns/op  |    0 allocs/op    |
| zap sugar               | 113.7 ns/op  |    1 allocs/op    |
| zap sugar (formatting)  |  3355 ns/op  |   108 allocs/op   |
| go-kit                  |  3628 ns/op  |   66 allocs/op    |
| log15                   | 12532 ns/op  |   130 allocs/op   |
| apex/log                | 14494 ns/op  |   53 allocs/op    |
| logrus                  | 16246 ns/op  |   68 allocs/op    |

# License
Logy is released under MIT License.