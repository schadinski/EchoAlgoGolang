package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	var logger logger

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

	chanUserInput := make(chan string)
	chanGotMsg := make(chan Msg)

	// Endless loop to wait for event
	for {
		go logger.receiveMsg(chanGotMsg)
		go logger.getUserInput(chanUserInput)
		select {
		case userInput := <-chanUserInput:
			// Trim user input in 3 fileds for instruction, ip and port
			splittedInput := strings.Fields(userInput)
			if len(splittedInput) == 3 {
				for i := 0; i < 3; i++ {
					splittedInput[i] = strings.Trim(splittedInput[i], " ")
					fmt.Println(splittedInput[i])
				}
				if ValidateUserInput(splittedInput) {
					initiatorAddr := BuildUDPAddr(splittedInput[1] + ":" + splittedInput[2])
					logger.startEchoAlgorithm(initiatorAddr)
				}
			}

		case currMsg := <-chanGotMsg:
			switch currMsg.MsgType {
			case logging:
				logger.printLoggingMsg(currMsg)
			case result:
				calc := 0
				for i := 1; i <= 15; i++ {
					calc += i
				}
				fmt.Println("Result should be ", calc)
				fmt.Println("Algorithm terminated. Result is: ", currMsg.Data)
			default:
				fmt.Println("Unknown messagetyp received")
			}
		}
	}
}
