//go:build darwin

package main

const (
	ioctlRead  = 0x40487413 // TIOCGETA
	ioctlWrite = 0x80487414 // TIOCSETA
	vMin       = 16         // VMIN
	vTime      = 17         // VTIME
)
