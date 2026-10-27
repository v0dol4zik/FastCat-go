//go:build linux

package main

import "syscall"

const (
	ioctlRead  = syscall.TCGETS
	ioctlWrite = syscall.TCSETS
	vMin       = syscall.VMIN  // 6
	vTime      = syscall.VTIME // 5
)
