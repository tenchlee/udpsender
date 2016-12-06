package main

import (
	"flag"
	"strings"
	"time"
	"fmt"
	"github.com/tenchlee/udpsender/mflag"
	"os"
	"github.com/tenchlee/udpsender/sender"
)


func main() {
	protocol := flag.String("proto", "udp", "send protocol")


	count := flag.Int("count", 10, "send count")
	dst := flag.String("dst", "", "send dst ip:port")
	duration := flag.Duration("duration", time.Second, "send duration time")
	sizeBeg := mflag.FlagBytes("size", 1, "send begin size (1kib=1024b, 1kb=1000b)")
	step := mflag.FlagBytes("step", 1, "send step size (1kib=1024b, 1kb=1000b)")
	flag.Parse()

	proto := strings.ToLower(*protocol)
	if len(*dst) < 2 {
		fmt.Printf("dst: %s is invalid\n", *dst)
		os.Exit(1)
	}

	var s sender.Sender
	switch proto {
	case "udp":
		s = sender.NewUDPSender(*dst, *count, *duration, *sizeBeg, *step)
	case "tcp":
		fmt.Println("not implement yet")
		os.Exit(1)
	}

	fmt.Printf("send \"%s\" to %s: %d times, duration: %s. begin: %s, step: %s \n\n", proto, *dst, *count, duration.String(),
		sizeBeg.String(), step.String())
	s.Send()
}
