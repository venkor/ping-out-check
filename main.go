package main

import (
	"flag"
	"fmt"
)

func main() {
	msgPtr := flag.String("message", "pinger", "message to use in output")
	filePathPtr := flag.String("filepath", "pinger_output.txt", "filepath used while writing the file")
	addressPtr := flag.String("address", "127.0.0.1", "address to ping")
	countPtr := flag.Int("count", 4, "number of echo requests to send")
	deadlinePtr := flag.Int("deadline", 5, "number of (seconds-Linux/miliseconds-Windows) to wait for each reply")
	stdoutPtr := flag.Bool("stdout", false, "if set to true, writes message to stdout instead of writing to file given in filepath")

	flag.Parse()

	fmt.Println("message:", *msgPtr)
	fmt.Println("filepath:", *filePathPtr)
	fmt.Println("address:", *addressPtr)
	fmt.Println("count:", *countPtr)
	fmt.Println("deadline:", *deadlinePtr)
	fmt.Println("stdout:", *stdoutPtr)

	target := pinger.Target{}

	fmt.Println(target)

}
