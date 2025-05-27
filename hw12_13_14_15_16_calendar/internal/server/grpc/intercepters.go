package internalgrpc

import (
	"context"
	"strings"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func loggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		remoteAddr := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			remoteAddr = p.Addr.String()
		}

		msg := strings.Builder{}
		msg.WriteString("gRPC call: ")
		msg.WriteString(remoteAddr + " ")
		msg.WriteString(info.FullMethod + " ")
		logger.Info(msg.String())

		return handler(ctx, req)
	}
}
