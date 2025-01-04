//go:build windows
// +build windows

package utils

import (
	"syscall"
	"unsafe"
)

var (
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")

	procGetTcpTable2 = modiphlpapi.NewProc("GetTcpTable2")
	procGetUdpTable  = modiphlpapi.NewProc("GetUdpTable")
	procGetTcp6Table = modiphlpapi.NewProc("GetTcp6Table")
	procGetUdp6Table = modiphlpapi.NewProc("GetUdp6Table")
)

// TCPRow represents a TCP connection entry
type TCPRow struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

// TCP6Row represents a TCP IPv6 connection entry
type TCP6Row struct {
	LocalAddr  [4]uint32
	LocalPort  uint32
	RemoteAddr [4]uint32
	RemotePort uint32
	State      uint32
	OwningPid  uint32
}

// UDPRow represents a UDP connection entry
type UDPRow struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

// UDP6Row represents a UDP IPv6 connection entry
type UDP6Row struct {
	LocalAddr [4]uint32
	LocalPort uint32
	OwningPid uint32
}

// GetTCPTable retrieves the TCP connection table
func GetTCPTable() ([]TCPRow, error) {
	var size uint32
	var tcpTable []byte

	// Get the size of the table
	r1, _, err := procGetTcpTable2.Call(0, uintptr(unsafe.Pointer(&size)), 1)
	if r1 != 0 && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	tcpTable = make([]byte, size)
	r1, _, err = procGetTcpTable2.Call(
		uintptr(unsafe.Pointer(&tcpTable[0])),
		uintptr(unsafe.Pointer(&size)),
		1,
	)
	if r1 != 0 {
		return nil, err
	}

	count := *(*uint32)(unsafe.Pointer(&tcpTable[0]))
	rows := make([]TCPRow, count)

	for i := uint32(0); i < count; i++ {
		offset := 4 + i*uint32(unsafe.Sizeof(TCPRow{}))
		row := (*TCPRow)(unsafe.Pointer(&tcpTable[offset]))
		rows[i] = *row
	}

	return rows, nil
}

// GetUDPTable retrieves the UDP connection table
func GetUDPTable() ([]UDPRow, error) {
	var size uint32
	var udpTable []byte

	// Get the size of the table
	r1, _, err := procGetUdpTable.Call(0, uintptr(unsafe.Pointer(&size)), 1)
	if r1 != 0 && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	udpTable = make([]byte, size)
	r1, _, err = procGetUdpTable.Call(
		uintptr(unsafe.Pointer(&udpTable[0])),
		uintptr(unsafe.Pointer(&size)),
		1,
	)
	if r1 != 0 {
		return nil, err
	}

	count := *(*uint32)(unsafe.Pointer(&udpTable[0]))
	rows := make([]UDPRow, count)

	for i := uint32(0); i < count; i++ {
		offset := 4 + i*uint32(unsafe.Sizeof(UDPRow{}))
		row := (*UDPRow)(unsafe.Pointer(&udpTable[offset]))
		rows[i] = *row
	}

	return rows, nil
}

// GetTCP6Table retrieves the TCP IPv6 connection table
func GetTCP6Table() ([]TCP6Row, error) {
	var size uint32
	var tcp6Table []byte

	// Get the size of the table
	r1, _, err := procGetTcp6Table.Call(0, uintptr(unsafe.Pointer(&size)), 1)
	if r1 != 0 && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	tcp6Table = make([]byte, size)
	r1, _, err = procGetTcp6Table.Call(
		uintptr(unsafe.Pointer(&tcp6Table[0])),
		uintptr(unsafe.Pointer(&size)),
		1,
	)
	if r1 != 0 {
		return nil, err
	}

	count := *(*uint32)(unsafe.Pointer(&tcp6Table[0]))
	rows := make([]TCP6Row, count)

	for i := uint32(0); i < count; i++ {
		offset := 4 + i*uint32(unsafe.Sizeof(TCP6Row{}))
		row := (*TCP6Row)(unsafe.Pointer(&tcp6Table[offset]))
		rows[i] = *row
	}

	return rows, nil
}

// GetUDP6Table retrieves the UDP IPv6 connection table
func GetUDP6Table() ([]UDP6Row, error) {
	var size uint32
	var udp6Table []byte

	// Get the size of the table
	r1, _, err := procGetUdp6Table.Call(0, uintptr(unsafe.Pointer(&size)), 1)
	if r1 != 0 && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	udp6Table = make([]byte, size)
	r1, _, err = procGetUdp6Table.Call(
		uintptr(unsafe.Pointer(&udp6Table[0])),
		uintptr(unsafe.Pointer(&size)),
		1,
	)
	if r1 != 0 {
		return nil, err
	}

	count := *(*uint32)(unsafe.Pointer(&udp6Table[0]))
	rows := make([]UDP6Row, count)

	for i := uint32(0); i < count; i++ {
		offset := 4 + i*uint32(unsafe.Sizeof(UDP6Row{}))
		row := (*UDP6Row)(unsafe.Pointer(&udp6Table[offset]))
		rows[i] = *row
	}

	return rows, nil
}
