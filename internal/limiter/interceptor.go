package limiter

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// достаём clientID (user-id из metadata или IP)
func getClientID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("user-id"); len(vals) > 0 {
			return vals[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		if addr, _, err := net.SplitHostPort(p.Addr.String()); err == nil {
			return addr
		}
		return p.Addr.String()
	}
	return "unknown"
}

func (l *Limiter) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	clientID := getClientID(ctx)

	if !l.Inc(clientID, info.FullMethod) {
		return nil, status.Errorf(codes.ResourceExhausted,
			"too many concurrent requests for client %s", clientID)
	}
	defer l.Dec(clientID, info.FullMethod)

	return handler(ctx, req)
}

func (l *Limiter) StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	clientID := getClientID(ss.Context())

	if !l.Inc(clientID, info.FullMethod) {
		return status.Errorf(codes.ResourceExhausted,
			"too many concurrent streams for client %s", clientID)
	}
	defer l.Dec(clientID, info.FullMethod)

	return handler(srv, ss)
}
