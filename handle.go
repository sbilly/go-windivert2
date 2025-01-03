package windivert

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -L${SRCDIR}/lib -lwindivert
// #include <windivert.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
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

// Recv receives a packet
func (h *Handle) Recv(packet []byte) (int, *Address, error) {
	var addr C.WINDIVERT_ADDRESS
	var readLen C.UINT

	if len(packet) == 0 {
		return 0, nil, fmt.Errorf("packet buffer is empty")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	ret := C.WinDivertRecv(h.handle, unsafe.Pointer(&packet[0]), C.UINT(len(packet)), &addr, &readLen)
	if ret == 0 {
		return 0, nil, getLastError()
	}

	goAddr := &Address{}
	*goAddr = *(*Address)(unsafe.Pointer(&addr))

	return int(readLen), goAddr, nil
}

// Send sends a packet
func (h *Handle) Send(packet []byte, addr *Address) (int, error) {
	var writeLen C.UINT

	if len(packet) == 0 {
		return 0, fmt.Errorf("packet buffer is empty")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	caddr := (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(addr))
	ret := C.WinDivertSend(h.handle, unsafe.Pointer(&packet[0]), C.UINT(len(packet)), caddr, &writeLen)
	if ret == 0 {
		return 0, getLastError()
	}

	return int(writeLen), nil
}
