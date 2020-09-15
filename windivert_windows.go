package windivert

import (
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Handle uintptr

type Address struct {
	IfIdx     uint32 // Packet's interface index
	SubIfIdx  uint32 // Packet's sub-interface index
	Direction uint8  // Packet's direction
}

const (
	false_             = 0
	InvalidHandleValue = -1

	LayerNetwork        = 0
	LayerNetworkForward = 1

	FlagSniff = 1
	FlagDrop  = 2

	DirectionOutbound = 0
	DirectionInbound  = 1

	ParamQueueLen  = 0
	ParamQueueTime = 1
	ParamQueueSize = 2

	ParamMajorVersion = 3
	ParamMinorVersion = 4
)

var (
	dll      *windows.DLL
	open     *windows.Proc
	recv     *windows.Proc
	send     *windows.Proc
	close_   *windows.Proc
	setParam *windows.Proc
	getParam *windows.Proc

	calcChecksums *windows.Proc

	DLLVersion string
)

func init() {
	dll = windows.MustLoadDLL("WinDivert")

	open = dll.MustFindProc("WinDivertOpen")
	recv = dll.MustFindProc("WinDivertRecv")
	send = dll.MustFindProc("WinDivertSend")
	close_ = dll.MustFindProc("WinDivertClose")
	setParam = dll.MustFindProc("WinDivertSetParam")
	getParam = dll.MustFindProc("WinDivertGetParam")

	calcChecksums = dll.MustFindProc("WinDivertHelperCalcChecksums")
}

func Open(filter string, layer, priority, flags int) (Handle, error) {
	str := make([]byte, len(filter)+1)
	copy(str, filter)
	r, _, err := open.Call(uintptr(unsafe.Pointer(&str[0])), uintptr(layer), uintptr(priority), uintptr(flags))
	if int(r) == InvalidHandleValue {
		return 0, err
	}
	return Handle(r), nil
}

func (h Handle) Close() error {
	r, _, err := close_.Call(uintptr(h))
	if r == false_ {
		return err
	}
	return nil
}

func (h Handle) Recv(packet []byte) (n int, addr Address, err error) {
	r, _, err := recv.Call(uintptr(h), uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)),
		uintptr(unsafe.Pointer(&n)), uintptr(unsafe.Pointer(&addr)))
	if r == false_ {
		return 0, addr, err
	}
	return n, addr, nil
}

func (h Handle) Send(packet []byte, addr Address) (n int, err error) {
	r, _, err := send.Call(uintptr(h), uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)),
		uintptr(unsafe.Pointer(&n)), uintptr(unsafe.Pointer(&addr)))
	if r == false_ {
		return 0, err
	}
	if len(packet) != n {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (h Handle) SetParam(param uintptr, value uint64) error {
	r, _, err := setParam.Call(uintptr(h), param, uintptr(value))
	if r == false_ {
		return err
	}
	return nil
}

func (h Handle) GetParam(param uintptr) (uint64, error) {
	var value uint64
	r, _, err := getParam.Call(uintptr(h), param, uintptr(unsafe.Pointer(&value)))
	if r == false_ {
		return 0, err
	}
	return value, nil
}

func CalcChecksums(packet []byte) []byte {
	calcChecksums.Call(uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)), 0)
	return packet
}

func (h Handle) GetVersion() (version string, err error) {
	var major, minor uint64
	r, _, err := getParam.Call(uintptr(h), ParamMajorVersion, uintptr(unsafe.Pointer(&major)))
	if r == false_ {
		return "", err
	}
	r, _, err = getParam.Call(uintptr(h), ParamMinorVersion, uintptr(unsafe.Pointer(&minor)))
	if r == false_ {
		return "", err
	}
	version = fmt.Sprintf("%d.%d.x", major, minor)
	return version, nil
}
