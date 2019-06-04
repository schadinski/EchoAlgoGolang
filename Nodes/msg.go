package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// BuildUDPAddr builds a UDPAdrress for use in net.ListenUDP
func BuildUDPAddr(service string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		fmt.Println("Error resolve address:", err)
	}
	return addr
}

// NetworkDelay simulates a network delay between 0 and 100
func NetworkDelay() {
	duration := rand.Intn(100)
	time.Sleep(time.Duration(duration) * time.Microsecond)
}

// Enum for message types
type msgEnum int

const (
	start   msgEnum = 0
	info    msgEnum = 1
	echo    msgEnum = 2
	result  msgEnum = 3
	logging msgEnum = 4
)

// Msg represents a msg to communicate with other nodes and the logger
type Msg struct {
	SenderAddr *net.UDPAddr
	MsgType    msgEnum
	Data       string
}

func (msg Msg) getStringForType() string {
	var typ string
	switch msg.MsgType {
	case start:
		typ = "start"
	case info:
		typ = "info"
	case echo:
		typ = "echo"
	case result:
		typ = "result"
	case logging:
		typ = "logging"
	default:
		fmt.Println("Msg haves unexpected type")
	}
	return typ
}
