package internalgrpc

import (
	"net"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server      *grpc.Server
	application Application
}

type Application interface {
	GetLogger() logger.Logger
	GetStorage() storage.Storage
}

func NewServer(app Application) *GrpcServer {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor(app.GetLogger()),
		),
	)
	service := &Service{
		app: app,
	}

	pb.RegisterEventsServer(server, service)
	reflection.Register(server)

	return &GrpcServer{
		server:      server,
		application: app,
	}
}

func (s *GrpcServer) Start(address string) error {
	lsn, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.application.GetLogger().Info("grpc server started")

	return s.server.Serve(lsn)
}

func (s *GrpcServer) Stop() {
	s.server.GracefulStop()
}
