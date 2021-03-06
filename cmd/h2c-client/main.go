package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"net/url"

	"net"

	"strconv"

	"github.com/lstoll/grpce/h2c"
	"github.com/lstoll/grpce/helloproto"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Provide remote host URL as first argument")
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing provided URL [%v]", err)
	}

	port := u.Port()
	if u.Port() == "" {
		p, err := net.LookupPort("tcp", u.Scheme)
		if err != nil {
			log.Fatalf("Error determining default port for scheme %s", u.Scheme)
		}
		port = strconv.Itoa(p)
	}
	conn, err := grpc.Dial(net.JoinHostPort(u.Hostname(), port),
		grpc.WithDialer(h2c.Dialer{URL: u}.DialGRPC),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Error connecting to remote host [%v]", err)
	}
	c := helloproto.NewHelloClient(conn)
	for i := 0; i < 10; i++ {
		resp, err := c.HelloWorld(context.Background(), &helloproto.HelloRequest{Name: "grpc-h2c client"})
		if err != nil {
			log.Fatalf("Error calling RPC: %v", err)
		}
		fmt.Printf("RPC call answered by %q and returned %q\n", resp.ServerName, resp.Message)
	}
	conn.Close()
}
