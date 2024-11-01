package usecase

import "SimpleShop/pkg/logger"

func (app *Application) Purchase(userId, productId int) error {

	err := app.ServiceDB.Purchase(userId, productId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "Purchase", "There is a problem of process of creating purchase connection between User and Product node ", err)
	}
	return nil
}

func (app *Application) Like(userId, productId int) error {

	err := app.ServiceDB.Like(userId, productId)
	if err != nil {
		return logger.ErrorWrapper("UseCase", "Like", "There is a problem of process of creating liked connection between User and Product node ", err)
	}
	return nil

}
