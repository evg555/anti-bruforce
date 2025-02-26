package internalgrpc

import (
	"context"

	"github.com/evg555/antibrutforce/api/pb"
)

type Handler struct {
	app    Application
	logger Logger

	pb.UnimplementedAppServiceServer
}

func (h Handler) Auth(ctx context.Context, request *pb.AuthRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}

func (h Handler) BucketReset(ctx context.Context, request *pb.BucketResetRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}

func (h Handler) AddIpWhitelist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIpWhitelist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}

func (h Handler) AddIpBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIpBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	return &pb.Response{Ok: true}, nil
}
