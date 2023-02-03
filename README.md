# ping-out-check


## About:

Pinging the other side of wireguard tunnel that connects from home assistant to the main location. If the other side is not reachable, write something to a file which home assistant is watching. Later automation will pick-up that the file has changed and will restart the wireguard client.

## Synopsis:

```
> .\ping-out-check.exe -h 

ping-out-check - pings the target host and determines if it's reachable based on the amount of packet loss. Can write specific message to a file if provided or prints the message to standard output.

Special thanks to Prometheus Community for maintaining pro-bing - https://github.com/prometheus-community/pro-bing

Usage:
    ping [--count] [--interval] [--timeout] [--ttl] [--size] [--privileged] [--message-positive] [--message-negative] [--skip-p] [--skip-n] [--filepath] [--stdout] [--packet-loss-ts] host
Examples:
        TBD
  -count int
        number of echo requests to send (default -1)
  -filepath string
        filepath used while writing the file
  -interval duration
        interval betweeen each echo request (default 1s)
  -message-negative string
        message to use in output when host is considered unreachable (default "Host unreachable")
  -message-positive string
        message to use in output when host is considered reachable (default "Host reachable")
  -packet-loss-ts float
        consider target host as unreachable after losing certain percentage of packets (default 50)
  -privileged
        send a privileged raw ICMP ping
  -size int
        send ICMP messages with certain payload (default 24)
  -skip-n
        Skips writing the unreachable message to output (either stdout or to filepath)
  -skip-p
        Skips writing the reachable message to output (either stdout or to filepath)
  -stdout
        if set to true, writes message to stdout instead of writing to file given in filepath
  -timeout duration
        ping for that amount of time (default 30s)
  -ttl int
        TTL (default 64)
```