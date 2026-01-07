package server

import (
	"net"

	"github.com/spf13/cobra"
)

type Options struct {
	endPoints []string
	ip        string
	port      int
	zone      string
	logPath   string
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) LogPath() string {
	return o.logPath
}

func (o *Options) Endpoints() []string {
	return o.endPoints
}

func (o *Options) Zone() string {
	return o.zone
}

func (o *Options) Port() int {
	return o.port
}

func (o *Options) IP() net.IP {
	return net.ParseIP(o.ip)
}

func (o *Options) AddFlags(command *cobra.Command) {
	command.Flags().StringArrayVar(&o.endPoints, "end_point", []string{"127.0.0.1:2379"}, "etcd endpoints")
	command.Flags().StringVar(&o.ip, "ip", "127.0.0.1", "grpc server ip")
	command.Flags().StringVar(&o.zone, "zone", "", "grpc server IPv6 scoped addressing zone")
	command.Flags().IntVar(&o.port, "port", 50051, "grpc server port")
	command.Flags().StringVar(&o.logPath, "log", "", "path of log file")
}

func (o *Options) Compelete() error {
	return nil
}
