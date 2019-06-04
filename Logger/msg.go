package main

import (
	"fmt"
	"net"
)

// BuildUDPAddr builds a UDPAdrress for use in net.ListenUDP
func BuildUDPAddr(service string) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		fmt.Println("Error resolve address:", err)
	}
	return addr
}

// ValidateUserInput checks given input is
// expected instruction and valid ip
func ValidateUserInput(input []string) bool {
	result := false

	testInput := net.ParseIP(input[1] + ":" + input[2])

	if input[0] == "SUM" && testInput.To4() == nil {
		result = true
	}

	return result
}

// Types of msgs
type msgEnum int

const (
	start   msgEnum = 0
	info    msgEnum = 1
	echo    msgEnum = 2
	result  msgEnum = 3
	logging msgEnum = 4
)

// Msg represents a message in the system
type Msg struct {
	SenderAddr *net.UDPAddr
	MsgType    msgEnum
	Data       string
}
