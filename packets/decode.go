package packets

import (
	"bufio"
	//"fmt"
	"math"
	//"net"
	"os"
	"strings"
)

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
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
