package common

import (
	"io"
	"os"
	"unsafe"

	"golang.org/x/term"
)

func AwaitEnter() error {
	code := 0
	defer func() {
		if code != 0 {
			os.Exit(code)
		}
	}()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState) //nolint:errcheck

	var b byte
	for {
		_, err = os.Stdin.Read(unsafe.Slice(&b, 1))
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		switch b {
		case '\r', '\n', '\x04' /* Ctrl+D */ :
			return nil
		case '\x03' /* Ctrl+C */, '\x1a' /*Ctrl+Z */ :
			os.Stderr.WriteString("signal: interrupt\r\n") //nolint:errcheck
			code = 1
			return nil
		}
	}
}
