package grpc

import (
	"context"
	"time"

	pb "test_nanimai/backend/internal/api/grpc/pb"
	"test_nanimai/backend/internal/service"
)

type BalanceGRPCServer struct {
	pb.UnimplementedBalanceServiceServer
	svc service.Balance
}

func NewBalanceGRPCServer(svc service.Balance) *BalanceGRPCServer {
	return &BalanceGRPCServer{svc: svc}
}

func (s *BalanceGRPCServer) UpdateLimit(ctx context.Context, req *pb.UpdateLimitRequest) (*pb.Empty, error) {
	err := s.svc.UpdateLimit(ctx, req.AccountId, req.Delta)
	return &pb.Empty{}, err
}

func (s *BalanceGRPCServer) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.Empty, error) {
	err := s.svc.UpdateBalance(ctx, req.AccountId, req.Delta)
	return &pb.Empty{}, err
}

func (s *BalanceGRPCServer) OpenReservation(ctx context.Context, req *pb.OpenReservationRequest) (*pb.ReservationResponse, error) {
	res, err := s.svc.OpenReservation(
		ctx,
		req.OwnerServiceId,
		req.AccountId,
		req.Amount,
		req.IdempotencyKey,
		time.Duration(req.TimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	return &pb.ReservationResponse{
		ReservationId:  res.ID,
		AccountId:      res.AccountID,
		OwnerServiceId: res.OwnerServiceID,
		Amount:         res.Amount,
		Status:         res.Status,
		ExpiresAt:      res.ExpiresAt.Unix(),
	}, nil
}

func (s *BalanceGRPCServer) ConfirmReservation(ctx context.Context, req *pb.ReservationRequest) (*pb.Empty, error) {
	err := s.svc.ConfirmReservation(ctx, req.ReservationId, req.OwnerServiceId)
	return &pb.Empty{}, err
}

func (s *BalanceGRPCServer) CancelReservation(ctx context.Context, req *pb.ReservationRequest) (*pb.Empty, error) {
	err := s.svc.CancelReservation(ctx, req.ReservationId, req.OwnerServiceId)
	return &pb.Empty{}, err
}
