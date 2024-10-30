package usecase

import (
	"SimpleShop/internal/repository"
)

type Application struct {
	ServiceDB repository.DbModule
}

func NewUseCase(repoObject repository.DbModule) *Application {
	return &Application{
		ServiceDB: repoObject,
	}
}
