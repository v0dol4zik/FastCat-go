//go:build linux || darwin

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// readKey waits for a single keypress (like bash's `read -r -n1`).
// Returns named keys for special input: "up", "down", "left", "right",
// "enter", "esc"; otherwise a single-character string; "" on EOF.
// Falls back to a whole-line read when stdin is not a terminal.
func readKey() string {
	fd := int(os.Stdin.Fd())
	var old syscall.Termios
	if err := ioctlTermios(fd, ioctlRead, &old); err != nil {
		return readLine() // not a terminal (pipe) or unsupported
	}
	raw := old
	raw.Lflag &^= syscall.ICANON | syscall.ECHO
	raw.Cc[vMin], raw.Cc[vTime] = 1, 0 // blocking single-byte reads
	if err := ioctlTermios(fd, ioctlWrite, &raw); err != nil {
		return readLine()
	}
	defer ioctlTermios(fd, ioctlWrite, &old)

	read1 := func() (byte, bool) {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return 0, false
		}
		return buf[0], true
	}

	c, ok := read1()
	if !ok {
		return ""
	}
	if c != 0x1b {
		if c == '\r' || c == '\n' {
			return "enter"
		}
		return string(c)
	}

	// Escape sequence: switch to 100 ms non-blocking reads so a bare
	// ESC press does not hang waiting for the rest of a sequence.
	raw.Cc[vMin], raw.Cc[vTime] = 0, 1
	ioctlTermios(fd, ioctlWrite, &raw)
	b1, ok1 := read1()
	if !ok1 {
		return "esc"
	}
	b2, ok2 := read1()
	if !ok2 {
		return "esc"
	}
	if b1 == '[' {
		switch b2 {
		case 'A':
			return "up"
		case 'B':
			return "down"
		case 'C':
			return "right"
		case 'D':
			return "left"
		}
	}
	return "esc"
}

func ioctlTermios(fd int, req uintptr, t *syscall.Termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), req, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}
