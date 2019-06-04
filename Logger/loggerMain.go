package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

func main() {

	fmt.Println("logger is up and running")
	var logger logger

	// logger.allAddrresses = []string{
	// 	"127.0.0.1:1053",
	// 	"127.0.0.1:1063",
	// 	"127.0.0.1:1073",
	// 	"127.0.0.1:1083",
	// 	"127.0.0.1:1093",
	// 	"127.0.0.1:2003",
	// 	"127.0.0.1:2013",
	// 	"127.0.0.1:2023",
	// 	"127.0.0.1:2033",
	// 	"127.0.0.1:2043",
	// 	"127.0.0.1:2053",
	// 	"127.0.0.1:2063",
	// 	"127.0.0.1:2073",
	// 	"127.0.0.1:2083",
	// 	"127.0.0.1:2093"}

	// Build loggers address
	hostName := "127.0.0.1"
	portNum := "1042"
	service := hostName + ":" + portNum

	// Get logger network addrress from string
	logger.loggerAddress = BuildUDPAddr(service)

	// Bind logger to its address
	conn, err := net.ListenUDP("udp", logger.loggerAddress)
	if err != nil {
		fmt.Println("Error at listenUDP:", err)
	}
	defer conn.Close()
	logger.conn = conn

	//chanUserInput := make(chan string)
	//chanGotMsg := make(chan Msg)
	//chanMsgReceiveDone := make(chan bool)
	//chanUserInputDone := make(chan bool)

	for {
		validInput := logger.getUserInput()
		if validInput {
			break
		}
	}
	// Endless loop to wait for event
	for {
		//		go logger.receiveMsg(chanGotMsg)
		//currMsg := logger.receiveMsg()
		byteArray := make([]byte, 1024)
		var currMsg Msg

		// Read newest msg from udp connection in byteArray
		//		_, nodeAddr, err := conn.ReadFromUDP(byteArray)
		_, _, err := logger.conn.ReadFromUDP(byteArray)
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

		//go logger.getUserInput(chanUserInput)
		//	select {
		//case userInput := <-chanUserInput:
		//fmt.Println("user input is ", userInput)
		// // Trim user input in 3 fileds for instruction, ip and port
		// splittedInput := strings.Fields(userInput)
		// //fmt.Println("input len is ", len(splittedInput))
		// if len(splittedInput) == 3 {
		// 	//fmt.Println("in if == 3")
		// 	for i := 0; i < 3; i++ {
		// 		//fmt.Println("in for")
		// 		splittedInput[i] = strings.Trim(splittedInput[i], " ")
		// 		fmt.Println(splittedInput[i])
		// 	}
		// 	//fmt.Println("after for")
		// 	if ValidateUserInput(splittedInput) {
		// 		//fmt.Println("input is valid")
		// 		initiatorAddr := BuildUDPAddr(splittedInput[1] + ":" + splittedInput[2])
		// 		logger.startEchoAlgorithm(initiatorAddr)
		// 	}
		// }

		//case currMsg := <-chanGotMsg:
		//fmt.Println("case upd msg")
		// Check msg type
		switch currMsg.MsgType {
		case logging:
			//	fmt.Println("Got log msg")
			logger.printLoggingMsg(currMsg)
		case result:
			logger.printResult(currMsg)
			//close(chanGotMsg)
		default:
			fmt.Println("Unknown messagetyp received")
		}
		//<-chanMsgReceiveDone

		//	}
	}

}
