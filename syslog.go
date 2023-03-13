package logy

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

type SyslogHandler struct {
	commonHandler
	endpoint         atomic.Value
	appName          atomic.Value
	hostname         atomic.Value
	facility         atomic.Value
	syslogType       atomic.Value
	protocol         atomic.Value
	blockOnReconnect atomic.Value
	mu               sync.RWMutex
}

func newSysLogHandler() *SyslogHandler {
	handler := &SyslogHandler{}
	handler.initializeHandler()

	handler.SetEnabled(false)
	handler.SetLevel(LevelInfo)
	handler.setWriter(&discarder{})

	handler.SetApplicationName(os.Args[0])
	handler.SetHostname("")
	handler.SetFormat(DefaultSyslogFormat)
	handler.setEndpoint(DefaultSyslogEndpoint)

	handler.SetFacility(FacilityUserLevel)
	handler.setLogType(RFC5424)
	handler.setProtocol(ProtocolTCP)
	handler.setBlockOnReconnect(false)
	return handler
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

	sWriter := newSyslogWriter(string(network), address, !h.IsBlockOnReconnect())
	h.setWriter(sWriter)
	return sWriter.connect()
}

func (h *SyslogHandler) reconnect() {
	network := h.Protocol()
	address := h.Endpoint()

	con, err := net.Dial(string(network), address)
	if err != nil {
		h.SetEnabled(false)
		h.setWriter(&discarder{})
		return
	}

	h.setWriter(con)
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
	h.formatText(encoder, format, record, false, true)

	putTextEncoder(encoder)
	_, err := h.writer.Write(*buf)
	return err
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
