package usecase

import (
	"github.com/codeedu/imersao/codepix-go/domain/model"
)

type PixKeyUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
}

func (pixKeyUseCase *PixKeyUseCase) RegisterKey(key string, kind string, accountId string) (*model.PixKey, error) {
	account, err := pixKeyUseCase.PixKeyRepository.FindAccount(accountId)

	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)

	if err != nil {
		return nil, err
	}

	pixKeyUseCase.PixKeyRepository.RegisterKey(pixKey)

	if pixKey.ID == "" {
		return nil, err
	}

	return pixKey, nil
}

func (pixKeyUseCase *PixKeyUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	pixKey, err := pixKeyUseCase.PixKeyRepository.FindKeyByKind(key, kind)

	if err != nil {
		return nil, err
	}

	return pixKey, nil
}
