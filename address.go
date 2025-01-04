package windivert

import (
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

// Ethernet represents ethernet layer information
type Ethernet struct {
	InterfaceIndex    uint32
	SubInterfaceIndex uint32
	_                 [7]uint64
}

// Network represents network layer information
type Network struct {
	InterfaceIndex    uint32
	SubInterfaceIndex uint32
	_                 [7]uint64
}

// Socket represents socket layer information
type Socket struct {
	EndpointID       uint64
	ParentEndpointID uint64
	ProcessID        uint32
	LocalAddress     [16]uint8
	RemoteAddress    [16]uint8
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint8
	_                [3]uint8
	_                uint32
}

// Flow represents flow layer information
type Flow struct {
	EndpointID       uint64
	ParentEndpointID uint64
	ProcessID        uint32
	LocalAddress     [16]uint8
	RemoteAddress    [16]uint8
	LocalPort        uint16
	RemotePort       uint16
	Protocol         uint8
	_                [3]uint8
	_                uint32
}

// Reflect represents reflect layer information
type Reflect struct {
	TimeStamp int64
	ProcessID uint32
	layer     uint32
	Flags     uint64
	Priority  int16
	_         int16
	_         int32
	_         [4]uint64
}

// Layer returns the layer type for reflect information
func (r *Reflect) Layer() Layer {
	return Layer(r.layer)
}

// Address represents a WinDivert address
type Address struct {
	Timestamp      int64
	LayerType      Layer // renamed from Layer
	EventType      Event // renamed from Event
	IsSniffed      uint8 // renamed from Sniffed
	IsOutbound     uint8 // renamed from Outbound
	HasIPChecksum  uint8 // renamed from IPChecksum
	HasTCPChecksum uint8 // renamed from TCPChecksum
	HasUDPChecksum uint8 // renamed from UDPChecksum
	Flags          uint8
	union          [64]byte
	length         uint64
}

// GetLayer returns the layer type
func (a *Address) Layer() Layer {
	return a.LayerType
}

// SetLayer sets the layer type
func (a *Address) SetLayer(layer Layer) {
	a.LayerType = layer
}

// GetEvent returns the event type
func (a *Address) Event() Event {
	return Event(a.EventType)
}

// SetEvent sets the event type
func (a *Address) SetEvent(event Event) {
	a.EventType = Event(event)
}

// IsSniffed returns whether the packet was sniffed
func (a *Address) Sniffed() bool {
	return (a.Flags & uint8(0x01<<0)) == uint8(0x01<<0)
}

// SetSniffed sets the sniffed flag
func (a *Address) SetSniffed() {
	a.Flags |= uint8(0x01 << 0)
}

// UnsetSniffed unsets the sniffed flag
func (a *Address) UnsetSniffed() {
	a.Flags &= ^uint8(0x01 << 0)
}

// IsOutbound returns whether the packet is outbound
func (a *Address) Outbound() bool {
	return (a.Flags & uint8(0x01<<1)) == uint8(0x01<<1)
}

// SetOutbound sets the outbound flag
func (a *Address) SetOutbound() {
	a.Flags |= uint8(0x01 << 1)
}

// UnsetOutbound unsets the outbound flag
func (a *Address) UnsetOutbound() {
	a.Flags &= ^uint8(0x01 << 1)
}

// HasIPChecksum returns whether IP checksum is present
func (a *Address) IPChecksum() bool {
	return (a.Flags & uint8(0x01<<5)) == uint8(0x01<<5)
}

// SetIPChecksum sets the IP checksum flag
func (a *Address) SetIPChecksum() {
	a.Flags |= uint8(0x01 << 5)
}

// UnsetIPChecksum unsets the IP checksum flag
func (a *Address) UnsetIPChecksum() {
	a.Flags &= ^uint8(0x01 << 5)
}

// HasTCPChecksum returns whether TCP checksum is present
func (a *Address) TCPChecksum() bool {
	return (a.Flags & uint8(0x01<<6)) == uint8(0x01<<6)
}

// SetTCPChecksum sets the TCP checksum flag
func (a *Address) SetTCPChecksum() {
	a.Flags |= uint8(0x01 << 6)
}

// UnsetTCPChecksum unsets the TCP checksum flag
func (a *Address) UnsetTCPChecksum() {
	a.Flags &= ^uint8(0x01 << 6)
}

// HasUDPChecksum returns whether UDP checksum is present
func (a *Address) UDPChecksum() bool {
	return (a.Flags & uint8(0x01<<7)) == uint8(0x01<<7)
}

// SetUDPChecksum sets the UDP checksum flag
func (a *Address) SetUDPChecksum() {
	a.Flags |= uint8(0x01 << 7)
}

// UnsetUDPChecksum unsets the UDP checksum flag
func (a *Address) UnsetUDPChecksum() {
	a.Flags &= ^uint8(0x01 << 7)
}

func (a *Address) Length() uint32 {
	return uint32(a.length >> 12)
}

func (a *Address) SetLength(n uint32) {
	a.length = uint64(n << 12)
}

func (a *Address) Ethernet() *Ethernet {
	return (*Ethernet)(unsafe.Pointer(&a.union))
}

func (a *Address) Network() *Network {
	return (*Network)(unsafe.Pointer(&a.union))
}

func (a *Address) Socket() *Socket {
	return (*Socket)(unsafe.Pointer(&a.union))
}

func (a *Address) Flow() *Flow {
	return (*Flow)(unsafe.Pointer(&a.union))
}

func (a *Address) Reflect() *Reflect {
	return (*Reflect)(unsafe.Pointer(&a.union))
}

type AddressHelper interface {
	CalcChecksums(packet []byte, flags uint64) error
	ParseIPv4Header(packet []byte) (*IPv4Header, error)
	ParseIPv6Header(packet []byte) (*IPv6Header, error)
	ParseTCPHeader(packet []byte) (*TCPHeader, error)
	ParseUDPHeader(packet []byte) (*UDPHeader, error)
}
