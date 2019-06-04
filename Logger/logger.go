package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type logger struct {
	//allAddrresses []string
	loggerAddress *net.UDPAddr
	conn          *net.UDPConn
}

func (l *logger) startEchoAlgorithm(startAddr *net.UDPAddr) {
	// build startMsg
	startMsg := Msg{l.loggerAddress, start, "0"}

	var buffer bytes.Buffer
	byteArray := make([]byte, 1024)

	// Encode
	e := gob.NewEncoder(&buffer)
	err := e.Encode(startMsg)
	if err != nil {
		fmt.Println("Error at encode start msg: ", err)
	}

	// Cast buffer to byte[]
	byteArray, err = ioutil.ReadAll(&buffer)
	if err != nil {
		fmt.Println("Error at ReadAll:", err)
	}
	// Send start msg
	_, err = l.conn.WriteToUDP(byteArray, startAddr)
	if err != nil {
		fmt.Println("Error at conn.Write:", err)
	}
}

// Print data from logging msg
// msg.data is build from sender node
func (l *logger) printLoggingMsg(msg Msg) {
	fmt.Println("\n ", time.Now())
	fmt.Println("logging msg received from ", msg.SenderAddr)
	fmt.Println("Data is:\n", msg.Data)
	//fmt.Println(msg.SenderAddr + "send msg to ", )
}

// print the final result of echo algorithm
func (l *logger) printResult(msg Msg) {
	fmt.Println(" Algorithm terminated. Result is: ", msg.Data)
}

// to get msg from node
// func (l *logger) receiveMsg(chanGotMsg chan Msg) {
// 	//for {
// 	//fmt.Println("msg received")
// 	byteArray := make([]byte, 1024)
// 	var currMsg Msg

// 	// Read newest msg from udp connection in byteArray
// 	//		_, nodeAddr, err := conn.ReadFromUDP(byteArray)
// 	_, _, err := l.conn.ReadFromUDP(byteArray)
// 	if err != nil {
// 		fmt.Println("Error at ReadFromUDP:", err)
// 	}

// 	// cast byte[] to buffer
// 	buffer := bytes.NewBuffer(byteArray)
// 	// Decoding
// 	d := gob.NewDecoder(buffer)
// 	if err := d.Decode(&currMsg); err != nil {
// 		fmt.Println("Error at decode: ", err)
// 	}

// 	chanGotMsg <- currMsg
// 	//}
// }

// Function in background to get user input
// func (l *logger) getUserInput(chanUserInput chan string) {
// 	reader := bufio.NewReader(os.Stdin)
// 	input, _ := reader.ReadString('\n')
// 	fmt.Println(input)
// 	chanUserInput <- input
// }

func (l *logger) getUserInput() bool {
	ret := false
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	fmt.Println(input)
	// Trim user input in 3 fileds for instruction, ip and port
	splittedInput := strings.Fields(input)
	//fmt.Println("input len is ", len(splittedInput))
	if len(splittedInput) == 3 {
		//fmt.Println("in if == 3")
		for i := 0; i < 3; i++ {
			//fmt.Println("in for")
			splittedInput[i] = strings.Trim(splittedInput[i], " ")
			fmt.Println(splittedInput[i])
		}
		//fmt.Println("after for")
		if ValidateUserInput(splittedInput) {
			//fmt.Println("input is valid")
			initiatorAddr := BuildUDPAddr(splittedInput[1] + ":" + splittedInput[2])
			l.startEchoAlgorithm(initiatorAddr)
			ret = true
		}
	}
	return ret
}

func (l *logger) receiveMsg() Msg {
	//for {
	//fmt.Println("msg received")
	byteArray := make([]byte, 1024)
	var currMsg Msg

	// Read newest msg from udp connection in byteArray
	//		_, nodeAddr, err := conn.ReadFromUDP(byteArray)
	_, _, err := l.conn.ReadFromUDP(byteArray)
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

	//chanGotMsg <- currMsg
	//}
	return currMsg
}
