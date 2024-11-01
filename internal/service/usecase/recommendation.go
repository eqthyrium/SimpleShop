package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"fmt"
)

func (app *Application) Recommendation(userId int) ([]domain.Product, []domain.Product, error) {
	fmt.Println("We entered to Recommendation with the UserId:", userId)
	behaviourProducts, err := app.ServiceDB.RetrieveBehaviourBasedProduct(userId)
	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "Recommendation", "There is a problem of process of getting the  behaviour based products of the given user from the db", err)
	}

	collaborativeProduct, err := app.ServiceDB.RetrieveCollaborativeProduct(userId)
	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "Recommendation", "There is a problem of process of getting the  collaborative filtering(user orientated) based products of the given user from the db", err)
	}

	return behaviourProducts, collaborativeProduct, nil

}
