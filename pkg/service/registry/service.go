package registry

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/stainton/service-center/pkg/registry"
	"github.com/stainton/service-center/pkg/storage"
)

type RegistryService struct {
	pb.UnimplementedRegistryServer
	storageCli storage.RegistryStorage
}

func NewRegistryService(storageCli storage.RegistryStorage) (*RegistryService, error) {
	return &RegistryService{
		storageCli: storageCli,
	}, nil
}

func (s *RegistryService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ttl := req.GetTtl()
	svc := &pb.ServiceInstance{
		InstanceId: req.GetInstanceId(),
		Address:    req.GetAddress(),
		Metadata:   req.GetMetadata(),
		Port:       req.GetPort(),
		Ttl:        ttl,
	}
	leaseID, err := s.storageCli.Grant(ttl)
	if err != nil {
		return nil, err
	}
	if err = s.storageCli.Put(svc.InstanceId, svc, storage.WithLeaseID(leaseID)); err != nil {
		return nil, err
	}
	fmt.Printf("new server registered: %+v\n", svc)
	return &pb.RegisterResponse{
		Success: true,
		LeaseId: leaseID,
	}, nil
}

func (s *RegistryService) Deregister(ctx context.Context, req *pb.DeregisterRequest) (*pb.DeregisterResponse, error) {
	if err := s.storageCli.Delete(req.GetInstanceId()); err != nil {
		return nil, err
	}
	fmt.Printf("server deregister: %+v\n", req)
	return &pb.DeregisterResponse{
		Success: true,
	}, nil
}

func (s *RegistryService) Heartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	if err := s.storageCli.KeepAliveOnce(req.LeaseId); err != nil {
		return nil, err
	}
	fmt.Printf("receive heartbeat: %+v\n", req)
	return &pb.HeartbeatResponse{
		Success: true,
	}, nil
}

func (s *RegistryService) Watch(req *pb.WatchRequest, stream pb.Registry_WatchServer) error {
	return s.storageCli.Watch(req.GetServiceName(), func(a any) {
		buffer, ok := a.([]byte)
		if !ok {
			return
		}
		svc := &pb.ServiceInstance{}
		err := json.Unmarshal(buffer, svc)
		if err != nil {
			return
		}
		stream.Send(svc)
	}, storage.WithPrefix())
}
