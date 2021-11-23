package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/diwuwudi123/golang-mini-project/nameresolve/hello"
	"github.com/diwuwudi123/golang-mini-project/nameresolve/resolver"
	"google.golang.org/grpc"
)

var x int
var addr = "127.0.0.1"

type SayServ struct {
}

func main() {
	rand.Seed(time.Now().Unix())
	x = rand.Intn(10)

	port := rand.Intn(100)
	port += 8000
	addrs := net.JoinHostPort(addr, strconv.Itoa(port))

	lis, err := net.Listen("tcp", addrs)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(addrs)
	rpcServer := &SayServ{}
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	hello.RegisterSayServer(srv, rpcServer)

	err = resolver.RegisterEtcd(resolver.Schema, resolver.EtcdAddr, "127.0.0.1", port, "hello", 10)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fmt.Sprintf("x is %d", x), "rpc get_token init success")

	err = srv.Serve(lis)
	if err != nil {
		log.Println(err)
	}

}

func (s *SayServ) Hello(ctx context.Context, request *hello.HelloRequest) (*hello.HelloResponse, error) {
	fmt.Println(request.Name)
	return &hello.HelloResponse{Data: fmt.Sprintf("hello %s %d", request.Name, x)}, nil
}
