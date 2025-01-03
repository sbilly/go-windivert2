package windivert

import (
	"fmt"
	"unsafe"
)

// PacketInfo contains parsed packet information
type PacketInfo struct {
	IPv4Header   *IPv4Header
	IPv6Header   *IPv6Header
	ICMPHeader   *ICMPHeader
	ICMPv6Header *ICMPv6Header
	TCPHeader    *TCPHeader
	UDPHeader    *UDPHeader
	Data         []byte
}

// IPv4Header represents an IPv4 header (WINDIVERT_IPHDR)
type IPv4Header struct {
	HdrLength uint8  // Header length
	Version   uint8  // Version
	TOS       uint8  // Type of service
	Length    uint16 // Total length
	Id        uint16 // Identification
	FragOff   uint16 // Fragment offset
	TTL       uint8  // Time to live
	Protocol  uint8  // Protocol
	Checksum  uint16 // Checksum
	SrcAddr   uint32 // Source address
	DstAddr   uint32 // Destination address
}

// IPv6Header represents an IPv6 header (WINDIVERT_IPV6HDR)
type IPv6Header struct {
	Version      uint8    // Version
	TrafficClass uint8    // Traffic class
	FlowLabel    uint32   // Flow label
	Length       uint16   // Payload length
	NextHdr      uint8    // Next header
	HopLimit     uint8    // Hop limit
	SrcAddr      [16]byte // Source address
	DstAddr      [16]byte // Destination address
}

// ICMPHeader represents an ICMP header (WINDIVERT_ICMPHDR)
type ICMPHeader struct {
	Type     uint8  // Type
	Code     uint8  // Code
	Checksum uint16 // Checksum
	Body     uint32 // Body
}

// ICMPv6Header represents an ICMPv6 header (WINDIVERT_ICMPV6HDR)
type ICMPv6Header struct {
	Type     uint8  // Type
	Code     uint8  // Code
	Checksum uint16 // Checksum
	Body     uint32 // Body
}

// TCPHeader represents a TCP header (WINDIVERT_TCPHDR)
type TCPHeader struct {
	SrcPort   uint16 // Source port
	DstPort   uint16 // Destination port
	SeqNum    uint32 // Sequence number
	AckNum    uint32 // Acknowledgement number
	Reserved1 uint8  // Reserved
	Reserved2 uint8  // Reserved
	Reserved3 uint8  // Reserved
	Flags     uint8  // Flags
	Window    uint16 // Window
	Checksum  uint16 // Checksum
	UrgPtr    uint16 // Urgent pointer
}

// UDPHeader represents a UDP header (WINDIVERT_UDPHDR)
type UDPHeader struct {
	SrcPort  uint16 // Source port
	DstPort  uint16 // Destination port
	Length   uint16 // Length
	Checksum uint16 // Checksum
}

// Helper functions for parsing addresses
func ParseIPv4Address(str string) (uint32, error) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	var addr C.UINT32
	if ret := C.WinDivertHelperParseIPv4Address(cstr, &addr); ret == 0 {
		return 0, fmt.Errorf("failed to parse IPv4 address")
	}
	return uint32(addr), nil
}

func ParseIPv6Address(str string) ([4]uint32, error) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	var addr [4]C.UINT32
	if ret := C.WinDivertHelperParseIPv6Address(cstr, &addr[0]); ret == 0 {
		return [4]uint32{}, fmt.Errorf("failed to parse IPv6 address")
	}
	return *(*[4]uint32)(unsafe.Pointer(&addr)), nil
}
