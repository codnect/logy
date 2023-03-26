package logy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net"
	"runtime"
	"testing"
	"time"
)

var (
	_, writerTestFilename, _, _ = runtime.Caller(0)
)

type mockDiscarder struct {
	mock.Mock
}

func (d *mockDiscarder) Write(b []byte) (n int, err error) {
	var (
		ok bool
	)

	args := d.MethodCalled("Write", b)
	if len(args) == 2 {
		n, ok = args[0].(int)
		if !ok {
			n = 0
		}

		err, ok = args[1].(error)
		if !ok {
			err = nil
		}
	}

	return
}

type mockConn struct {
	mock.Mock
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (c *mockConn) Write(b []byte) (n int, err error) {
	var (
		ok bool
	)

	args := c.MethodCalled("Write", b)
	if len(args) == 2 {
		n, ok = args[0].(int)
		if !ok {
			n = 0
		}

		err, ok = args[1].(error)
		if !ok {
			err = nil
		}
	}

	return
}

func (c *mockConn) Close() error {
	return nil
}

func (c *mockConn) LocalAddr() net.Addr {
	return nil
}

func (c *mockConn) RemoteAddr() net.Addr {
	return nil
}

func (c *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestSyncWriter_WriteShouldNotOutputDataToWriterIfItIsDiscarded(t *testing.T) {
	writer := newSyncWriter(nil, true)
	writer.Write([]byte("anyMessage"))
}

func TestSyncWriter_WriteShouldOutputDataToWriterIfItIsNotDiscarded(t *testing.T) {
	mockDiscarder := &mockDiscarder{}
	mockDiscarder.On("Write", []byte("anyMessage")).Return(0, nil)

	writer := newSyncWriter(mockDiscarder, false)
	writer.Write([]byte("anyMessage"))

	mockDiscarder.AssertCalled(t, "Write", []byte("anyMessage"))
}

func TestSyncWriter_isDiscarded(t *testing.T) {
	writer := newSyncWriter(nil, true)
	assert.True(t, writer.isDiscarded())

	writer.setDiscarded(false)
	assert.False(t, writer.isDiscarded())
}

func TestSyslogWriter_WriteShouldNotOutputDataToWriterIfItIsDiscarded(t *testing.T) {
	writer := newSyslogWriter("", "", false, true)
	writer.Write([]byte("anyMessage"))
}

func TestSyslogWriter_WriteShouldOutputDataToWriterIfItIsNotDiscarded(t *testing.T) {
	mockConn := &mockConn{}
	mockConn.On("Write", []byte("anyMessage")).Return(0, nil)

	writer := newSyslogWriter("", "", false, false)
	writer.writer = mockConn
	writer.Write([]byte("anyMessage"))

	mockConn.AssertCalled(t, "Write", []byte("anyMessage"))
}

func TestSyslogWriter_isDiscarded(t *testing.T) {
	writer := newSyslogWriter("", "", false, true)
	assert.True(t, writer.isDiscarded())

	writer.setDiscarded(false)
	assert.False(t, writer.isDiscarded())
}

func TestGlobalWriter_Write(t *testing.T) {
	current := time.Now()
	now = func() time.Time {
		return current
	}

	logger := &Logger{
		isDefault: true,
	}
	logger.SetLevel(LevelAll)
	logger.includesCaller.Store(false)
	writer := newGlobalWriter(logger)
	log.SetOutput(writer)

	mockTestHandler.On("Handle", mock.AnythingOfType("Record")).Return(nil)

	log.SetFlags(0)
	log.Print("anyMessage")

	mockTestHandler.AssertCalled(t, "Handle", Record{
		Time:       current,
		Level:      LevelDebug,
		Message:    fmt.Sprintf("%s", "anyMessage"),
		Context:    nil,
		LoggerName: "github.com/procyon-projects/logy",
		StackTrace: "",
		Error:      nil,
		Caller: Caller{
			defined:  true,
			file:     writerTestFilename,
			line:     160,
			function: "github.com/procyon-projects/logy.TestGlobalWriter_Write",
		},
	})
}
