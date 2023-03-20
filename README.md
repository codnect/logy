![Logy logo](https://user-images.githubusercontent.com/5354910/224152840-2c913efa-f7c3-41ea-b0cc-a215a7ec02cf.png)

# Logy

The Logy is a fast, configurable, and easy-to-use logger for Go applications. It supports various logging levels,
handlers, and hierarchically named loggers.

## Getting Started

### Installation

To use Logy in your Go project, you need to first install it using the following command:

```bash
go get -u github.com/procyon-projects/logy
```

After installing Logy, you can import it in your Go code like this:

```go
import "github.com/procyon-projects/logy"
```

Once you've imported Logy, you can start using it to log messages.

### Usage

Here's an example of how to use Logy:

```go
package main

import (
    "context"
    "github.com/procyon-projects/logy"
)

func main() {
    // logy.Get() creates a logger with the name of the package it is called from
    log := logy.Get()

    // Logging messages with different log levels
    log.Info("This is an information message")
    log.Warn("This is a warning message")
    log.Error("This is an error message")
    log.Debug("This is a debug message")
    log.Trace("This is a trace message")
}
```

The above code produces the following result.

![Output](https://user-images.githubusercontent.com/5354910/226194630-778278b0-80a5-48bd-81f7-e22e4caa96db.png)

If you want to add contextual fields to your log messages, you can use the `logy.WithValue()` function to create a new
context with the desired fields.
This function returns a new context with the additional field(s) and copies any existing contextual fields from the
original context.

You can then use the new context to log messages with contextual fields using the `I()`, `W()`, `E()`, `D()`, and `T()`
methods.
These methods accept the context as the first argument, followed by the log message itself.

```go
package main

import (
	"context"
	"github.com/procyon-projects/logy"
)

func main() {
    // logy.Get() creates a logger with the name of the package it is called from
    log := logy.Get()

    // Change console log format
    err := logy.LoadConfig(&logy.Config{
        Console: &logy.ConsoleConfig{
            Enabled: true,
            Format:  "%d{2006-01-02 15:04:05.000} %l [%x{traceId},%x{spanId}] %p %c : %s%e%n",
            Color:   true,
        },
    })
	
    if err != nil {
        panic(err)
    }
	
    // logy.WithValue() returns a new context with the given field and copies any 
    // existing contextual fields if they exist.
    // This ensures that the original context is not modified and avoids any potential 
    // issues.
    ctx := logy.WithValue(context.Background(), "traceId", "anyTraceId")
    // It will create a new context with the spanId and copies the existing fields
    ctx = logy.WithValue(ctx, "spanId", "anySpanId")
	
    // Logging messages with contextual fields
    log.I(ctx, "info message")
    log.W(ctx, "warning message")
    log.E(ctx, "error message")
    log.D(ctx, "debug message")
    log.T(ctx, "trace message")
}
```

The above code produces the following result.

![Output](https://user-images.githubusercontent.com/5354910/226194821-bd4e7211-b829-4927-a230-eca72701957a.png)

if you want to add a contextual field to an existing context, you can use the `logy.PutValue()` function
to directly modify the context. However, note that this should not be done across multiple goroutines as it is not safe.

```go
// It will put the field into original context, so the original context is changed.
logy.PutValue(ctx, "traceId", "anotherTraceId")
```

### Parameterized Logging

Logy provides support for parameterized log messages.

```go
package main

import (
    "context"
    "github.com/procyon-projects/logy"
)

func main() {
    // logy.Get() creates a logger with the name of the package it is called from
    log := logy.Get()
	
    value := 30
    // Logging a parameterized message
    log.Info("The value {} should be between {} and {}", value, 128, 256)
}
```

The output of the above code execution looks as follows:

```bash
2023-03-19 21:13:06.029186  INFO github.com/procyon-projects/logy/test    : The value 30 should be between 128 and 256
```

### JSON Logging Format

Here is an example of how to enable the JSON formatting for console handler.

```go
package main

import (
    "context"
    "github.com/procyon-projects/logy"
)

func main() {
    // logy.Get() creates a logger with the name of the package it is called from
    log := logy.Get()
	
    // Enabled the Json logging
    err := logy.LoadConfig(&logy.Config{
        Console: &logy.ConsoleConfig{
            Enabled: true,
            Json: &logy.JsonConfig{
                Enabled: true,
            },
        },
    })

    if err != nil {
       panic(err)
    }

    // logy.WithValue() returns a new context with the given field and copies any
    // existing contextual fields if they exist.
    ctx := logy.WithValue(context.Background(), "traceId", "anyTraceId")
    // It will create a new context with the spanId and copies the existing fields
    ctx = logy.WithValue(ctx, "spanId", "anySpanId")

    // Logging an information message
    log.Info("This is an information message")

    // Logging an information message with contextual fields
    log.I(ctx, "info message")
}
```

The output of the above code execution looks as follows:

```bash
{"timestamp":"2023-03-20T20:59:02+03:00","level":"INFO","logger":"github.com/procyon-projects/logy/test","message":"This is an information message"}
{"timestamp":"2023-03-20T20:59:02+03:00","level":"INFO","logger":"github.com/procyon-projects/logy/test","message":"info message","mappedContext":{"traceId":"anyTraceId","spanId":"anySpanId"}}
```

### Error and Stack Trace Logging
If you pass an error to the logging methods as the last argument, it will print the full stack trace along with the given error.

Here is an example:

```go
package main

import (
    "context"
    "errors"
    "github.com/procyon-projects/logy"
)

func main() {
    // logy.Get() creates a logger with the name of the package it is called from
    log := logy.Get()
	
    value := "anyValue"
    err := errors.New("an error occurred")

    // Note that there must not be placeholder(curly braces) for the error. 
    // Otherwise, the only error string will be printed.
    log.Info("The value {} was not inserted", value, err)
    log.Warn("The value {} was not inserted", value, err)
    log.Error("The value {} was not inserted", value, err)
    log.Debug("The value {} was not inserted", value, err)
    log.Error("The value {} was not inserted", value, err)
}
```

The output of the above code execution looks as follows:

```bash
2023-03-20 21:17:03.165347  INFO github.com/procyon-projects/logy/test    : The value anyValue was not inserted
Error: an error occurred
main.main()
    /Users/burakkoken/GolandProjects/procyon-projects/logy/test/main.go:19
2023-03-20 21:17:03.165428  WARN github.com/procyon-projects/logy/test    : The value anyValue was not inserted
Error: an error occurred
main.main()
    /Users/burakkoken/GolandProjects/procyon-projects/logy/test/main.go:20
2023-03-20 21:17:03.165434 ERROR github.com/procyon-projects/logy/test    : The value anyValue was not inserted
Error: an error occurred
main.main()
    /Users/burakkoken/GolandProjects/procyon-projects/logy/test/main.go:21
2023-03-20 21:17:03.165438 DEBUG github.com/procyon-projects/logy/test    : The value anyValue was not inserted
Error: an error occurred
main.main()
    /Users/burakkoken/GolandProjects/procyon-projects/logy/test/main.go:22
2023-03-20 21:17:03.165441 ERROR github.com/procyon-projects/logy/test    : The value anyValue was not inserted
Error: an error occurred
main.main()
    /Users/burakkoken/GolandProjects/procyon-projects/logy/test/main.go:23
```

### Loggers

A Logger instance is used to log messages for an application. Loggers are named,
using a hierarchical dot and slash separated namespace.

For example, the logger named `github.com/procyon-projects` is a parent of the logger
named `github.com/procyon-projects/logy`.
Similarly, `net` is a parent of `net/http` and an ancestor of `net/http/cookiejar`

Logger names can be arbitrary strings, however it's recommended that they are based on the package name or struct name
of the logged component.

Logging messages will be forwarded to handlers attached to the loggers.

### Creating Logger

Logy provides multiple ways of creating a logger.
You can either create a new logger with the `logy.Get()` function,
which creates a named logger with the name of the package it is called from:

```go
log := logy.Get()
```

For example, a logger created in the `github.com/procyon-projects/logy` package would have the
name `github.com/procyon-projects/logy`.

Alternatively, you can use the `logy.Named()` function to create a named logger with a specific name:

```go
log := logy.Named("myLogger")
```

This will create a logger with the given name.

You can also use the `logy.Of()` method to create a logger for a specific type:

```go
log := logy.Of[http.Client]
```

This will create a logger with the name `net/http.Client`.

Invoking the `logy.Named()` function with the same name or the `logy.Get()`function in the same package will return
always the exact same Logger.

```go
// x and y loggers will be the same
x = logy.Named("foo")
y = logy.Named("foo")

// z and q loggers will be the same
z = logy.Get()
q = logy.Get()
```

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

If you want a custom log handler, you can create your own log handler by implementing `logy.Handler` interface.
After completing its implementation, you must register it by using `logy.RegisterHandler` function.
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

### Logging Package

Logging is done on a per-package basis. Each package can be independently configured.
A configuration which applies to a package will also apply to all sub-categories of that package,
unless there is a more specific matching sub-package configuration.

For every package the same settings that are configured on ( console / file / syslog ) apply.

| Property                                            |                         Description                         |           Type | Default |
|:----------------------------------------------------|:-----------------------------------------------------------:|---------------:|:-------:|
| `logy.package`.`package-path`.`level`               |               The log level for this package                |           bool | `TRACE` |
| `logy.package`.`package-path`.`use-parent-handlers` | Specify whether this logger should user its parent handlers |           bool | `true`  |
| `logy.package`.`package-path`.`handlers`            |      The names of the handlers to link to this package      | list of string |         |

### Root Logger Configuration

The root logger is handled separately, and is configured via the following properties:

| Property        |                Description                |           Type |  Default  |
|:----------------|:-----------------------------------------:|---------------:|:---------:|
| `logy.level`    |    The log level for every log package    |           bool |  `TRACE`  |
| `logy.handlers` | The names of handlers to link to the root | list of string | [console] |

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

| Property                                              |                                         Description                                         |                                                     Type |                   Default                    |
|:------------------------------------------------------|:-------------------------------------------------------------------------------------------:|---------------------------------------------------------:|:--------------------------------------------:|
| `logy.console.enabled`                                |                                 Enable the console logging                                  |                                                     bool |                    `true`                    |
| `logy.console.target`                                 |                              Override keys with custom values                               |                    Target(`stdout`, `stderr`, `discard`) |                   `stdout`                   |
| `logy.console.format`                                 | The console log format. Note that this value will be ignored if json is enabled for console |                                                   string | `d{2006-01-02 15:04:05.000000} %p %c : %m%n` |
| `logy.console.color`                                  |                Enable color coded output if the target terminal supports it                 |                                                     bool |                    `true`                    |
| `logy.console.level`                                  |                                    The console log level                                    | Level(`OFF`,`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`,`ALL`) |                   `TRACE`                    |
| `logy.console.json.enabled`                           |                             Enable the JSON console formatting                              |                                                     bool |                   `false`                    |
| `logy.console.json.excluded-keys`                     |                          Keys to be excluded from the Json output                           |                                           list of string |                                              |
| `logy.console.json.key-overrides`.`property-name`     |                              Override keys with custom values                               |                                        map[string]string |                                              |
| `logy.console.json.additional-fields`.`property-name` |                                   Additional field values                                   |                                           map[string]any |                                              |

### File Handler Properties

You can configure the file handler with the following configuration properties:

| Property                                           |                                          Description                                           |                                                     Type |                   Default                    |
|:---------------------------------------------------|:----------------------------------------------------------------------------------------------:|---------------------------------------------------------:|:--------------------------------------------:|
| `logy.file.enabled`                                |                                    Enable the file logging                                     |                                                     bool |                   `false`                    |
| `logy.file.format`                                 |     The file log format. Note that this value will be ignored if json is enabled for file      |                                                   string | `d{2006-01-02 15:04:05.000000} %p %c : %m%n` |
| `logy.file.name`                                   |                       The name of the file in which logs will be written                       |                                                   string |                  `logy.log`                  |
| `logy.file.path`                                   |                       The path of the file in which logs will be written                       |                                                   string |              Working directory               |
| `logy.file.level`                                  |                         The level of logs to be written into the file                          | Level(`OFF`,`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`,`ALL`) |                   `TRACE`                    |
| `logy.file.json.enabled`                           |                                Enable the JSON file formatting                                 |                                                     bool |                   `false`                    |
| `logy.file.json.excluded-keys`                     |                            Keys to be excluded from the Json output                            |                                           list of string |                                              |
| `logy.file.json.key-overrides`.`property-name`     |                                Override keys with custom values                                |                                        map[string]string |                                              |
| `logy.file.json.additional-fields`.`property-name` |                                    Additional field values                                     |                                           map[string]any |                                              |

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
| `logy.syslog.level`              |                        The level of the logs to be logged by syslog logger                        |                                                                                                                                                                                                                                                                      Level(`OFF`,`ERROR`,`WARN`,`INFO`,`DEBUG`,`TRACE`,`ALL`) |                   `TRACE`                    |

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
    level: DEBUG
    json:
      enabled: true
      excluded-keys:
        - level
        - logger
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
      excluded-keys:
        - level
        - logger
      key-overrides:
        timestamp: "@timestamp"
      additional-fields:
        application-name: "my-logy-app"
```

Note that file log will only contain `INFO` or higher order logs because we set the root logger level to `INFO`.

## Performance

Here is the benchmark results.

**Log a message without context fields:**

| Package                 |    Time     | Objects Allocated |
|:------------------------|:-----------:|:-----------------:|
| :star: logy             | 62.04 ns/op |    0 allocs/op    |
| :star: logy(formatting) | 1287 ns/op  |    7 allocs/op    |
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

**Log a message with a logger that already has 10 fields of context:**

| Package                 |     Time     | Objects Allocated |
|:------------------------|:------------:|:-----------------:|
| :star: logy             | 85.29 ns/op  |    0 allocs/op    |
| :star: logy(formatting) | 1369.0 ns/op |    7 allocs/op    |
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