package windivert

// IPv4Header represents an IPv4 header
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

// IPv6Header represents an IPv6 header
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

// ICMPHeader represents an ICMP header
type ICMPHeader struct {
	Type     uint8  // Type
	Code     uint8  // Code
	Checksum uint16 // Checksum
	Body     uint32 // Body
}

// ICMPv6Header represents an ICMPv6 header
type ICMPv6Header struct {
	Type     uint8  // Type
	Code     uint8  // Code
	Checksum uint16 // Checksum
	Body     uint32 // Body
}

// TCPHeader represents a TCP header
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

// UDPHeader represents a UDP header
type UDPHeader struct {
	SrcPort  uint16 // Source port
	DstPort  uint16 // Destination port
	Length   uint16 // Length
	Checksum uint16 // Checksum
}
