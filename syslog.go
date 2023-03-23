package logy

import (
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

type syslogDiscarder struct {
}

func (d *syslogDiscarder) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (d *syslogDiscarder) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (d *syslogDiscarder) Close() error {
	return nil
}

func (d *syslogDiscarder) LocalAddr() net.Addr {
	return nil
}

func (d *syslogDiscarder) RemoteAddr() net.Addr {
	return nil
}

func (d *syslogDiscarder) SetDeadline(tm time.Time) error {
	return nil
}

func (d *syslogDiscarder) SetReadDeadline(tm time.Time) error {
	return nil
}

func (d *syslogDiscarder) SetWriteDeadline(tm time.Time) error {
	return nil
}

type SyslogHandler struct {
	writer  io.Writer
	enabled atomic.Value
	level   atomic.Value
	format  atomic.Value

	endpoint         atomic.Value
	appName          atomic.Value
	hostname         atomic.Value
	facility         atomic.Value
	syslogType       atomic.Value
	protocol         atomic.Value
	blockOnReconnect atomic.Value
	mu               sync.RWMutex

	underTest atomic.Value
}

func newSysLogHandler(underTest bool) *SyslogHandler {
	handler := &SyslogHandler{}
	handler.initializeHandler(underTest)
	return handler
}

func (h *SyslogHandler) initializeHandler(underTest bool) {
	h.SetEnabled(false)
	h.SetLevel(LevelInfo)
	h.setWriter(newSyslogWriter(string(ProtocolTCP), DefaultSyslogEndpoint, false))

	h.SetApplicationName(os.Args[0])
	h.SetHostname("")
	h.SetFormat(DefaultSyslogFormat)
	h.setEndpoint(DefaultSyslogEndpoint)

	h.SetFacility(FacilityUserLevel)
	h.setLogType(RFC5424)
	h.setProtocol(ProtocolTCP)
	h.setBlockOnReconnect(false)
	h.underTest.Store(underTest)
}

func (h *SyslogHandler) Handle(record Record) error {
	facility := h.Facility()
	priority := int(facility-1)*8 + int(record.Level.syslogLevel())

	buf := newBuffer()
	defer buf.Free()

	encoder := getTextEncoder(buf)

	syslogType := h.syslogType.Load().(SysLogType)
	if syslogType == RFC3164 {
		h.writeRFC3164Header(buf, priority, record)
	} else if syslogType == RFC5424 {
		h.writeRFC5424Header(buf, priority, record)
	}

	format := h.format.Load().(string)
	formatText(encoder, format, record, false, true)

	putTextEncoder(encoder)
	_, err := h.writer.Write(*buf)
	return err
}

func (h *SyslogHandler) SetLevel(level Level) {
	h.level.Store(level)
}

func (h *SyslogHandler) Level() Level {
	return h.level.Load().(Level)
}

func (h *SyslogHandler) SetEnabled(enabled bool) {
	h.enabled.Store(enabled)
}

func (h *SyslogHandler) IsEnabled() bool {
	return h.enabled.Load().(bool)
}

func (h *SyslogHandler) IsLoggable(record Record) bool {
	if !h.IsEnabled() {
		return false
	}

	return record.Level <= h.Level()
}

func (h *SyslogHandler) setWriter(writer io.Writer) {
	h.writer = writer
}

func (h *SyslogHandler) Writer() io.Writer {
	return h.writer
}

func (h *SyslogHandler) SetFormat(format string) {
	h.format.Store(format)
}

func (h *SyslogHandler) Format() string {
	return h.format.Load().(string)
}

func (h *SyslogHandler) setEndpoint(endpoint string) {
	h.endpoint.Store(endpoint)
}

func (h *SyslogHandler) Endpoint() string {
	return h.endpoint.Load().(string)
}

func (h *SyslogHandler) SetApplicationName(name string) {
	h.appName.Store(name)
}

func (h *SyslogHandler) ApplicationName() string {
	return h.appName.Load().(string)
}

func (h *SyslogHandler) SetHostname(hostname string) {
	h.hostname.Store(hostname)
}

func (h *SyslogHandler) Hostname() string {
	return h.hostname.Load().(string)
}

func (h *SyslogHandler) SetFacility(facility Facility) {
	h.facility.Store(facility)
}

func (h *SyslogHandler) Facility() Facility {
	return h.facility.Load().(Facility)
}

func (h *SyslogHandler) setLogType(logType SysLogType) {
	h.syslogType.Store(logType)
}

func (h *SyslogHandler) LogType() SysLogType {
	return h.syslogType.Load().(SysLogType)
}

func (h *SyslogHandler) setProtocol(protocol Protocol) {
	h.protocol.Store(protocol)
}

func (h *SyslogHandler) Protocol() Protocol {
	return h.protocol.Load().(Protocol)
}

func (h *SyslogHandler) setBlockOnReconnect(blockOnReconnect bool) {
	h.blockOnReconnect.Store(blockOnReconnect)
}

func (h *SyslogHandler) IsBlockOnReconnect() bool {
	return h.blockOnReconnect.Load().(bool)
}

func (h *SyslogHandler) OnConfigure(config Config) error {
	h.SetEnabled(config.Syslog.Enabled)
	h.SetLevel(config.Syslog.Level)
	h.SetFormat(config.Syslog.Format)

	h.SetApplicationName(config.Syslog.AppName)
	h.SetHostname(config.Syslog.Hostname)
	h.setProtocol(config.Syslog.Protocol)
	h.setLogType(config.Syslog.LogType)
	h.SetFacility(config.Syslog.Facility)
	h.setBlockOnReconnect(config.Syslog.BlockOnReconnect)
	h.setEndpoint(config.Syslog.Endpoint)

	network := h.Protocol()
	address := h.Endpoint()

	underTest := h.underTest.Load().(bool)

	sysWriter := h.writer.(*syslogWriter)

	defer sysWriter.mu.Unlock()
	sysWriter.mu.Lock()

	if !underTest {
		sysWriter.network = string(network)
		sysWriter.address = address
		sysWriter.retry = !h.IsBlockOnReconnect()

		return sysWriter.connect()
	} else {
		sysWriter.writer = &syslogDiscarder{}
	}

	return nil
}

func (h *SyslogHandler) writeRFC5424Header(buf *buffer, priority int, record Record) {
	hostname := h.Hostname()
	appName := h.ApplicationName()

	buf.WriteByte('<')
	buf.WriteInt(int64(priority))
	buf.WriteString(">1 ")

	buf.WriteTimeLayout(record.Time, time.RFC3339)
	buf.WriteByte(' ')

	if hostname == "" {
		buf.WriteByte('-')
	} else {
		buf.WriteString(hostname)
	}

	buf.WriteByte(' ')

	if appName == "" {
		buf.WriteByte('-')
	} else {
		buf.WriteString(appName)
	}

	buf.WriteByte(' ')
	buf.WriteInt(int64(os.Getpid()))
	buf.WriteByte(' ')

	appendLoggerAsText(buf, record.LoggerName, false, true)

	// structured data
	buf.WriteString(" - ")

	*buf = utf8.AppendRune(*buf, 0xFEFF)
}

func (h *SyslogHandler) writeRFC3164Header(buf *buffer, priority int, record Record) {
	hostname := h.Hostname()
	appName := h.ApplicationName()

	buf.WriteByte('<')
	buf.WriteInt(int64(priority))
	buf.WriteByte('>')

	buf.WriteTimeLayout(record.Time, time.Stamp)
	buf.WriteByte(' ')

	if hostname == "" {
		buf.WriteString("UNKNOWN_HOSTNAME")
	} else {
		buf.WriteString(hostname)
	}

	buf.WriteByte(' ')

	if appName == "" {
		buf.WriteString(appName)
	}

	buf.WriteByte('[')
	buf.WriteInt(int64(os.Getpid()))
	buf.WriteString("]: ")
}
