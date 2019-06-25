package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {

	node := newNode(os.Args)
	// "Bind" to local address
	conn, err := net.ListenUDP("udp", node.localAddr)
	if err != nil {
		fmt.Print("Error at listenUDP:", err)
	}
	defer conn.Close()
	node.conn = conn

	// Send logger startup msg
	var startupMsg Msg
	startupMsg.SenderAddr = node.localAddr
	startupMsg.MsgType = logging
	startupMsg.Data = "Node " + node.localAddr.String() + " is up.\n"
	node.sendLogMsg(startupMsg)

	for {
		byteArray := make([]byte, 1024)
		var currMsg Msg
		// Read newest msg from udp connection in byteArray
		_, senderAddr, err := conn.ReadFromUDP(byteArray)
		if err != nil {
			fmt.Println("Error at ReadFromUDP:", err)
		}
		// cast byte[] to buffer
		buffer := bytes.NewBuffer(byteArray)
		// Decoding
		d := gob.NewDecoder(buffer)
		if err := d.Decode(&currMsg); err != nil {
			fmt.Println("Error at decode: ", err)
		}

		// Simulate network delay 0<x<=99 milliseconds
		NetworkDelay()

		node.sendLogMsg(currMsg)

		switch currMsg.MsgType {
		case start:
			node.receiveStartMsg(currMsg)
		case info:
			node.receiveIEMsg(&currMsg, senderAddr)
		case echo:
			peersMem, err := strconv.Atoi(currMsg.Data)
			if err != nil {
				fmt.Println("Error at Atoi: ", err)
			}
			node.sumOfMem += peersMem
			node.receiveIEMsg(&currMsg, senderAddr)
		default:
			fmt.Println("Received msg with unexpected type ", currMsg.getStringForType())
		}
	}
}
