package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"net/url"

	"time"

	"github.com/lstoll/grpce/h2c"
	"github.com/lstoll/grpce/helloproto"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Provide remote host address")
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing provided URL [%v]", err)
	}
	o := []grpc.DialOption{grpc.WithDialer(h2c.Dialer{URL: u}.DialGRPC)} // dialer to URL?
	if u.Scheme != "https" {
		o = append(o, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(u.Host, o...)
	if err != nil {
		log.Fatalf("Error connecting to remote host [%v]", err)
	}
	time.Sleep(1 * time.Second)
	c := helloproto.NewHelloClient(conn)
	resp, err := c.HelloWorld(context.Background(), &helloproto.HelloRequest{Name: "grpc-h2c client"})
	if err != nil {
		log.Fatalf("Error calling RPC: %v", err)
	}
	fmt.Printf("RPC call answered by %q and returned %q\n", resp.ServerName, resp.Message)
	conn.Close()
}
