package windivert

// Layer represents WinDivert layers
type Layer int

const (
	LayerNetwork        Layer = 0
	LayerNetworkForward Layer = 1
	LayerFlow           Layer = 2
	LayerSocket         Layer = 3
	LayerReflect        Layer = 4
)

func (l Layer) String() string {
	switch l {
	case LayerNetwork:
		return "WINDIVERT_LAYER_NETWORK"
	case LayerNetworkForward:
		return "WINDIVERT_LAYER_NETWORK_FORWARD"
	case LayerFlow:
		return "WINDIVERT_LAYER_FLOW"
	case LayerSocket:
		return "WINDIVERT_LAYER_SOCKET"
	case LayerReflect:
		return "WINDIVERT_LAYER_REFLECT"
	default:
		return ""
	}
}

// Event represents WinDivert events
type Event int

const (
	EventNetworkPacket   Event = 0
	EventFlowEstablished Event = 1
	EventFlowDeleted     Event = 2
	EventSocketBind      Event = 3
	EventSocketConnect   Event = 4
	EventSocketListen    Event = 5
	EventSocketAccept    Event = 6
	EventSocketClose     Event = 7
	EventReflectOpen     Event = 8
	EventReflectClose    Event = 9
	EventEthernetFrame   Event = 10
)

// ShutdownType represents WinDivert shutdown types
type ShutdownType uint32

const (
	ShutdownRecv ShutdownType = 0
	ShutdownSend ShutdownType = 1
	ShutdownBoth ShutdownType = 2
)

// Param represents WinDivert parameters
type Param uint32

const (
	QueueLength  Param = 0
	QueueTime    Param = 1
	QueueSize    Param = 2
	VersionMajor Param = 3
	VersionMinor Param = 4
)

// Flags for WinDivertOpen()
const (
	FlagDefault   uint64 = 0
	FlagSniff     uint64 = 1
	FlagDrop      uint64 = 2
	FlagDebug     uint64 = 4
	FlagRecvOnly  uint64 = 8
	FlagSendOnly  uint64 = 16
	FlagNoInstall uint64 = 32
	FlagFragments uint64 = 64
)

// Default values
const (
	PriorityDefault    = 0
	QueueLengthDefault = 512
	QueueLengthMin     = 32
	QueueLengthMax     = 16384
	QueueTimeDefault   = 2000
	QueueTimeMin       = 100
	QueueTimeMax       = 16000
	QueueSizeDefault   = 4194304
	QueueSizeMin       = 65535
	QueueSizeMax       = 33554432
)

const (
	BatchMax = 0xff
	MTUMax   = 40 + 0xffff
)

// CtlCode represents a control code
type CtlCode uint32

const (
	CtlCodeInitialize = CtlCode(0x921)
	CtlCodeStartup    = CtlCode(0x922)
	CtlCodeRecv       = CtlCode(0x923)
	CtlCodeSend       = CtlCode(0x924)
	CtlCodeSetParam   = CtlCode(0x925)
	CtlCodeGetParam   = CtlCode(0x926)
	CtlCodeShutdown   = CtlCode(0x927)
)
