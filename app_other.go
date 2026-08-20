//go:build !windows && !linux && !darwin

package main

import (
	"log"
	"runtime/debug"
)

// initPlatformSettings is a no-op fallback for non-Windows, non-Linux platforms.
func (a *App) initPlatformSettings() {
	// No-op
}

// sendNativeWheelZoom is a no-op fallback on non-Windows platforms.
func (a *App) sendNativeWheelZoom(delta int32) {
	log.Printf("[Gotion Go] sendNativeWheelZoom called with delta %d (non-Windows platform)", delta)
}

// trimProcessMemory triggers Go runtime GC & returns OS memory on non-Windows platforms.
func (a *App) trimProcessMemory() {
	debug.FreeOSMemory()
}
