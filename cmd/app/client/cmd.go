package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stainton/service-center/cmd/app/cli"
	"github.com/stainton/service-center/pkg/registry"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct {
	target     string
	watch      string // 要监听的服务名称
	name       string // 服务名称
	instanceId string // 实例ID
	address    string // 服务地址
	port       int32  // 服务端口
	ttl        int64  // 心跳时间
}

func NewOptions() Options {
	return Options{}
}

func (o *Options) InstanceId() string {
	return o.instanceId
}

func (o *Options) Address() string {
	return o.address
}

func (o *Options) Port() int32 {
	return o.port
}

func (o *Options) AddFlags(command *cobra.Command) {
	command.Flags().StringVar(&o.target, "target", "localhost:50051", "service center target address")
	command.Flags().StringVar(&o.watch, "watch", "", "service name to watch")
	command.Flags().StringVar(&o.name, "name", "", "service name to register")
	command.Flags().Int64Var(&o.ttl, "ttl", 5, "service instance heartbeat ttl")
	command.Flags().StringVar(&o.instanceId, "instance-id", "", "service instance id")
	command.Flags().StringVar(&o.address, "address", "", "service instance address")
	command.Flags().Int32Var(&o.port, "port", 0, "service instance port")
}

func (o *Options) Target() string {
	return o.target
}

func (o *Options) Watch() string {
	return o.watch
}

func (o *Options) Name() string {
	return o.name
}

func (o *Options) Ttl() int64 {
	return o.ttl
}

func NewClient(o *Options) (registry.RegistryClient, error) {
	conn, err := grpc.NewClient(o.Target(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := registry.NewRegistryClient(conn)
	return c, nil
}

func ClientCommandTest(c registry.RegistryClient, o *Options) error {
	ctx := cli.RootContext()
	resp, err := c.Register(ctx, &registry.RegisterRequest{
		Name:       o.Name(),
		InstanceId: o.InstanceId(),
		Address:    o.Address(),
		Metadata:   map[string]string{},
		Port:       o.Port(),
		Ttl:        o.Ttl(),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Register response: %+v\n", resp)
	ch, err := c.Watch(ctx, &registry.WatchRequest{
		ServiceName: o.Watch(),
	})
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			watchResp, err := ch.Recv()
			if err != nil {
				fmt.Printf("Watch error: %v\n", err)
				return
			}
			fmt.Printf("Watch response: %+v\n", watchResp)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		tick := time.NewTicker(time.Duration(o.Ttl()/2) * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				_, err := c.Heartbeat(ctx, &registry.HeartbeatRequest{
					LeaseId: resp.LeaseId,
				})
				if err != nil {
					fmt.Printf("Heartbeat error: %v\n", err)
					return
				}
			}
		}
	}()
	wg.Wait()
	dregResp, err := c.Deregister(context.Background(), &registry.DeregisterRequest{
		Name:       o.Name(),
		InstanceId: o.InstanceId(),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Deregister response: %+v\n", dregResp)
	return nil
}
