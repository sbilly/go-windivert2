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

// Recv receives a single packet
func (h *Handle) Recv(packet []byte, addr *Address) (uint, error) {
	nr, _, err := h.RecvEx([][]byte{packet}, []Address{*addr}, 0)
	if err != nil {
		return 0, err
	}
	return nr, nil
}

// Send sends a single packet
func (h *Handle) Send(packet []byte, addr *Address) (uint, error) {
	nw, err := h.SendEx([][]byte{packet}, []Address{*addr}, 0)
	if err != nil {
		return 0, err
	}
	return nw, nil
}
