package pinger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Target - holds pro-bing ("github.com/prometheus-community/pro-bing") parameters enriched with Message details and other.
type Target struct {
	Address, Filepath   string
	Count, TTL, Size    int
	Timeout, Interval   time.Duration
	Privileged, Stdout  bool
	PacketLossThreshold float64
	Message             MessageStruct
}

// MessageStruct - holds all relevant information about passed message details.
type MessageStruct struct {
	MessagePositive, MessageNegative                   string
	SkipReachableMessageOut, SkipUnreachableMessageOut bool
}

// Validate - validates passed values to flags.
func (t *Target) Validate() (err error) {
	if t.PacketLossThreshold > 100.0 || t.PacketLossThreshold < 0.0 {
		return fmt.Errorf("packet-loss-ts range is [0.0-100.0]")
	}
	return nil
}

// SetExtendedDefaults - extends default MessagePositive and MessageNegative with Address if no custom message value was specified.
func (t *Target) SetExtendedDefaults() {
	if t.Message.MessagePositive == "" {
		t.Message.MessagePositive = fmt.Sprintf("%s is reachable", t.Address)
	}
	if t.Message.MessageNegative == "" {
		t.Message.MessageNegative = fmt.Sprintf("%s is unreachable", t.Address)
	}
}

// PingIt - creates a NawPinger, sets various parameters for it and returns statistics gathered. Writes ping details to standard output.
func (t Target) PingIt() (stats *probing.Statistics, err error) {
	pinger, err := probing.NewPinger(t.Address)
	if err != nil {
		log.Printf("Failed to ping: %v, %s\n", t.Address, err.Error())
		return &probing.Statistics{}, err
	}
	// listen for ctrl-C signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pinger.Stop()
		}
	}()
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}
	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	pinger.Count = t.Count
	pinger.Size = t.Size
	pinger.Interval = t.Interval
	pinger.Timeout = t.Timeout
	pinger.TTL = t.TTL
	pinger.SetPrivileged(t.Privileged)
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		log.Printf("Failed to ping target host: %v, %s\n", t.Address, err.Error())
		return &probing.Statistics{}, err
	}
	return pinger.Statistics(), nil
}

// IsTargetReachable - returns true or false based if the PacketLossThreshold has been breached.
func (t Target) IsTargetReachable(packetLoss float64) bool {
	return packetLoss < t.PacketLossThreshold
}

func (t Target) pickMessage(reachable bool) string {
	if reachable {
		return t.Message.MessagePositive
	}
	return t.Message.MessageNegative
}

// Out - determines if and what message should be printed, either to standard output or to the file specified in the filepath.
func (t Target) Out(reachable bool) (err error) {
	if t.Message.SkipReachableMessageOut && t.Message.SkipUnreachableMessageOut {
		return nil
	}
	if reachable && t.Message.SkipReachableMessageOut {
		return nil
	}
	if !reachable && t.Message.SkipUnreachableMessageOut {
		return nil
	}
	message := t.pickMessage(reachable)
	if t.Stdout {
		fmt.Printf("%s\n", message)
		return
	}
	if t.Filepath != "" {
		if err := os.MkdirAll(filepath.Dir(t.Filepath), os.ModePerm); err != nil {
			return err
		}
		err = os.WriteFile(t.Filepath, []byte(message), 0660)
		if err != nil {
			return err
		}
	}
	return nil
}
