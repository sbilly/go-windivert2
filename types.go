package windivert

/*
#include <windows.h>
#include <windivert.h>
*/
import "C"

import (
	"sync"
	"unsafe"
)

// Handle represents a WinDivert handle
type Handle struct {
	handle C.HANDLE
	mutex  sync.Mutex
}

// Lock locks the handle
func (h *Handle) Lock() {
	h.mutex.Lock()
}

// Unlock unlocks the handle
func (h *Handle) Unlock() {
	h.mutex.Unlock()
}

// Close closes the WinDivert handle
func (h *Handle) Close() error {
	if h.handle != C.INVALID_HANDLE_VALUE {
		if C.WinDivertClose(h.handle) == 0 {
			return getLastError()
		}
		h.handle = C.INVALID_HANDLE_VALUE
	}
	return nil
}

// Open opens a WinDivert handle
func Open(filter string, layer Layer, priority int16, flags uint64) (*Handle, error) {
	cfilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cfilter))

	handle := C.WinDivertOpen(cfilter, C.WINDIVERT_LAYER(layer), C.INT16(priority), C.UINT64(flags))
	if handle == C.INVALID_HANDLE_VALUE {
		return nil, getLastError()
	}

	return &Handle{handle: handle}, nil
}
