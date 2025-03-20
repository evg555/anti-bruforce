package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *Server) loggingMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	latency := time.Since(start)
	timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")

	msg := fmt.Sprintf("%s [%s] %s %s %d %d %s",
		getClientIP(ctx),
		timestamp,
		info.FullMethod,
		"HTTP/2",
		status.Code(err),
		latency.Milliseconds(),
		getUserAgent(ctx),
	)

	s.logger.Info(msg)

	return resp, err
}

func getClientIP(ctx context.Context) string {
	peerInfo, ok := peer.FromContext(ctx)
	if ok && peerInfo.Addr != nil {
		return peerInfo.Addr.String()
	}

	return "unknown IP address"
}

func getUserAgent(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgent, exists := md["user-agent"]; exists && len(userAgent) > 0 {
			return userAgent[0]
		}
	}
	return "unknown user agent"
}
