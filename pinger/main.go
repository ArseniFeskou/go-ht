package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sparrc/go-ping"
)

var usage = `
Usage:
    [-c count] [-t timeout] host
Examples:
    # ping google continuously
    www.google.com
    # ping google 5 times
    ping -c 5 -t 5 www.google.com
`

func main() {
	timeout := flag.Duration("t", time.Second*100000, "")
	count := flag.Int("c", -1, "")
	flag.Usage = func() {
		fmt.Printf(usage)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	i := 0
	for ok := true; ok; ok = flag.Arg(i) != "" {

		host := flag.Arg(i)

		pinger, err := ping.NewPinger(host)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}
		pinger.OnRecv = func(pkt *ping.Packet) {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
				pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}
		pinger.OnFinish = func(stats *ping.Statistics) {
			fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
				stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			fmt.Printf("round-trip min/avg/max = %v/%v/%v\n",
				stats.MinRtt, stats.AvgRtt, stats.MaxRtt)
		}
		i++
		pinger.Count = *count
		pinger.Timeout = *timeout
		pinger.SetPrivileged(true)

		fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
		pinger.Run()
	}
}
