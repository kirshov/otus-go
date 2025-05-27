package internalgrpc

import (
	"context"
	"errors"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/domain"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedEventsServer
	app Application
}

func (s *Service) Create(ctx context.Context, req *pb.BaseEventRequest) (*pb.BaseEventResponse, error) {
	event, err := getEventFromRequest(req)
	if err != nil {
		return nil, err
	}

	id, err := s.app.GetStorage().Add(ctx, *event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BaseEventResponse{
		Response: &pb.BaseEventResponse_Id{
			Id: id,
		},
	}, nil
}

func (s *Service) Update(ctx context.Context, req *pb.BaseEventRequest) (*pb.BaseEventResponse, error) {
	event, err := getEventFromRequest(req)
	if err != nil {
		return nil, err
	}

	if event.ID == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is required")
	}

	err = s.app.GetStorage().Update(ctx, *event)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BaseEventResponse{
		Response: &pb.BaseEventResponse_Id{
			Id: event.ID,
		},
	}, nil
}

func (s *Service) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if req.DeleteEvent == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}

	id := req.DeleteEvent.Id
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "ID is required")
	}

	if err := s.app.GetStorage().Remove(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteResponse{
		Response: &pb.DeleteResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *Service) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	if req.ListEvent == nil {
		return nil, status.Error(codes.InvalidArgument, "event is required")
	}

	days := req.ListEvent.Days

	l, err := s.app.GetStorage().List(ctx, int(days))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Event, 0, len(l))
	for _, v := range l {
		result = append(result, &pb.Event{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			DateStart:   timestamppb.New(v.DateStart),
			DateEnd:     timestamppb.New(v.DateEnd),
		})
	}

	return &pb.ListResponse{
		Events: result,
	}, nil
}

func getEventFromRequest(req *pb.BaseEventRequest) (*domain.Event, error) {
	if req == nil {
		return nil, errors.New("request is empty")
	}

	reqEvent := req.Event
	if reqEvent == nil {
		return nil, status.Error(codes.InvalidArgument, "vote is not specified")
	}

	return &domain.Event{
		ID:          reqEvent.Id,
		Title:       reqEvent.Title,
		Description: reqEvent.Description,
		DateStart:   reqEvent.DateStart.AsTime(),
		DateEnd:     reqEvent.DateEnd.AsTime(),
	}, nil
}
