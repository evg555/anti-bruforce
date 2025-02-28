package internalgrpc

import (
	"context"
	"fmt"

	"github.com/evg555/antibrutforce/api/pb"
	"github.com/evg555/antibrutforce/internal/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.AddIpWhitelist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("add ip whitelist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIpWhitelist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.DeleteIpWhitelist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("delete ip whitelist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) AddIpBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.AddIpBlacklist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("add ip blacklist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIpBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.DeleteIpBlacklist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("delete ip blacklist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}
