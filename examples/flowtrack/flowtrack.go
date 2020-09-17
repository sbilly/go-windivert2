package main

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/sbilly/go-windivert2"
)

func init() {}

func main() {
	handle, err := windivert.Open("true", windivert.LayerFlow, 0, windivert.FlagSniff|windivert.FlagReadOnly)
	if err != nil {
		log.Fatal("Open() error", err)
	}
	defer func() {
		err = handle.Close()
		if err != nil {
			log.Fatal("Close() error: ", err)
		}
	}()

	packet := make([]byte, 4096)
	log.Println("Recv() loop started.")
	for {
		n, addr, err := handle.Recv(packet)
		if err != nil {
			log.Fatal("Recv() error", err)
		}
		buff := bytes.NewBuffer(addr.Data[:])

		if n != 0 {
			log.Println(packet, addr, n)
		}
		log.Printf("Timestamp: %d", addr.Timestamp)
		log.Printf("Layer: 0x%x", addr.Layer)
		log.Printf("Event: 0x%x", addr.Event)

		log.Printf("Flags: 0x%x", addr.Flags)
		sniff, _ := addr.IsSniffed()
		log.Printf("IsSniffed: %t", sniff)
		outbound, _ := addr.IsOutbound()
		log.Printf("IsOutbound: %t", outbound)
		loopback, _ := addr.IsLoopback()
		log.Printf("IsLoopback: %t", loopback)
		imposter, _ := addr.IsImpostor()
		log.Printf("IsImpostor: %t", imposter)
		ipv6, _ := addr.IsIPv6()
		log.Printf("IsIPv6: %t", ipv6)
		ipchecksum, _ := addr.IsIPChecksum()
		log.Printf("IsIPChecksum: %t", ipchecksum)
		tcpchecksum, _ := addr.IsTCPChecksum()
		log.Printf("IsTCPChecksum: %t", tcpchecksum)
		udpchecksum, _ := addr.IsUDPChecksum()
		log.Printf("IsUDPChecksum: %t", udpchecksum)

		log.Printf("Reserved1: 0x%x", addr.Reserved1)
		log.Printf("Reserved2: 0x%04x", addr.Reserved2)
		// log.Printf("Data: 0x%x", addr.Data)

		switch addr.Layer {
		case windivert.LayerFlow:
			flow := windivert.DataFlow{}

			err = binary.Read(buff, binary.LittleEndian, &flow)
			if err != nil {
				log.Fatal("Read addr.Event error:", err)
			}

			switch addr.Event {
			case windivert.EventNetworkFlowEstablished:
				// Flow established:
				log.Printf("================= Flow established =================")
			case windivert.EventNetworkFlowDelete:
				// Flow deleted:
				log.Printf("================= Flow deleted =================")
			}

			log.Printf("flow.EndpointId: %d", flow.EndpointId)
			log.Printf("flow.ParentEndpointId: %d", flow.ParentEndpointId)
			log.Printf("flow.ProcessId: %d", flow.ProcessId)
			log.Printf("flow.LocalAddr: %s", windivert.FormatIPAddress(flow.LocalAddr))
			log.Printf("flow.RemoteAddr: %s", windivert.FormatIPAddress(flow.RemoteAddr))
			log.Printf("flow.LocalPort: %d", flow.LocalPort)
			log.Printf("flow.RemotePort: %d", flow.RemotePort)
			log.Printf("flow.Protocol: %d", flow.Protocol)
		}

	}
}
