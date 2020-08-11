package main

import (
	"bufio"
	//"fmt"
	"math"
	"net"
	"os"
	"strings"
)

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func SendChat(c net.Conn, msg string) {
	SendPacket(c, "CHAT")
}

func ReadUserInput() (input string) {
	reader := bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')
	return input
}

func RemoveNewline(stringIn string) (stringOut string) {
	return stringIn[:len(stringIn)-1]
}

func RemoveNewlines(stringIn []string) (stringOut []string) {
	stringOu := stringIn[len(stringIn)-1]
	stringOu = stringOu[:len(stringOu)-1]
	stringIn = stringIn[:len(stringIn)-1]
	stringIn = append(stringIn, stringOu)
	return stringIn
}

func DecodePacket(packet string) (packetType string, packetData []string) {
	splitPacket := strings.Split(packet, " ")

	return splitPacket[0], splitPacket[1:]
}

func FormatPacket(packetType string, packetData []string) (packet string) {
	packetDataString := strings.Join(packetData, " ")

	return packetType + packetDataString
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
