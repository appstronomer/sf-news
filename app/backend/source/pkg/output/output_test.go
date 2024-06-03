package output

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestOutput_Close(t *testing.T) {
	wantRes := false

	output := Make(&MockWriter{}, &MockWriter{})
	output.Close()

	gotRes := output.Log("hello")
	if gotRes != wantRes {
		t.Errorf("Log() after Close() result got %t; want %t", gotRes, wantRes)
	}

	gotRes = output.Err("hello")
	if gotRes != wantRes {
		t.Errorf("Err() after Close() result got %t; want %t", gotRes, wantRes)
	}
}

func TestOutput_Log(t *testing.T) {
	sendMessage := "hello"
	wantMessage := fmt.Sprintln(strings.Clone(sendMessage))
	wantLenLog := len(wantMessage)
	wantLenErr := 0

	writerLog, writerErr := &MockWriter{}, &MockWriter{}
	output := Make(writerLog, writerErr)
	output.Log(sendMessage)

	helpCheckLengths(t, wantLenLog, wantLenErr, writerLog, writerErr)

	gotMessage := string(writerLog.sl)
	if gotMessage != wantMessage {
		t.Fatalf("message got %#v; want %#v", gotMessage, wantMessage)
	}
}

func TestOutput_Err(t *testing.T) {
	sendMessage := "hello"
	wantMessage := fmt.Sprintln(strings.Clone(sendMessage))
	wantLenLog := 0
	wantLenErr := len(wantMessage)

	writerLog, writerErr := &MockWriter{}, &MockWriter{}
	output := Make(writerLog, writerErr)
	output.Err(sendMessage)

	helpCheckLengths(t, wantLenLog, wantLenErr, writerLog, writerErr)

	gotMessage := string(writerErr.sl)
	if gotMessage != wantMessage {
		t.Fatalf("message got %#v; want %#v", gotMessage, wantMessage)
	}
}

func helpCheckLengths(t *testing.T, wantLenLog, wantLenErr int, writerLog, writerErr *MockWriter) {
	if len(writerLog.sl) != wantLenLog {
		t.Errorf("Log() message length got %d; want %d", len(writerLog.sl), wantLenLog)
	}
	if len(writerErr.sl) != wantLenErr {
		t.Errorf("Err() message length got %d; want %d", len(writerErr.sl), wantLenErr)
	}
}

type MockWriter struct {
	sl       []byte
	isClosed bool
}

func (w *MockWriter) Write(p []byte) (n int, err error) {
	if w.isClosed {
		return 0, io.ErrClosedPipe
	}
	w.sl = make([]byte, len(p))
	return copy(w.sl, p), nil
}

func (w *MockWriter) Close() error {
	w.isClosed = true
	return nil
}
