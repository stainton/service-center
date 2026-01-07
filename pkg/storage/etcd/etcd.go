package etcd

import (
	"context"
	"encoding/json"

	"github.com/stainton/service-center/pkg/storage"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdBackEnd struct {
	root       context.Context
	etcdClient *clientv3.Client
}

func NewEtcdBackEnd(ctx context.Context, client *clientv3.Client) *EtcdBackEnd {
	return &EtcdBackEnd{
		root:       ctx,
		etcdClient: client,
	}
}

func (e *EtcdBackEnd) getEtcdOptions(opts ...storage.CallOption) []clientv3.OpOption {
	options := &storage.CallOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var etcdOpts []clientv3.OpOption
	if options.Prefix {
		etcdOpts = append(etcdOpts, clientv3.WithPrefix())
	}
	if options.LeaseID > 0 {
		etcdOpts = append(etcdOpts, clientv3.WithLease(clientv3.LeaseID(options.LeaseID)))
	}
	return etcdOpts
}

func (e *EtcdBackEnd) Put(key string, val any, opts ...storage.CallOption) error {
	etcdOpts := e.getEtcdOptions(opts...)

	value, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = e.etcdClient.Put(e.root, key, string(value), etcdOpts...)
	return err
}

func (e *EtcdBackEnd) Grant(ttl int64, opts ...storage.CallOption) (int64, error) {
	resp, err := e.etcdClient.Grant(e.root, ttl)
	if err != nil {
		return 0, err
	}
	return int64(resp.ID), nil
}

func (e *EtcdBackEnd) Revoke(leaseID int64, opts ...storage.CallOption) error {
	_, err := e.etcdClient.Revoke(e.root, clientv3.LeaseID(leaseID))
	return err
}

func (e *EtcdBackEnd) Delete(key string, opts ...storage.CallOption) error {
	etcdOpts := e.getEtcdOptions(opts...)
	_, err := e.etcdClient.Delete(e.root, key, etcdOpts...)
	return err
}

func (e *EtcdBackEnd) KeepAliveOnce(leaseID int64, opts ...storage.CallOption) error {
	_, err := e.etcdClient.KeepAliveOnce(e.root, clientv3.LeaseID(leaseID))
	return err
}

func (e *EtcdBackEnd) Watch(key string, handler func(any), opts ...storage.CallOption) error {
	etcdOpts := e.getEtcdOptions(opts...)

	watchCh := e.etcdClient.Watch(e.root, key, etcdOpts...)
	for watchResp := range watchCh {
		for _, event := range watchResp.Events {
			handler(event.Kv.Value)
		}
	}
	return nil
}

var _ storage.RegistryStorage = (*EtcdBackEnd)(nil)
