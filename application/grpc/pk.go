package grpc

import (
	"context"

	"github.com/codeedu/imersao/codepix-go/application/grpc/pb"
	"github.com/codeedu/imersao/codepix-go/application/usecase"
)

type PixKeyGrpcService struct {
	PixKeyUseCase usecase.PixKeyUseCase
	pb.UnimplementedPixKeyServiceServer
}

func (pixKeyGrpcService *PixKeyGrpcService) RegisterPixKey(ctx context.Context, in *pb.PixKeyRegistration) (*pb.PixKeyCreatedResult, error) {
	pixKey, err := pixKeyGrpcService.PixKeyUseCase.RegisterKey(in.Key, in.Kind, in.AccountId)

	if err != nil {
		return &pb.PixKeyCreatedResult{
			Status: "not created",
			Error:  err.Error(),
		}, err
	}

	return &pb.PixKeyCreatedResult{
		Id:     pixKey.ID,
		Status: "created",
	}, nil
}

func (pixKeyGrpcService *PixKeyGrpcService) Find(ctx context.Context, in *pb.PixKey) (*pb.PixKeyInfo, error) {
	pixKey, err := pixKeyGrpcService.PixKeyUseCase.FindKey(in.Key, in.Kind)

	if err != nil {
		return &pb.PixKeyInfo{}, err
	}

	return &pb.PixKeyInfo{
		Id:   pixKey.ID,
		Kind: pixKey.Kind,
		Key:  pixKey.Key,
		Account: &pb.Account{
			AccountId:     pixKey.AccountID,
			AccountNumber: pixKey.Account.Number,
			BankId:        pixKey.Account.BankID,
			BankName:      pixKey.Account.Bank.Name,
			OwnerName:     pixKey.Account.OwnerName,
			CreatedAt:     pixKey.Account.CreatedAt.String(),
		},
		CreatedAt: pixKey.Account.CreatedAt.String(),
	}, nil
}

func NewPixKeyGrpcService(usecase usecase.PixKeyUseCase) *PixKeyGrpcService {
	return &PixKeyGrpcService{
		PixKeyUseCase: usecase,
	}
}
