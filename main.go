package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/venkor/ping-out-check/pinger"
)

var usage = `
ping-out-check - pings the target host and determines if it's reachable based on the amount of packet loss. Can write specific message to a file if provided or prints the message to standard output.
				 
Special thanks to Prometheus Community for maintaining pro-bing - https://github.com/prometheus-community/pro-bing

Usage:
    ping [--count] [--interval] [--timeout] [--ttl] [--size] [--privileged] [--message-positive] [--message-negative] [--skip-p] [--skip-n] [--filepath] [--stdout] [--packet-loss-ts] host
Examples:
	TBD
`

func main() {
	countPtr := flag.Int("count", -1, "number of echo requests to send")
	intervalPtr := flag.Duration("interval", time.Second, "interval betweeen each echo request")
	timeoutPtr := flag.Duration("timeout", time.Second*30, "ping for that amount of time")
	ttlPtr := flag.Int("ttl", 64, "TTL")
	sizePtr := flag.Int("size", 24, "send ICMP messages with certain payload")
	privilegedPtr := flag.Bool("privileged", false, "send a privileged raw ICMP ping")
	messagePositivePtr := flag.String("message-positive", "Host reachable", "message to use in output when host is considered reachable")
	messageNegativePtr := flag.String("message-negative", "Host unreachable", "message to use in output when host is considered unreachable")
	skipReachableMessageOut := flag.Bool("skip-p", false, "Skips writing the reachable message to output (either stdout or to filepath)")
	skipUnreachableMessageOut := flag.Bool("skip-n", false, "Skips writing the unreachable message to output (either stdout or to filepath)")
	stdoutPtr := flag.Bool("stdout", false, "if set to true, writes message to stdout instead of writing to file given in filepath")
	filePathPtr := flag.String("filepath", "", "filepath used while writing the file")
	packetLossThresholdPtr := flag.Float64("packet-loss-ts", 50.0, "consider target host as unreachable after losing certain percentage of packets")
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	target := &pinger.Target{
		Address:             flag.Arg(0),
		Count:               *countPtr,
		Interval:            *intervalPtr,
		Timeout:             *timeoutPtr,
		TTL:                 *ttlPtr,
		Size:                *sizePtr,
		Privileged:          *privilegedPtr,
		Stdout:              *stdoutPtr,
		Filepath:            *filePathPtr,
		PacketLossThreshold: *packetLossThresholdPtr,
		Message: pinger.MessageStruct{
			MessagePositive:           *messagePositivePtr,
			MessageNegative:           *messageNegativePtr,
			SkipReachableMessageOut:   *skipReachableMessageOut,
			SkipUnreachableMessageOut: *skipUnreachableMessageOut,
		},
	}
	err := target.Validate()
	if err != nil {
		fmt.Printf("error occured during flag validation: %s\n", err.Error())
		flag.Usage()
		return
	}
	target.SetExtendedDefaults()
	stats, err := target.PingIt()
	if err != nil {
		fmt.Printf("failed while pinging host: %s. Error: %s\n", target.Address, err.Error())
		return
	}
	reachable := target.IsTargetReachable(stats.PacketLoss)
	fmt.Printf("(%v%% packets lost vs. packet loss threshold %v%%). Target %s considered as reachable: %t \n", stats.PacketLoss, target.PacketLossThreshold, target.Address, reachable)
	err = target.Out(reachable)
	if err != nil {
		log.Fatal(err)
	}
}
