// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"bufio"
	"io"
	"os"

	"golang.org/x/term"
)

func AwaitEnter() error {
	code := 0
	defer func() {
		if code != 0 {
			os.Exit(code)
		}
	}()

	fd := int(os.Stdin.Fd()) //nolint:gosec // G115
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	r := bufio.NewReader(os.Stdin)
	for {
		b, err := r.ReadByte()
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
