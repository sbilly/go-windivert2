package windivert

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// CtlCode represents a control code
type CtlCode uint32

// IoCtl represents an IO control structure
type IoCtl struct {
	// ... fields
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

// Error represents a WinDivert error
type Error syscall.Errno

// Error returns the error string
func (e Error) Error() string {
	return syscall.Errno(e).Error()
}

var (
	winDivert     = (*windows.DLL)(nil)
	winDivertOpen = (*windows.Proc)(nil)
	windivertsys  = ""
	windivertdll  = ""
	DeviceName    = windows.StringToUTF16Ptr("WinDivert")
)

func init() {
	if err := checkForWow64(); err != nil {
		panic(err)
	}

	system32, err := windows.GetSystemDirectory()
	if err != nil {
		panic(err)
	}
	windivertsys = filepath.Join(system32, "WinDivert"+strconv.Itoa(32<<(^uint(0)>>63))+".sys")
	windivertdll = filepath.Join(system32, "WinDivert.dll")

	if err := InstallDriver(); err != nil {
		panic(err)
	}

	winDivert = windows.MustLoadDLL("WinDivert.dll")
	winDivertOpen = winDivert.MustFindProc("WinDivertOpen")

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

type Handle struct {
	sync.Mutex
	windows.Handle
	rOverlapped windows.Overlapped
	wOverlapped windows.Overlapped
}

func Open(filter string, layer Layer, priority int16, flags uint64) (*Handle, error) {
	if priority < PriorityLowest || priority > PriorityHighest {
		return nil, fmt.Errorf("Priority %v is not Correct, Max: %v, Min: %v", priority, PriorityHighest, PriorityLowest)
	}

	filterPtr, err := windows.BytePtrFromString(filter)
	if err != nil {
		return nil, err
	}

	runtime.LockOSThread()
	hd, _, err := winDivertOpen.Call(uintptr(unsafe.Pointer(filterPtr)), uintptr(layer), uintptr(priority), uintptr(flags))
	runtime.UnlockOSThread()

	if windows.Handle(hd) == windows.InvalidHandle {
		return nil, Error(err.(windows.Errno))
	}

	rEvent, _ := windows.CreateEvent(nil, 0, 0, nil)
	wEvent, _ := windows.CreateEvent(nil, 0, 0, nil)

	return &Handle{
		Mutex:  sync.Mutex{},
		Handle: windows.Handle(hd),
		rOverlapped: windows.Overlapped{
			HEvent: rEvent,
		},
		wOverlapped: windows.Overlapped{
			HEvent: wEvent,
		},
	}, nil
}

func (h *Handle) Recv(buffer []byte, address *Address) (uint, error) {
	addrLen := uint(unsafe.Sizeof(Address{}))
	recv := recv{
		Addr:       uint64(uintptr(unsafe.Pointer(address))),
		AddrLenPtr: uint64(uintptr(unsafe.Pointer(&addrLen))),
	}

	iolen, err := IoControlEx(h.Handle, IoCtlRecv, unsafe.Pointer(&recv), &buffer[0], uint32(len(buffer)), &h.rOverlapped)
	if err != nil {
		return uint(iolen), Error(err.(syscall.Errno))
	}

	return uint(iolen), nil
}

func (h *Handle) Send(buffer []byte, address *Address) (uint, error) {
	send := send{
		Addr:    uint64(uintptr(unsafe.Pointer(address))),
		AddrLen: uint64(unsafe.Sizeof(Address{})),
	}

	iolen, err := IoControlEx(h.Handle, IoCtlSend, unsafe.Pointer(&send), &buffer[0], uint32(len(buffer)), &h.wOverlapped)
	if err != nil {
		return uint(iolen), Error(err.(syscall.Errno))
	}

	return uint(iolen), nil
}

func (h *Handle) Close() error {
	windows.CloseHandle(h.rOverlapped.HEvent)
	windows.CloseHandle(h.wOverlapped.HEvent)

	err := windows.CloseHandle(h.Handle)
	if err != nil {
		return Error(err.(syscall.Errno))
	}

	return nil
}
