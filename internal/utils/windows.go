//go:build windows
// +build windows

package utils

import (
	"syscall"
	"unsafe"
)

// GetTCPTable retrieves the TCP connection table
func GetTCPTable() ([]TCPRow, error) {
	var size uint32
	var tcpTable []byte

	// Get the size of the table
	err := syscall.GetTcpTable2(nil, &size, true)
	if err != nil && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	tcpTable = make([]byte, size)
	err = syscall.GetTcpTable2(&tcpTable[0], &size, true)
	if err != nil {
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
	err := syscall.GetUdpTable(&udpTable[0], &size, true)
	if err != nil && err != syscall.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	udpTable = make([]byte, size)
	err = syscall.GetUdpTable(&udpTable[0], &size, true)
	if err != nil {
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

// TCPRow represents a TCP connection entry
type TCPRow struct {
	State      uint32
	LocalAddr  uint32
	LocalPort  uint32
	RemoteAddr uint32
	RemotePort uint32
	OwningPid  uint32
}

// UDPRow represents a UDP connection entry
type UDPRow struct {
	LocalAddr uint32
	LocalPort uint32
	OwningPid uint32
}

// GetTCP6Table retrieves the TCP IPv6 connection table
func GetTCP6Table() ([]TCP6Row, error) {
	// Implementation for IPv6 TCP table
	return nil, nil
}

// GetUDP6Table retrieves the UDP IPv6 connection table
func GetUDP6Table() ([]UDP6Row, error) {
	// Implementation for IPv6 UDP table
	return nil, nil
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

// UDP6Row represents a UDP IPv6 connection entry
type UDP6Row struct {
	LocalAddr [4]uint32
	LocalPort uint32
	OwningPid uint32
}
