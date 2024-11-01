package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
)

func (app *Application) History(userId int) ([]domain.Product, []domain.Product, error) {
	purchasedProducts, err := app.ServiceDB.RetrievePurchasedProduct(userId)
	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "History", "There is a problem of process of getting the  purchased products from the db", err)
	}

	likedProducts, err := app.ServiceDB.RetrieveLikedProduct(userId)

	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "History", "There is a problem of process of getting the  liked products from the db", err)
	}
	return purchasedProducts, likedProducts, nil
}
