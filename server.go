package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
	"math/rand"
	//"reflect"
)

type Player struct {
	x, y int
	id int
	used bool
}

func SendPacket(c net.Conn, packet string) {
	_, err := c.Write([]byte(packet))
	if err != nil {
		panic(err)
	}
}

func main() {
	clientCount := 0
	allClients := make(map[net.Conn]int)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	var pls [10002]Player
	packets := make(chan string)
	var tiles [21][21]int8
	server, err := net.Listen("tcp", ":61435")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnections:

			log.Printf("Accepted new client, #%d", clientCount)
			Tid:=rand.Intn(9990)
			allClients[conn] = Tid
			clientCount += 1

			time.Sleep(time.Millisecond*500)
			for i:=0;i!=10000;i++ {
				if pls[i].used {
					go SendPacket(conn, "SPWN "+strconv.Itoa(pls[i].x)+" "+strconv.Itoa(pls[i].y)+" "+strconv.Itoa(pls[i].id)+"\n")
					time.Sleep(time.Millisecond*50)
				}
			}
			
			pls[Tid] = Player{10, 10, Tid, true}
			paket := "SPWN 10 10 " + strconv.Itoa(Tid) + "\n"
			for conn, _ := range allClients {
				time.Sleep(time.Millisecond*30)
				go SendPacket(conn, paket)
			}
			
			time.Sleep(time.Millisecond*500)
			paket = "YOUI " + strconv.Itoa(Tid) + "\n"
			SendPacket(conn, paket)
			for x:=0;x!=21;x++ {
				for y:=0;y!=21;y++ {
					if tiles[x][y] == 1 {
						_, err := conn.Write([]byte(fmt.Sprintf("TILE %v %v\n", x, y)))
						if err != nil {
							deadConnections <- conn
						}
						time.Sleep(time.Millisecond*50)
					}
				}
			}
			go func(conn net.Conn, clientId int) {
				reader := bufio.NewReader(conn)
				for {
					incoming, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					packets <- fmt.Sprintf(incoming)
				}
				deadConnections <- conn

			}(conn, allClients[conn])
		case packet := <-packets:
			for conn, _ := range allClients {
				go func(conn net.Conn, packet string) {
					_, err := conn.Write([]byte(packet))
					if err != nil {
						deadConnections <- conn
					}
				}(conn, packet)
			}
			//log.Printf("New packet: %s", packet)
			pktSit := strings.Split(packet, " ")
			switch pktSit[0] {
			case "TILE":
				x, err:=strconv.Atoi(pktSit[1])
				y, err:=strconv.Atoi(pktSit[2][:len(pktSit[2])-1])
				if err!=nil {
					panic(err)
				}
				if x>=20 {
					x=20
				}
				if x<=0 {
					x=0
				}
				if y>=20 {
					y=20
				}
				if y<=0 {
					y=0
				}
				switch tiles[x][y] {
				case 1:
					tiles[x][y] = 0
				case 0:
					tiles[x][y] = 1
				}
			case "MOVE":
				id, err := strconv.Atoi(strings.Split(pktSit[2], "\n")[0])
				if err!=nil {
					panic(err)
				}
				switch pktSit[1] {
				case "N":
					pls[id].y = pls[id].y + 1
				case "S":
					pls[id].y = pls[id].y - 1
				case "E":
					pls[id].x = pls[id].x + 1
				case "W":
					pls[id].x = pls[id].x - 1
				}
			default:
			}
			//log.Printf("Broadcast to %d clients", len(allClients))
		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", allClients[conn])
			pls[allClients[conn]] = Player{0, 0, 9997, false}
			for conn, _ := range allClients {
				go SendPacket(conn, "DISC "+strconv.Itoa(allClients[conn])+"\n")
			}
			delete(allClients, conn)
		}
	}
	for conn, _ := range allClients {
		go func(conn net.Conn, packet string) {
			_, _ = conn.Write([]byte(packet))
		}(conn, "KICK")
	}
}