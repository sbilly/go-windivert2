package windivert

/*
#include <windows.h>
#include <windivert.h>

HANDLE getLastError() {
    return GetLastError();
}
*/
import "C"

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CtlCode represents a control code
type CtlCode uint32

// IoCtl represents an IO control structure
type IoCtl struct {
	Code   CtlCode
	Pkt    uint64
	Addr   uint64
	Param  uint32
	Length uint32
}

// recv represents a receive operation
type recv struct {
	Addr       uint64
	AddrLenPtr uint64
}

// send represents a send operation
type send struct {
	Addr    uint64
	AddrLen uint64
}

var (
	// WinDivert is the DLL instance
	WinDivert = (*windows.DLL)(nil)
	// WinDivertOpen is the WinDivertOpen procedure
	WinDivertOpen = (*windows.Proc)(nil)
	// WinDivertSys is the path to WinDivert sys file
	WinDivertSys = ""
	// WinDivertDll is the path to WinDivert dll file
	WinDivertDll = ""
	// DeviceName is the WinDivert device name
	DeviceName = windows.StringToUTF16Ptr("WinDivert")
)

func init() {
	if err := checkForWow64(); err != nil {
		panic(err)
	}

	system32, err := windows.GetSystemDirectory()
	if err != nil {
		panic(err)
	}
	WinDivertSys = filepath.Join(system32, "WinDivert"+strconv.Itoa(32<<(^uint(0)>>63))+".sys")
	WinDivertDll = filepath.Join(system32, "WinDivert.dll")

	if err := InstallDriver(); err != nil {
		panic(err)
	}

	WinDivert = windows.MustLoadDLL("WinDivert.dll")
	WinDivertOpen = WinDivert.MustFindProc("WinDivertOpen")

	var vers = map[string]struct{}{
		"2.0": struct{}{},
		"2.1": struct{}{},
		"2.2": struct{}{},
	}

	hd, err := Open("false", LayerNetwork, PriorityDefault, FlagDefault)
	if err != nil {
		panic(err)
	}
	defer hd.Close()

	major, err := hd.GetParam(VersionMajor)
	if err != nil {
		panic(err)
	}

	minor, err := hd.GetParam(VersionMinor)
	if err != nil {
		panic(err)
	}

	if err := hd.Shutdown(ShutdownBoth); err != nil {
		panic(err)
	}

	ver := strings.Join([]string{strconv.Itoa(int(major)), strconv.Itoa(int(minor))}, ".")
	if _, ok := vers[ver]; !ok {
		s := ""
		for k, _ := range vers {
			s += k
		}
		panic(fmt.Errorf("unsupported version %v of windivert, only support %v", ver, s))
	}
}

func checkForWow64() error {
	var b bool
	err := windows.IsWow64Process(windows.CurrentProcess(), &b)
	if err != nil {
		return fmt.Errorf("Unable to determine whether the process is running under WOW64: %v", err)
	}
	if b {
		return fmt.Errorf("You must use the 64-bit version of WireGuard on this computer.")
	}
	return nil
}

func IoControlEx(h windows.Handle, code CtlCode, ioctl unsafe.Pointer, buf *byte, bufLen uint32, overlapped *windows.Overlapped) (iolen uint32, err error) {
	err = windows.DeviceIoControl(h, uint32(code), (*byte)(ioctl), uint32(unsafe.Sizeof(IoCtl{})), buf, bufLen, &iolen, overlapped)
	if err != windows.ERROR_IO_PENDING {
		return
	}

	err = windows.GetOverlappedResult(h, overlapped, &iolen, true)

	return
}

func IoControl(h windows.Handle, code CtlCode, ioctl unsafe.Pointer, buf *byte, bufLen uint32) (iolen uint32, err error) {
	event, _ := windows.CreateEvent(nil, 0, 0, nil)

	overlapped := windows.Overlapped{
		HEvent: event,
	}

	iolen, err = IoControlEx(h, code, ioctl, buf, bufLen, &overlapped)

	windows.CloseHandle(event)
	return
}

// RecvEx receives multiple packets
func (h *Handle) RecvEx(packets [][]byte, addrs []Address, flags uint64) (uint, uint, error) {
	if len(packets) == 0 || len(addrs) == 0 {
		return 0, 0, fmt.Errorf("empty packets or addresses buffer")
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

	return uint(readLen), uint(addrLen), nil
}

// SendEx sends multiple packets
func (h *Handle) SendEx(packets [][]byte, addrs []Address, flags uint64) (uint, error) {
	if len(packets) == 0 || len(addrs) == 0 {
		return 0, fmt.Errorf("empty packets or addresses buffer")
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

	return uint(writeLen), nil
}

// Recv receives a single packet
func (h *Handle) Recv(packet []byte, addr *Address) (uint, error) {
	nr, _, err := h.RecvEx([][]byte{packet}, []Address{*addr}, 0)
	if err != nil {
		return 0, err
	}
	return nr, nil
}

// SetParam sets a WinDivert parameter
func (h *Handle) SetParam(param Param, value uint64) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	ret := C.WinDivertSetParam(h.handle, C.WINDIVERT_PARAM(param), C.UINT64(value))
	if ret == 0 {
		return getLastError()
	}
	return nil
}

// Shutdown shuts down a WinDivert handle
func (h *Handle) Shutdown(how ShutdownType) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	ret := C.WinDivertShutdown(h.handle, C.WINDIVERT_SHUTDOWN(how))
	if ret == 0 {
		return getLastError()
	}
	return nil
}
