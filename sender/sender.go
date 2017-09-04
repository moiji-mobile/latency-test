package sender

import (
	"encoding/binary"
	"fmt"
	"net"
	"sort"
	"time"
	"github.com/moiji-mobile/latency-test/result"
)

type SentInfo struct {
	pkt	int
	sentAt	time.Time
}
type BySentPkt []SentInfo

type RecvInfo struct {
	pkt	int
	rcvdAt	time.Time
}
type ByRecvPkt []RecvInfo


// be able to sort by pkt info
func (a BySentPkt) Len() int { return len(a) }
func (a BySentPkt) Swap(i int, j int) { a[i], a[j] = a[j], a[i] }
func (a BySentPkt) Less(i int, j int) bool { return a[i].pkt < a[j].pkt }
func (a ByRecvPkt) Len() int { return len(a) }
func (a ByRecvPkt) Swap(i int, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRecvPkt) Less(i int, j int) bool { return a[i].pkt < a[j].pkt }

const (
	PKT_SIZE = 16
)

func send(id int, conn net.Conn, resChannel chan<-SentInfo) {
	b := make([]byte, PKT_SIZE)
	binary.BigEndian.PutUint16(b[0:], uint16(len(b) - 2))
	binary.BigEndian.PutUint32(b[2:], uint32(id))
	sentInfo := SentInfo{id, time.Now()}
	conn.Write(b)
	resChannel <- sentInfo
}

func recvAll(total int, conn net.Conn, resChannel chan<-RecvInfo) {
	b := make([]byte, PKT_SIZE)
	for i := 0; i < total; i++ {
		l, err := conn.Read(b)
		if l != PKT_SIZE {
			panic(fmt.Sprintf("Short read.. %v %v", l, err))
		}
		now := time.Now()
		id := int(binary.BigEndian.Uint32(b[2:]))
		resChannel <- RecvInfo{id, now}
	}
}

func buildResult(sent []SentInfo, recv []RecvInfo) result.Result {
	sort.Sort(BySentPkt(sent))
	sort.Sort(ByRecvPkt(recv))

	if len(sent) != len(recv) {
		panic(fmt.Sprintf("Not equal length: %v %v", len(sent), len(recv)))
	}

	items := make([]result.Item, 0, len(sent))
	for i := 0; i < len(sent); i++ {
		if sent[i].pkt != recv[i].pkt {
			panic("Not equal..")
		}
		elapsed := recv[i].rcvdAt.Sub(sent[i].sentAt)
		items = append(items, result.Item{sent[i].pkt, elapsed.Nanoseconds()})
	}
	return result.Result{items}
}

func Run(dest string, msgs int, delay time.Duration) (*result.Result, error) {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		fmt.Println("Failed to connect")
		return nil, err
	}

	sentChan := make(chan SentInfo)
	recvChan := make(chan RecvInfo)
	sentMsg := 0
	sent := make([]SentInfo, 0, msgs)
	recv := make([]RecvInfo, 0, msgs)
	tick := time.Tick(delay)

	// Start a receiver...
	go recvAll(msgs, conn, recvChan)

	// Sent and fetch whatever is there
	for sentMsg < msgs {
		select {
		case <- tick:
			go send(len(sent), conn, sentChan)
			sentMsg += 1
		case info := <-sentChan:
			sent = append(sent, info)
			fmt.Print("+")
		case info := <-recvChan:
			recv = append(recv, info)
			fmt.Print(".")
		}
	}

	// fetch pending results
	for len(recv) < sentMsg || len(sent) < sentMsg {
		select {
		case info := <-sentChan:
			sent = append(sent, info)
			fmt.Print("+")
		case info := <-recvChan:
			recv = append(recv, info)
			fmt.Print(".")
		}
	}
	fmt.Println("")

	conn.Close()
	res := buildResult(sent, recv)
	return &res, nil
}
