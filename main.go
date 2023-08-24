package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	BroadCastPort = 4000
	SendSleep     = 3 * time.Second
)

var (
	remoteAddr = net.UDPAddr{}
	bindAddr   = net.UDPAddr{}
)

func init() {
	remoteAddr = net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: BroadCastPort,
	}
	bindAddr = net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: BroadCastPort,
	}
}

func getTemperature() string {
	file, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
	defer file.Close()
	if err != nil {
		return "0"
	}
	res, _ := io.ReadAll(file)
	tempInt, _ := strconv.Atoi(strings.TrimSpace(string(res)))
	tempFloat := float64(tempInt) / 1000
	return fmt.Sprintf("%.2f", tempFloat)
}

func Sender(wg *sync.WaitGroup) {
	// send broadcast
	defer wg.Done()
	conn, err := net.Dial("udp", remoteAddr.String())
	if err != nil {
		println(err.Error())
		return
	}
	defer conn.Close()
	for {
		hostName, err := os.Hostname()
		if err != nil {
			return
		}
		_, err = conn.Write([]byte(`{"host": "` + hostName + `, "time:"` + time.Now().Format("2006-01-02 15:04:05") + `", "from": "go"` + `", "temp": "` + getTemperature() + `"}`))
		if err != nil {
			println(err.Error())
			return
		}
		time.Sleep(SendSleep)
	}
}

func Receiver(wg *sync.WaitGroup) {
	// receive broadcast
	defer wg.Done()
	conn, err := net.ListenUDP("udp", &bindAddr)
	if err != nil {
		println(err.Error())
		return
	}
	for {
		var buf [1024]byte
		length, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			println(err.Error())
			return
		}
		fmt.Println("recv: ", string(buf[:length]), addr)
	}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go Sender(wg)
	go Receiver(wg)
	wg.Wait()
}

