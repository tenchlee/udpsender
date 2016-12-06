package sender

import (
	"time"
	"github.com/tenchlee/udpsender/mflag"
	"fmt"
	"net"
	"os"
)

type UDPSender struct {
	Count int
	Duration time.Duration
	SizeBegin mflag.Bytes
	Step mflag.Bytes
	DstAddress *net.UDPAddr
	ch chan uint64
}

func NewUDPSender(dst string, count int, duration time.Duration,
					sizeBeg mflag.Bytes, step mflag.Bytes) (sender *UDPSender) {
	sender = new(UDPSender)
	sender.Count = count
	sender.Duration = duration
	sender.SizeBegin = sizeBeg
	sender.Step = step
	sender.ch = make(chan uint64)

	addr, err := net.ResolveUDPAddr("udp4", dst)
	if err != nil {
		fmt.Printf("dst: %s error: %s\n", dst, err.Error())
		os.Exit(1)
	}
	sender.DstAddress = addr
	return
}

func makeString(len mflag.Bytes) []byte {
	data := make([]byte, len)
	for i := uint64(0); i < uint64(len); i++ {
		data[i] = byte(i%26 + 97)
	}
	return data
}

func recv(conn *net.UDPConn, expect mflag.Bytes, ch chan bool) {
	data := make([]byte, expect)
	var recv mflag.Bytes
	conn.SetReadBuffer(int(expect))
	for expect != recv {
		b, err := conn.Read(data)
		if err != nil {
			fmt.Println("read fail:", err.Error())
			os.Exit(1)
		}
		recv += mflag.Bytes(b)
	}
	ch <- true
}

func (s *UDPSender) listen(conn *net.UDPConn) {
	ch := make(chan bool)
	for i:=0; i < s.Count; i++ {
		go recv(conn, mflag.Bytes(i)*s.Step+s.SizeBegin, ch)
		select {
		case <-ch:
		case <-time.After(s.Duration*2):
			fmt.Println("udp read timeout")
			os.Exit(1)
		}
		time.Sleep(s.Duration)
	}
	s.ch <- 1
}

func (s *UDPSender) Send() {
	conn, err := net.DialUDP("udp4", nil, s.DstAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("send begin...")
	go s.listen(conn)
	sndByte := s.SizeBegin
	for i:=0; i < s.Count; i++ {
		fmt.Printf("send %d\n", sndByte)
		_, err := conn.Write(makeString(sndByte))
		if err != nil {
			panic(err)
		}
		sndByte += s.Step
		time.Sleep(s.Duration)
	}
	<- s.ch
	conn.Close()
	fmt.Println("send finish")
}