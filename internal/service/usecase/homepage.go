package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
)

func (app *Application) Homepage(userId int) ([]domain.Product, error) {
	var products []domain.Product
	var err error

	if userId >= 0 {
		products, err = app.ServiceDB.RetrieveProducts("User", userId)
	} else {
		products, err = app.ServiceDB.RetrieveProducts("Guest", userId)
	}

	if err != nil {
		return nil, logger.ErrorWrapper("UseCase", "Homepage", "There is problem with retrieveing the data from the db", err)
	}

	return products, nil
}
