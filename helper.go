package windivert

import (
	"fmt"
	"unsafe"
)

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/lib -lwindivert
#include <windivert.h>
#include <stdlib.h>
*/
import "C"

// CalcChecksums calculates checksums for the packet
func CalcChecksums(packet []byte, addr *Address, flags uint64) error {
	if len(packet) == 0 {
		return fmt.Errorf("empty packet buffer")
	}

	caddr := (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(addr))
	ret := C.WinDivertHelperCalcChecksums(
		unsafe.Pointer(&packet[0]),
		C.UINT(len(packet)),
		caddr,
		C.UINT64(flags),
	)

	if ret == 0 {
		return getLastError()
	}

	return nil
}

// ParsePacket parses a network packet
func ParsePacket(packet []byte) (*PacketInfo, error) {
	if len(packet) == 0 {
		return nil, fmt.Errorf("empty packet buffer")
	}

	var info C.WINDIVERT_PACKET_INFO
	ret := C.WinDivertHelperParsePacket(
		unsafe.Pointer(&packet[0]),
		C.UINT(len(packet)),
		&info.IPv4Header,
		&info.IPv6Header,
		&info.ICMPHeader,
		&info.ICMPv6Header,
		&info.TCPHeader,
		&info.UDPHeader,
		&info.Data,
		&info.DataLen,
	)

	if ret == 0 {
		return nil, getLastError()
	}

	return (*PacketInfo)(unsafe.Pointer(&info)), nil
}

// FormatIPv4Address formats an IPv4 address
func FormatIPv4Address(addr uint32) string {
	var buf [16]C.char
	C.WinDivertHelperFormatIPv4Address(C.UINT32(addr), &buf[0], 16)
	return C.GoString(&buf[0])
}

// FormatIPv6Address formats an IPv6 address
func FormatIPv6Address(addr [4]uint32) string {
	var buf [46]C.char
	C.WinDivertHelperFormatIPv6Address((*C.UINT32)(&addr[0]), &buf[0], 46)
	return C.GoString(&buf[0])
}

// NtohIPv4Address converts a network byte order IPv4 address to host byte order
func NtohIPv4Address(addr uint32) uint32 {
	return uint32(C.WinDivertHelperNtohIPv4Address(C.UINT32(addr)))
}

// NtohIPv6Address converts a network byte order IPv6 address to host byte order
func NtohIPv6Address(addr [4]uint32) [4]uint32 {
	var result [4]uint32
	C.WinDivertHelperNtohIPv6Address((*C.UINT32)(&addr[0]), (*C.UINT32)(&result[0]))
	return result
}

// HtonIPv4Address converts a host byte order IPv4 address to network byte order
func HtonIPv4Address(addr uint32) uint32 {
	return uint32(C.WinDivertHelperHtonIPv4Address(C.UINT32(addr)))
}

// HtonIPv6Address converts a host byte order IPv6 address to network byte order
func HtonIPv6Address(addr [4]uint32) [4]uint32 {
	var result [4]uint32
	C.WinDivertHelperHtonIPv6Address((*C.UINT32)(&addr[0]), (*C.UINT32)(&result[0]))
	return result
}

// HashPacket calculates a 64bit hash value of the given packet
func HashPacket(packet []byte, seed uint64) (uint64, error) {
	if len(packet) == 0 {
		return 0, fmt.Errorf("empty packet buffer")
	}

	hash := C.WinDivertHelperHashPacket(
		unsafe.Pointer(&packet[0]),
		C.UINT(len(packet)),
		C.UINT64(seed),
	)

	return uint64(hash), nil
}

// DecrementTTL decrements the TTL/HopLimit field of an IP packet
func DecrementTTL(packet []byte) error {
	if len(packet) == 0 {
		return fmt.Errorf("empty packet buffer")
	}

	ret := C.WinDivertHelperDecrementTTL(
		unsafe.Pointer(&packet[0]),
		C.UINT(len(packet)),
	)

	if ret == 0 {
		return fmt.Errorf("TTL/HopLimit would become 0")
	}

	return nil
}

// CompileFilter compiles a filter string into an object representation
func CompileFilter(filter string, layer Layer) (string, error) {
	cfilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cfilter))

	var errorStr *C.char
	var errorPos C.UINT
	var objBuf [1024]C.char

	ret := C.WinDivertHelperCompileFilter(
		cfilter,
		C.WINDIVERT_LAYER(layer),
		&objBuf[0],
		C.UINT(len(objBuf)),
		&errorStr,
		&errorPos,
	)

	if ret == 0 {
		return "", fmt.Errorf("filter compilation failed at position %d: %s",
			uint(errorPos), C.GoString(errorStr))
	}

	return C.GoString(&objBuf[0]), nil
}

// EvalFilter evaluates a packet against a filter string
func EvalFilter(filter string, packet []byte, addr *Address) (bool, error) {
	if len(packet) == 0 {
		return false, fmt.Errorf("empty packet buffer")
	}

	cfilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cfilter))

	caddr := (*C.WINDIVERT_ADDRESS)(unsafe.Pointer(addr))
	ret := C.WinDivertHelperEvalFilter(
		cfilter,
		unsafe.Pointer(&packet[0]),
		C.UINT(len(packet)),
		caddr,
	)

	return ret != 0, nil
}

// FormatFilter formats a filter string
func FormatFilter(filter string, layer Layer) (string, error) {
	cfilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cfilter))

	var buf [1024]C.char
	ret := C.WinDivertHelperFormatFilter(
		cfilter,
		C.WINDIVERT_LAYER(layer),
		&buf[0],
		C.UINT(len(buf)),
	)

	if ret == 0 {
		return "", getLastError()
	}

	return C.GoString(&buf[0]), nil
}

// Htons converts a 16-bit number from host to network byte order
func Htons(x uint16) uint16 {
	return uint16(C.WinDivertHelperHtons(C.UINT16(x)))
}

// Htonl converts a 32-bit number from host to network byte order
func Htonl(x uint32) uint32 {
	return uint32(C.WinDivertHelperHtonl(C.UINT32(x)))
}

// Htonll converts a 64-bit number from host to network byte order
func Htonll(x uint64) uint64 {
	return uint64(C.WinDivertHelperHtonll(C.UINT64(x)))
}
