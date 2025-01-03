package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 打开一个 WinDivert 句柄
	handle, err := windivert.Open("true", windivert.LayerNetwork, 0, 0)
	if err != nil {
		fmt.Printf("Error opening WinDivert handle: %v\n", err)
		return
	}
	defer handle.Close()

	// 处理 Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 数据包处理循环
	packet := make([]byte, 1500)
	for {
		select {
		case <-sigCh:
			return
		default:
			n, addr, err := handle.Recv(packet)
			if err != nil {
				fmt.Printf("Error receiving packet: %v\n", err)
				continue
			}

			// 重新注入数据包
			_, err = handle.Send(packet[:n], addr)
			if err != nil {
				fmt.Printf("Error sending packet: %v\n", err)
			}
		}
	}
}
