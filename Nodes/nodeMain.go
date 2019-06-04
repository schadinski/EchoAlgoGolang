package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
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
	fmt.Println("Startup msg sent")

	// loop to get msgs while not send an echo or result msg
	for {
		byteArray := make([]byte, 1024)
		var currMsg Msg
		//fmt.Println("in for loop")
		// Read newest msg from udp connection in byteArray
		_, senderAddr, err := conn.ReadFromUDP(byteArray)
		if err != nil {
			fmt.Println("Error at ReadFromUDP:", err)
		}
		//fmt.Println("after read from udp")
		// cast byte[] to buffer
		buffer := bytes.NewBuffer(byteArray)
		// Decoding
		d := gob.NewDecoder(buffer)
		if err := d.Decode(&currMsg); err != nil {
			fmt.Println("Error at decode: ", err)
		}

		// Simulate network delay 0<x<=99 milliseconds
		NetworkDelay()

		//fmt.Println("after network delay")
		node.sendLogMsg(currMsg)

		switch currMsg.MsgType {
		case start:
			fmt.Println("got start msg")
			node.receiveStartMsg(currMsg)
		case info, echo:
			//fmt.Println("got info or echo msg")
			node.receiveInfoMsg(currMsg, senderAddr)
		default:
			fmt.Println("Received msg with unexpected type ", currMsg.MsgType)
		}
	}
	fmt.Println("Finished")
	for {

	}
}
