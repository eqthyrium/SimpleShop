package usecase

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
)

func (app *Application) Homepage(userId int, searchingValue string) ([]domain.Product, map[string]bool, error) {
	var products []domain.Product
	var err error

	if userId >= 0 {
		products, err = app.ServiceDB.RetrieveProducts("User", userId)
	} else {
		products, err = app.ServiceDB.RetrieveProducts("Guest", userId)
	}

	if err != nil {
		return nil, nil, logger.ErrorWrapper("UseCase", "Homepage", "There is problem with retrieveing the data from the db", err)
	}

	if searchingValue == "" {
		return products, nil, nil
	} else {
		var filteredProduct []domain.Product
		var mapping map[string]bool = make(map[string]bool)
		for i := 0; i < len(products); i++ {
			mapping[products[i].Category] = true
			if products[i].Category == searchingValue {
				filteredProduct = append(filteredProduct, products[i])
			}
		}
		return filteredProduct, mapping, nil
	}

}
