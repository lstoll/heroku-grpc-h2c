package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"net"

	"github.com/lstoll/grpce/h2c"
	"github.com/lstoll/grpce/helloproto"
)

type hs struct{}

func (h *hs) HelloWorld(ctx context.Context, req *helloproto.HelloRequest) (*helloproto.HelloResponse, error) {
	return &helloproto.HelloResponse{
		Message:    fmt.Sprintf("Hello, %s!", req.Name),
		ServerName: "h2c",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	helloproto.RegisterHelloServer(s, &hs{})

	srv := &h2c.Server{
		HTTP2Handler:      s,
		NonUpgradeHandler: http.HandlerFunc(http.NotFound),
	}

	log.Printf("Serving at %s", lis.Addr().String())
	http.Serve(lis, srv)
}
