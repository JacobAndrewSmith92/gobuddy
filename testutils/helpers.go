package testutils

import (
	"bytes"
	"os"
)

// Capture captures the standard output of a function for testing purposes.
func Capture(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outputChan <- buf.String()
	}()

	f()
	_ = w.Close()
	return <-outputChan
}
