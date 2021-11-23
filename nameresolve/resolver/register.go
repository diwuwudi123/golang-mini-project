package resolver

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func RegisterEtcd(schema, etcdAddr, host string, port int, serviceName string, ttl int) error {
	log.Println("RegisterEtcd")

	ctx, _ := context.WithCancel(context.Background())
	resp, err := etcdCli.Grant(ctx, int64(ttl))
	if err != nil {
		log.Println(err)
		return err
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	name := getName(schema, serviceName) + addr
	if _, err := etcdCli.Put(ctx, name, addr, clientv3.WithLease(resp.ID)); err != nil {
		log.Println(err)
		return err
	} else {
		res, err := etcdCli.Get(ctx, name)
		log.Println(res.Kvs, err)
	}
	//自动续期
	var kresp <-chan *clientv3.LeaseKeepAliveResponse
	kresp, err = etcdCli.KeepAlive(ctx, resp.ID)
	if err != nil {
		log.Println(err)
	}
	go func() {
	FLOOP:
		for {
			select {
			case data, ok := <-kresp:
				if ok == true {
					fmt.Println("data is ", data)
				} else {
					break FLOOP
				}
			}
		}
	}()

	go Watch(name)
	return nil
}
func Watch(name string) {
	for range time.Tick(time.Second * 5) {
		ctx, _ := context.WithCancel(context.Background())
		res, err := etcdCli.Get(ctx, name)
		log.Println(res.Kvs, err)
	}

}
