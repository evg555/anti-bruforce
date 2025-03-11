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
	ipAddress := request.Ip

	if !common.IsValidIPAddress(ipAddress) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip address")
	}

	if h.app.IsInWhitelist(ctx, ipAddress) {
		h.logger.Info(fmt.Sprintf("ip address %s in whitelist", ipAddress))
		return &pb.Response{Ok: true}, nil
	}

	if h.app.IsInBlacklist(ctx, ipAddress) {
		h.logger.Info(fmt.Sprintf("ip address %s in blacklist", ipAddress))
		return &pb.Response{Ok: false}, nil
	}

	if !h.app.HasLimits(request.Login, request.Password, request.Ip) {
		h.logger.Info(fmt.Sprintf("limits are ended for login %s or password %s or ip %s",
			request.Login, request.Password, request.Ip))
		return &pb.Response{Ok: false}, nil
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) BucketReset(_ context.Context, request *pb.BucketResetRequest) (*pb.Response, error) {
	ipAddress := request.Ip

	if !common.IsValidIPAddress(ipAddress) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ip address")
	}

	h.app.ResetBucket(request.Password, ipAddress)
	h.logger.Info(fmt.Sprintf("reset bucket for password %s and ip %s", request.Password, ipAddress))

	return &pb.Response{Ok: true}, nil
}

func (h Handler) AddIPWhitelist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.AddIPWhitelist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("add ip whitelist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIPWhitelist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.DeleteIPWhitelist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("delete ip whitelist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) AddIPBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.AddIPBlacklist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("add ip blacklist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}

func (h Handler) DeleteIPBlacklist(ctx context.Context, request *pb.IpRequest) (*pb.Response, error) {
	subnet := request.NetAddress

	if !common.IsValidSubnet(subnet) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid net address")
	}

	err := h.app.DeleteIPBlacklist(ctx, subnet)
	if err != nil {
		h.logger.Error(fmt.Sprintf("delete ip blacklist failed: %v", err))
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &pb.Response{Ok: true}, nil
}
