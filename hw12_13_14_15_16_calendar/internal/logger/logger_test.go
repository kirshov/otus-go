package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	msg := "test log"
	logger := New("WARNING")
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	logger.Warning(msg)
	require.Regexp(t, "WARNING "+msg, buf.String())

	buf.Reset()
	logger.Error(msg)
	require.Regexp(t, "ERROR "+msg, buf.String())

	buf.Reset()
	logger.Debug(msg)
	require.Empty(t, buf.String())
}
