package main

import (
	"context"
	"log"

	"github.com/diwuwudi123/golang-mini-project/nameresolve/hello"
	"github.com/diwuwudi123/golang-mini-project/nameresolve/resolver"
)

func main() {
	cli := resolver.GetConn("etcd", "hello")
	ctx := context.Background()
	client := hello.NewSayClient(cli)
	res, err := client.Hello(ctx, &hello.HelloRequest{Name: "wudi"})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res.Data)

	}
}
