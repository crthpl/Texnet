package packets

import (
	"bufio"
	//"fmt"
	"math"
	"net"
	"os"
	"strings"
)

func SendChat(c net.Conn, msg string) {
	SendPacket(c, "CHAT"+msg)
}

func StartConnection(host string) net.Conn {
	c, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	return c
}

func RecievePackets(c net.Conn) (packet string) {
	reader := bufio.NewReader(c)
	incoming, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return incoming
}

func SendPacket(c net.Conn, packet string) {
	_, err := c.Write([]byte(packet))
	if err != nil {
		panic(err)
	}
}