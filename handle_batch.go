package windivert

import (
	"fmt"
	"unsafe"
)

// RecvEx receives multiple packets
func (h *Handle) RecvEx(packets [][]byte, addrs []Address, flags uint64) (int, int, error) {
	if len(packets) == 0 || len(addrs) == 0 {
		return 0, 0, fmt.Errorf("empty packets or addresses buffer")
	}

	maxPackets := len(packets)
	if len(addrs) < maxPackets {
		maxPackets = len(addrs)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	var readLen C.UINT
	var addrLen C.UINT

	ret := C.WinDivertRecvEx(
		h.handle,
		unsafe.Pointer(&packets[0][0]),
		C.UINT(len(packets[0])*len(packets)),
		&readLen,
		C.UINT64(flags),
		(*C.WINDIVERT_ADDRESS)(unsafe.Pointer(&addrs[0])),
		&addrLen,
		nil,
	)

	if ret == 0 {
		return 0, 0, getLastError()
	}

	return int(readLen), int(addrLen), nil
}

// SendEx sends multiple packets
func (h *Handle) SendEx(packets [][]byte, addrs []Address, flags uint64) (int, error) {
	if len(packets) == 0 || len(addrs) == 0 {
		return 0, fmt.Errorf("empty packets or addresses buffer")
	}

	maxPackets := len(packets)
	if len(addrs) < maxPackets {
		maxPackets = len(addrs)
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	var writeLen C.UINT

	ret := C.WinDivertSendEx(
		h.handle,
		unsafe.Pointer(&packets[0][0]),
		C.UINT(len(packets[0])*len(packets)),
		&writeLen,
		C.UINT64(flags),
		(*C.WINDIVERT_ADDRESS)(unsafe.Pointer(&addrs[0])),
		C.UINT(len(addrs)),
		nil,
	)

	if ret == 0 {
		return 0, getLastError()
	}

	return int(writeLen), nil
}
