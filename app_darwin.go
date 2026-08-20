//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdint.h>
#include <mach/mach.h>
#if __has_include(<mach/task.h>)
#include <mach/task.h>
#endif
#if __has_include(<mach/task_info.h>)
#include <mach/task_info.h>
#endif

// sendCGEventScroll sends a scroll wheel event to simulate Ctrl+Wheel for zoom
void sendCGEventScroll(int32_t deltaY) {
    CGEventRef scrollEvent = CGEventCreateScrollWheelEvent(
        NULL,
        kCGScrollEventUnitLine,
        1,  // number of axes
        deltaY
    );
    if (scrollEvent) {
        CGEventSetFlags(scrollEvent, kCGEventFlagMaskControl);
        CGEventPost(kCGSessionEventTap, scrollEvent);
        CFRelease(scrollEvent);
    }
}

// trimProcessMemory trims memory using mach task APIs
long long getResidentMemory() {
    struct mach_task_basic_info info;
    mach_msg_type_number_t count = MACH_TASK_BASIC_INFO_COUNT;
    kern_return_t kr = task_info(mach_task_self(), MACH_TASK_BASIC_INFO, (task_info_t)&info, &count);
    if (kr == KERN_SUCCESS) {
        return (long long)info.resident_size;
    }
    return 0;
}
*/
import "C"

import (
	"log"
	"os"
	"runtime/debug"
)

// initPlatformSettings initializes macOS-specific WKWebView settings
func (a *App) initPlatformSettings() {
	// On macOS, WKWebView manages cookies through its own API
	// Wails handles this automatically via the WKWebView configuration
	// We just need to ensure the profile directory exists
	profileDir := a.userDataDir
	if profileDir != "" {
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			log.Printf("[Gotion macOS] Warning: Failed to create profile directory: %v", err)
		} else {
			log.Printf("[Gotion macOS] Profile directory initialized: %s", profileDir)
		}
	}

	// Log memory info before trimming
	memBefore := C.getResidentMemory()
	log.Printf("[Gotion macOS] Initial resident memory: %d bytes", int64(memBefore))
}

// sendNativeWheelZoom sends a scroll wheel event to simulate zoom via CGEvent
func (a *App) sendNativeWheelZoom(delta int32) {
	log.Printf("[Gotion macOS] Sending native scroll wheel event with delta %d", delta)

	// Send scroll event via CGEvent API
	C.sendCGEventScroll(C.int(delta))

	// Small delay to let the event propagate
	// (CGEvent is synchronous, but we log for debugging)
	log.Printf("[Gotion macOS] Scroll wheel event sent successfully")
}

// trimProcessMemory triggers Go runtime GC and reports memory usage on macOS
func (a *App) trimProcessMemory() {
	// Get memory before GC
	memBefore := C.getResidentMemory()

	// Trigger Go GC
	debug.FreeOSMemory()

	// Get memory after GC
	memAfter := C.getResidentMemory()

	// Log memory change
	if int64(memBefore) > 0 && int64(memAfter) > 0 {
		saved := int64(memBefore) - int64(memAfter)
		if saved > 0 {
			log.Printf("[Gotion macOS] Memory trimmed: %d -> %d bytes (saved %d bytes)", int64(memBefore), int64(memAfter), saved)
		} else {
			log.Printf("[Gotion macOS] Memory after GC: %d bytes (no reduction)", int64(memAfter))
		}
	} else {
		log.Printf("[Gotion macOS] Memory trim completed (freeOSMemory called)")
	}
}

// init is called when the app starts on macOS
func init() {
	// macOS-specific initialization if needed
	// WKWebView is managed by Wails automatically
}
