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
	flag.Parse()

	runResult, err := sender.Run(*destFlag, *msgsFlag, *timeFlag)
	if err != nil {
		fmt.Println("Testing failed!", err)
	} else {
		result.PrintResult(*runResult)
	}
}
