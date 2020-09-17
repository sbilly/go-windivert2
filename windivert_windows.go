package windivert

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Handle uintptr

// https://github.com/basil00/Divert/blob/master/include/windivert.h
type DataNetwork struct {
	IfIdx    uint32 /* Packet's interface index. */
	SubIfIdx uint32 /* Packet's sub-interface index. */
}

type DataFlow struct {
	EndpointId       uint64    /* Endpoint ID. */
	ParentEndpointId uint64    /* Parent Endpoint ID. */
	ProcessId        uint32    /* Process ID. */
	LocalAddr        [4]uint32 /* Local address. */
	RemoteAddr       [4]uint32 /* Remote address. */
	LocalPort        uint16    /* Local port. */
	RemotePort       uint16    /* Remote port. */
	Protocol         uint8     /* Protocol. */
}

type DataSocket struct {
	EndpointId       uint64    /* Endpoint ID. */
	ParentEndpointId uint64    /* Parent Endpoint ID. */
	ProcessId        uint32    /* Process ID. */
	LocalAddr        [4]uint32 /* Local address. */
	RemoteAddr       [4]uint32 /* Remote address. */
	LocalPort        uint16    /* Local port. */
	RemotePort       uint16    /* Remote port. */
	Protocol         uint8     /* Protocol. */
}

type DataReflect struct {
	Timestamp int64  /* Handle open time. */
	ProcessId uint32 /* Handle process ID. */
	Layer     int    /* Handle layer. */
	Flags     uint64 /* Handle flags. */
	Priority  int16  /* Handle priority. */
}

type Address struct {
	Timestamp int64  /* Packet's timestamp. */
	Layer     uint8  /* Packet's layer. */
	Event     uint8  /* Packet event. */
	Flags     uint8  /* Packet Flags: Sniffed, Outbound, Loopback, Impostor, IPv6, IPChecksum, TCPChecksum, UDPChecksum */
	Reserved1 uint8  /* Reserved1 */
	Reserved2 uint32 /* Reserved2 */
	Data      [64]byte
}

const (
	false_             = 0
	InvalidHandleValue = -1

	// https://github.com/basil00/Divert/blob/master/include/windivert.h
	PacketFlagSniffed     = 1      /* Packet was sniffed? */
	PacketFlagOutbound    = 1 << 1 /* Packet is outound? */
	PacketFlagLoopback    = 1 << 2 /* Packet is loop? */
	PacketFlagImpostor    = 1 << 3 /* Packet is Impostor? */
	PacketFlagIPv6        = 1 << 4 /* Packet is IPv6? */
	PacketFlagIPChecksum  = 1 << 5 /* Packet has valid IPv4 checksum? */
	PacketFlagTCPChecksum = 1 << 6 /* Packet has valid TCP checksum? */
	PacketFlagUDPChecksum = 1 << 7 /* Packet has valid UDP checksum? */

	// WinDivert events.
	EventNetworkPacket          = 0 /* Network packet. */
	EventNetworkFlowEstablished = 1 /* Flow established. */
	EventNetworkFlowDelete      = 2 /* Flow deleted. */
	EventSocketBind             = 3 /* Socket bind. */
	EventSocketConnect          = 4 /* Socket connect. */
	EventSocketListen           = 5 /* Socket listen. */
	EventSocketAccept           = 6 /* Socket accept. */
	EventSocketClose            = 7 /* Socket close. */
	EventReflectOpen            = 8 /* WinDivert handle opened. */
	EventReflectClose           = 9 /* WinDivert handle closed. */

	LayerNetwork        = 0
	LayerNetworkForward = 1
	LayerFlow           = 2
	LayerSocket         = 3
	LayerReflect        = 4

	FlagSniff     = 1
	FlagDrop      = 2
	FlagRecvOnly  = 4
	FlagReadOnly  = 4
	FlagSendOnly  = 8
	FlagWriteOnly = 8
	FlagNoInstall = 16
	FlagFragment  = 32

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

func CalcChecksums(packet []byte) []byte {
	calcChecksums.Call(uintptr(unsafe.Pointer(&packet[0])), uintptr(len(packet)), 0)
	return packet
}

func FormatIPAddress(addr [4]uint32) net.IP {
	b := make([]byte, 16)
	binary.BigEndian.PutUint32(b[0:4], addr[3])
	binary.BigEndian.PutUint32(b[4:8], addr[2])
	binary.BigEndian.PutUint32(b[8:12], addr[1])
	binary.BigEndian.PutUint32(b[12:16], addr[0])
	return net.IP(b)
}

func FormatIPv4Address(addr [4]uint32) net.IP {
	return FormatIPAddress(addr).To4()
}

func FormatIPv6Address(addr [4]uint32) net.IP {
	return FormatIPAddress(addr).To16()
}

func (addr Address) IsSniffed() (bool, error) {
	result := addr.Flags & PacketFlagSniffed
	return result == PacketFlagSniffed, nil
}

func (addr Address) IsOutbound() (bool, error) {
	result := addr.Flags & PacketFlagOutbound
	return result == PacketFlagOutbound, nil
}

func (addr Address) IsLoopback() (bool, error) {
	result := addr.Flags & PacketFlagLoopback
	return result == PacketFlagLoopback, nil
}

func (addr Address) IsImpostor() (bool, error) {
	result := addr.Flags & PacketFlagImpostor
	return result == PacketFlagImpostor, nil
}

func (addr Address) IsIPv6() (bool, error) {
	result := addr.Flags & PacketFlagIPv6
	return result == PacketFlagIPv6, nil
}

func (addr Address) IsIPChecksum() (bool, error) {
	result := addr.Flags & PacketFlagIPChecksum
	return result == PacketFlagIPChecksum, nil
}

func (addr Address) IsTCPChecksum() (bool, error) {
	result := addr.Flags & PacketFlagTCPChecksum
	return result == PacketFlagTCPChecksum, nil
}

func (addr Address) IsUDPChecksum() (bool, error) {
	result := addr.Flags & PacketFlagUDPChecksum
	return result == PacketFlagUDPChecksum, nil
}
