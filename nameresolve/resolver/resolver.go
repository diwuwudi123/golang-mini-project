package resolver

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	cc          resolver.ClientConn
	etcdAddr    string
	schema      string
	serviceName string
	cli         *clientv3.Client
	conn        *grpc.ClientConn
}

var (
	nameResolver = make(map[string]*Resolver)
	etcdCli      *clientv3.Client
)

func init() {
	var err error
	etcdConfig := clientv3.Config{
		Endpoints: strings.Split(EtcdAddr, ","),
	}
	etcdCli, err = clientv3.New(etcdConfig)
	if err != nil {
		log.Fatalln(err)
	}
}
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	r.cc = cc
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := r.cli.Get(ctx, getName(r.schema, r.serviceName), clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(getName(r.schema, r.serviceName))
	var addList []resolver.Address
	for i := range resp.Kvs {
		log.Println(string(resp.Kvs[i].Value))
		addList = append(addList, resolver.Address{Addr: string(resp.Kvs[i].Value)})
	}
	r.cc.UpdateState(resolver.State{Addresses: addList})
	return r, nil
}
func (r *Resolver) Scheme() string {
	return ""
}
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {

}
func (r *Resolver) Close() {

}
func GetConn(schema, serviceName string) *grpc.ClientConn {
	name := getName(schema, serviceName)
	if r, ok := nameResolver[name]; ok {
		return r.conn
	}
	r := NewResolver(schema, serviceName)
	if r != nil {
		nameResolver[name] = r
	}
	return r.conn

}

func getName(schema, serviceName string) string {
	return fmt.Sprintf("%s/%s", schema, serviceName)
}
func NewResolver(schema, serviceName string) *Resolver {

	// var r Resolver
	r := &Resolver{
		etcdAddr:    EtcdAddr,
		schema:      schema,
		serviceName: serviceName,
		cli:         etcdCli,
	}
	resolver.Register(r)
	name := getName(schema, serviceName)
	conn, err := grpc.Dial(name,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Duration(5)*time.Second),
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	r.conn = conn
	return r
}
