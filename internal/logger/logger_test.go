package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var errSome = fmt.Errorf("some error")

func TestError(t *testing.T) {
	result := wrapLogTest(t, func(logger *Logger) {
		logger.Error(errSome)
	})

	require.Contains(t, result, errSome.Error())
}

func wrapLogTest(t *testing.T, callback func(logger *Logger)) string {
	t.Helper()

	rescueStdout := os.Stdout
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = writer
	logger := NewLogger()

	callback(logger)

	writer.Close()
	out, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	os.Stdout = rescueStdout

	return string(out)
}
