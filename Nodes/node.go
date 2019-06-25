package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

type node struct {
	localAddr          *net.UDPAddr
	loggerAddr         *net.UDPAddr
	neighbourAddrs     []*net.UDPAddr
	conn               *net.UDPConn
	echoNodeAddr       *net.UDPAddr
	mem                int
	sumOfMem           int
	informed           bool
	initiator          bool
	neighboursInformed int
}

// NewNode is constructor for struct node
func newNode(args []string) *node {
	n := &node{}
	n.localAddr = BuildUDPAddr(os.Args[1])
	n.loggerAddr = BuildUDPAddr(os.Args[2])
	noOfNeighbours := len(os.Args) - 4
	n.neighbourAddrs = make([]*net.UDPAddr, noOfNeighbours)
	n.setNeighbours(os.Args[4:len(os.Args)], noOfNeighbours)
	n.setMem(os.Args[3])
	n.sumOfMem = 0
	n.informed = false
	n.initiator = false
	n.neighboursInformed = 0
	return n
}

// Setter for nodes field mem
func (n *node) setMem(s string) {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error string to int: ", err)
	}
	n.mem = i
}

// Setter for nodes field neighbourAddrs
func (n *node) setNeighbours(allAddrs []string, noOfNeighbours int) {
	for i := 0; i < noOfNeighbours; i++ {
		n.neighbourAddrs[i] = BuildUDPAddr(allAddrs[i])
	}
}

// Sends any type of msg to given addr
func (n *node) sendMsg(msg *Msg, addr *net.UDPAddr) {
	var buffer bytes.Buffer
	byteArray := make([]byte, 1024)

	//Encoding
	e := gob.NewEncoder(&buffer)
	err := e.Encode(msg)
	if err != nil {
		fmt.Println("Error at encode msg:", err)
	}

	// Cast buffer to byte[]
	byteArray, err = ioutil.ReadAll(&buffer)
	if err != nil {
		fmt.Println("Error at ReadAll:", err)
	}

	// Send msg
	_, err = n.conn.WriteToUDP(byteArray, addr)
	if err != nil {
		fmt.Println("Error at send:", err)
	}
}

// String sum is result got from peer in echo msg
// Convert string sum to int, add nodes mem, return both as string
func (n *node) setDataForEcho(sumFromPeer string) string {
	peersMem, err := strconv.Atoi(sumFromPeer)
	if err != nil {
		fmt.Println("Error at Atoi: ", err)
	}
	result := peersMem + n.mem
	return strconv.Itoa(result)
}

// SendLogMsg send the given msg to logger
func (n *node) sendLogMsg(msg Msg) {
	var logMsg Msg
	logMsg.SenderAddr = n.localAddr
	logMsg.MsgType = logging
	logInfo := msg.SenderAddr.String() + " send msg\n" +
		"to " + n.localAddr.String() +
		"\nwith type " + msg.getStringForType() +
		"\ndata: " + msg.Data
	logMsg.Data = logInfo
	n.sendMsg(&logMsg, n.loggerAddr)
}

func (n *node) sendInfoMsg(peerAddr *net.UDPAddr) {
	// Build info msg
	var infoMsg Msg
	infoMsg.SenderAddr = n.localAddr
	infoMsg.MsgType = info
	infoMsg.Data = "0"
	n.sendMsg(&infoMsg, peerAddr)
}

func (n *node) receiveStartMsg(msg Msg) {
	n.initiator = true
	n.informed = true
	for _, addr := range n.neighbourAddrs {
		n.sendInfoMsg(addr)
	}
}

func (n *node) receiveIEMsg(msg *Msg, addr *net.UDPAddr) {
	n.neighboursInformed++

	// Got the first info msg
	if n.informed == false {
		n.informed = true
		// Safe edge in spanning tree
		n.echoNodeAddr = addr
		// Inform all neighbours except echo node
		for _, neighbour := range n.neighbourAddrs {
			if neighbour.String() != addr.String() {
				n.sendInfoMsg(neighbour)
			}
		}
	}

	// If all neighbours are finished
	// send echo msg
	// to logger if node is initiator
	// or to echo node if node is not initiator
	if n.neighboursInformed == len(n.neighbourAddrs) {
		if n.initiator == true {
			var resultMsg Msg
			resultMsg.SenderAddr = n.localAddr
			resultMsg.MsgType = result
			currSum := n.sumOfMem + n.mem
			resultMsg.Data = strconv.Itoa(currSum)
			n.sendMsg(&resultMsg, n.loggerAddr)
		} else {
			var echoMsg Msg
			echoMsg.SenderAddr = n.localAddr
			echoMsg.MsgType = echo
			currSum := n.sumOfMem + n.mem
			echoMsg.Data = strconv.Itoa(currSum)
			n.sendMsg(&echoMsg, n.echoNodeAddr)
		}
	}
}
