package server

import (
	"context"
	"net"

	pb "github.com/stainton/service-center/pkg/registry"
	"github.com/stainton/service-center/pkg/service/registry"
	"github.com/stainton/service-center/pkg/storage"
	"github.com/stainton/service-center/pkg/storage/etcd"

	"github.com/spf13/cobra"
	"github.com/stainton/logger"
	"github.com/stainton/logger/impls/zapimpl"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func etcdBackEnd(o *Options, l logger.Logger) (storage.RegistryStorage, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: o.Endpoints(),
	})
	if err != nil {
		l.Errorf("create etcd failed with error: %+v", err)
		return nil, err
	}
	backend := etcd.NewEtcdBackEnd(context.Background(), etcdClient)
	return backend, nil
}

func Run(o *Options, l logger.Logger) error {
	backend, err := etcdBackEnd(o, l)
	if err != nil {
		l.Errorf("create storage backend failed with error: %+v", err)
		return err
	}
	registryService, err := registry.NewRegistryService(backend)
	if err != nil {
		l.Errorf("new registry failed with error: %+v", err)
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRegistryServer(grpcServer, registryService)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   o.IP(),
		Port: o.Port(),
		Zone: o.Zone(),
	})
	if err != nil {
		l.Errorf("listen TCP failed with error: %+v", err)
		return err
	}
	return grpcServer.Serve(listener)
}

func NewCommand() *cobra.Command {
	opts := NewOptions()
	command := &cobra.Command{
		Use:   "service-center",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(cmd *cobra.Command, args []string) error {
			o := logger.NewOptions(
				logger.WithEncoding(logger.EncodingJSON),
				logger.WithMinLevel(logger.Debug),
				logger.WithOutputPaths([]string{opts.LogPath()}),
			)
			l, err := zapimpl.NewDevelopLogger(o)
			if err != nil {
				return err
			}
			if err := Run(opts, l); err != nil {
				l.Errorf("service-center exited with err: %+v", err)
				return err
			}
			return nil
		},
	}
	opts.AddFlags(command)
	return command
}
