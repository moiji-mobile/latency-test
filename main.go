package main

import (
	"flag"
	"fmt"
	"time"
	"github.com/moiji-mobile/latency-test/result"
	"github.com/moiji-mobile/latency-test/sender"
)

func main() {
	destFlag := flag.String("destination", "localhost:8000", "Destination to connect to")
	msgsFlag := flag.Int("messages", 5000, "Number of messages to send")
	timeFlag := flag.Duration("interval", 500*time.Millisecond, "Wait between messages")
	pktsFlag := flag.Int("packet-size", 16, "Bytes to send")
	flag.Parse()

	if (*pktsFlag < 2) {
		panic("packet-size need to be at least 6 bytes")
	}

	runResult, err := sender.Run(*destFlag, *msgsFlag, *timeFlag, *pktsFlag)
	if err != nil {
		fmt.Println("Testing failed!", err)
	} else {
		result.PrintResult(*runResult)
	}
}
