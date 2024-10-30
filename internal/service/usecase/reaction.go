package usecase

import "SimpleShop/pkg/logger"

func (app *Application) Purchase(userId, productId int) error {
	err := app.Purchase(userId, productId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "Purchase", "There is a problem of process of creating purchase connection between User and Product node ", err)
	}
	return nil
}

func (app *Application) Like(userId, productId int) error {
	return nil
}
