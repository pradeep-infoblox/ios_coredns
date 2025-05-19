package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx interface{}, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.Dial("udp", "127.0.0.1:5353")
		},
	}

	ips, err := r.LookupHost(nil, "example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("IPs for example.com: %v\n", ips)
}
