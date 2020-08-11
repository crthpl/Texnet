package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	x, y int
	used bool
}

type ItemStack struct {
	amnt int8		//can be at most 85 (legally) and at the least 1
	itype int32		//allows for four billion different types of items
	//nbt string	// unused (for now!)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func run() {
	contOut := make(chan bool)
	contIn := make(chan bool)
	pktReceive := make(chan string)
	//BEGIN NETWORKING CONNECTING
	fmt.Print("\n\n\nPlease Enter Server IP ----------------------------------------------------------->")
	host := ReadUserInput()
	if host == "l\n" { // if the user type l, it is localhost
		host = "localhost\n"
	}
	host = host[:len(host)-1]
	host = host + ":61435"
	c := StartConnection(host) //connect to the server
	//END NETWORKING CONNECTING
	cfg := pixelgl.WindowConfig{ //the settings for the window
		Title:  "Online Game",
		Bounds: pixel.R(0, 0, 640, 704),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	grass, err := loadPicture("grass.png")	//loading the grass tile
	tile, err := loadPicture("wood.png")	//loading the tile tile
	you, err := loadPicture("you.png")		//loading you
	inv, err := loadPicture("inv.png")		//loading the inventory hotbar slors
	if err != nil {
		panic(err)
	}
	grasses := pixel.NewBatch(&pixel.TrianglesData{}, grass)
	grassSpr := pixel.NewSprite(grass, grass.Bounds())
	tiles := pixel.NewBatch(&pixel.TrianglesData{}, tile)
	invTiles := pixel.NewBatch(&pixel.TrianglesData{}, tile)
	tileSpr := pixel.NewSprite(tile, tile.Bounds())
	yous := pixel.NewBatch(&pixel.TrianglesData{}, you)
	youSpr := pixel.NewSprite(you, you.Bounds())
	invs := pixel.NewBatch(&pixel.TrianglesData{}, inv)
	invSpr := pixel.NewSprite(inv, inv.Bounds())

	//BEGIN NETWORKING STUFF
	go func() {
		for {
			go func(c net.Conn) {
				for {
					incoming := RecievePackets(c)
					pktReceive <- incoming
				}
			}(c)

			go func() {
				for {
					outComing := ReadUserInput()
					SendChat(c, outComing)
					contOut <- true
				}
			}()
			for <-contIn {
			}
		}
	}()
	//END NETWORKING STUFF

	var (
		//playerPos = pixel.ZV
		frames = 0
		second = time.Tick(time.Second)
		selSlot = 0
	)
	var tilePos [21][21]int32
	var pls [10000]Player
	var youID int
	randAngles := [4]float64{0, 1.5708, 3.14159, 4.71239}
	rand.Seed(time.Now().UnixNano())
	var inventory [52]ItemStack	//in order: 10 hotbar slots, 30 inventory slots, 4 armor slots, 8 misceloanous bauble slots, 1 error slot
	inventory[1] = ItemStack{1, 1}
	grasses.Clear()
	for x := 16; x != 656; x += 32 {
		for y := 80; y != 720; y += 32 {
			grassSpr.Draw(grasses, pixel.IM.Moved(pixel.V(float64(x), float64(y))).ScaledXY(pixel.V(float64(x), float64(y)), pixel.V(2, 2)).Rotated(pixel.V(float64(x), float64(y)), randAngles[rand.Intn(3)]))
		}
	}
	grasses.Draw(win)

	//last := time.Now()
	for !win.Closed() {
		//dt := time.Since(last).Seconds()
		//last = time.Now()

		// placing/breaking
		if win.JustPressed(pixelgl.KeyUp) {
			SendPacket(c, "PLCB "+strconv.Itoa(pls[youID].x)+" "+strconv.Itoa(pls[youID].y+1)+"\n")
		}
		if win.JustPressed(pixelgl.KeyDown) {
			SendPacket(c, "PLCB "+strconv.Itoa(pls[youID].x)+" "+strconv.Itoa(pls[youID].y-1)+"\n")
		}
		if win.JustPressed(pixelgl.KeyRight) {
			SendPacket(c, "PLCB "+strconv.Itoa(pls[youID].x+1)+" "+strconv.Itoa(pls[youID].y)+"\n")
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			SendPacket(c, "PLCB "+strconv.Itoa(pls[youID].x-1)+" "+strconv.Itoa(pls[youID].y)+"\n")
		}
		//selecting slots
		if win.JustPressed(pixelgl.Key1) {
			selSlot = 0
		}
		if win.JustPressed(pixelgl.Key2) {
			selSlot = 1
		}
		if win.JustPressed(pixelgl.Key3) {
			selSlot = 2
		}
		if win.JustPressed(pixelgl.Key4) {
			selSlot = 3
		}
		if win.JustPressed(pixelgl.Key5) {
			selSlot = 4
		}
		if win.JustPressed(pixelgl.Key6) {
			selSlot = 5
		}
		if win.JustPressed(pixelgl.Key7) {
			selSlot = 6
		}
		if win.JustPressed(pixelgl.Key8) {
			selSlot = 7
		}
		if win.JustPressed(pixelgl.Key9) {
			selSlot = 8
		}
		if win.JustPressed(pixelgl.Key0) {
			selSlot = 9
		}
		//moving
		if win.JustPressed(pixelgl.KeyW) {
			if pls[youID].y!=19 {
				if tilePos[pls[youID].x][pls[youID].y+1] == 0 {
					SendPacket(c, "MOVE N " + strconv.Itoa(youID)+"\n")
				}
			}
		}
		if win.JustPressed(pixelgl.KeyA) {
			if pls[youID].x!=0 {
				if tilePos[pls[youID].x-1][pls[youID].y] == 0 {
					SendPacket(c, "MOVE W " + strconv.Itoa(youID)+"\n")
				}
			}
		}
		if win.JustPressed(pixelgl.KeyS) {
			if pls[youID].y!=0 {
				if tilePos[pls[youID].x][pls[youID].y-1] == 0 {
					SendPacket(c, "MOVE S " + strconv.Itoa(youID)+"\n")
				}
			}
		}
		if win.JustPressed(pixelgl.KeyD) {
			if pls[youID].x!=19 {
				if tilePos[pls[youID].x+1][pls[youID].y] == 0 {
					SendPacket(c, "MOVE E " + strconv.Itoa(youID)+"\n")
				}
			}
		}

		tiles.Clear()
		for x := 0; x != 21; x++ {
			for y := 0; y != 21; y++ {
				if tilePos[x][y] == 1 {
					tileSpr.Draw(tiles, pixel.IM.Moved(pixel.V(float64(x*32)+8, float64(y*32)+40)).ScaledXY(pixel.V(float64(x*32), float64(y*32)), pixel.V(2, 2)))
				}
			}
		}
		yous.Clear()
		for x := 0; x != 10000; x++ {
			if pls[x].used {
				youSpr.Draw(yous, pixel.IM.Moved(pixel.V(float64(pls[x].x*32)+8, float64(pls[x].y*32)+40)).ScaledXY(pixel.V(float64(pls[x].x*32), float64(pls[x].y*32)), pixel.V(2, 2)))
			}
		}
		invs.Clear()
		for x := 32; x != 672; x += 64 {
			if (x+32)/64==selSlot+1 {
				invSpr.Draw(invs, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4.4, 4.4)))
			} else {
				invSpr.Draw(invs, pixel.IM.Moved(pixel.V(float64(x), float64(32))).ScaledXY(pixel.V(float64(x), float64(32)), pixel.V(4, 4)))
			}
			switch inventory[selSlot].itype {
			case 0:
			
			case 1:
				tileSpr.Draw(invTiles, pixel.IM.Moved(pixel.V(float64(x), 32)).ScaledXY(pixel.V(float64(x), 32), pixel.V(2.3, 2.3)))
			}
		}

		//drawing everything
		win.Clear(colornames.Forestgreen)
		grasses.Draw(win)
		tiles.Draw(win)
		yous.Draw(win)
		invs.Draw(win)
		invTiles.Draw(win)
		win.Update()

		frames++ //Fps displaying stuff
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0 //End Fps displaying stuff
		case packet := <-pktReceive:
			pktType, pktData := DecodePacket(packet)
			fmt.Print(packet)
			switch pktType {
			case "CHAT":
				fmt.Print(strings.Join(pktData, " "))
			case "PLCB":
				x, _ := strconv.Atoi(pktData[0])
				y, _ := strconv.Atoi(RemoveNewline(pktData[1]))
				if x >= 20 {
					x = 20
				}
				if x <= 0 {
					x = 0
				}
				if y >= 20 {
					y = 20
				}
				if y <= 0 {
					y = 0
				}
				tileSpr.Draw(tiles, pixel.IM.Moved(pixel.V(float64(x*32), float64(y*32))).ScaledXY(pixel.V(float64(x*32), float64(y*32)), pixel.V(2, 2)))
				switch tilePos[x][y] {
				case 1:
					tilePos[x][y] = 0
				case 0:
					tilePos[x][y] = 1
				}
			case "MOVE":
				id, err := strconv.Atoi(RemoveNewline(pktData[1]))
				if err!=nil {
					panic(err)
				}
				switch pktData[0] {
				case "N":
					pls[id].y = pls[id].y + 1
				case "S":
					pls[id].y = pls[id].y - 1
				case "E":
					pls[id].x = pls[id].x + 1
				case "W":
					pls[id].x = pls[id].x - 1
				}
			case "SPWN":
				x, err := strconv.Atoi(pktData[0])
				y, err := strconv.Atoi(pktData[1])
				id, err := strconv.Atoi(RemoveNewline(pktData[2]))
				if err!=nil {
					panic(err)
				}
				pls[id] = Player{x, y, true}
			case "YOUI":
				youID, err = strconv.Atoi(RemoveNewline(pktData[0]))
				if err!=nil {
					panic(err)
				}
				if !pls[youID].used {
					pls[youID] = Player{10, 10, true}
				}
			case "DISC":
				id, err := strconv.Atoi(RemoveNewline(pktData[0]))
				if err!=nil {
					panic(err)
				}
				pls[id] = Player{0, 0, false}
			}
			
		default:
		}
	}
	os.Exit(12897)
}

func main() {
	pixelgl.Run(run) // all the graphics stuff (end everything else...)
}
