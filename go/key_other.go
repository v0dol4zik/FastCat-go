//go:build !linux && !darwin

package main

// readKey on unsupported platforms falls back to a whole-line read.
func readKey() string { return readLine() }
