package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	SendPort = 4000
	RecvPort = 4001
)

var (
	localAddr  = net.UDPAddr{}
	remoteAddr = net.UDPAddr{}
	bindAddr   = net.UDPAddr{}
)

func init() {
	localAddr = net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: SendPort,
	}
	remoteAddr = net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: RecvPort,
	}
	bindAddr = net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: RecvPort,
	}
}

func Sender(wg *sync.WaitGroup) {
	// send broadcast
	defer wg.Done()
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
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
		_, err = conn.Write([]byte(`{"name":"` + hostName + `","port":` + fmt.Sprintf("%d", localAddr.Port) + `})`))
		if err != nil {
			println(err.Error())
			return
		}
		time.Sleep(time.Second * 3)
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
		var buf [128]byte
		length, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			println(err.Error())
			return
		}
		fmt.Println(string(buf[:length]), addr)
	}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go Sender(wg)
	go Receiver(wg)
	wg.Wait()
}
