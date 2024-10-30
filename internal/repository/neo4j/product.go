package neo4j

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (rp *Repository) RetrieveProducts(role string, userId int) ([]domain.Product, error) {
	var statement string
	mapping := make(map[string]interface{})
	if role == "User" {
		statement = `		
			MATCH (u:User {UserID: $userId}), (p:Product)
			WHERE NOT (u)-[:purchased]->(p)
			RETURN p
		`
		mapping = map[string]interface{}{"userId": userId}
	} else if role == "Guest" {
		statement = `MATCH (p : Product) RETURN p`
		mapping = nil
	}

	products, err := retrieveProductsOperation(rp, statement, mapping)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "RetrieveProducts", "There is a problem with executing the statement into db", err)
	}
	return products, nil
}

func retrieveProductsOperation(rp *Repository, statement string, mapping map[string]interface{}) ([]domain.Product, error) {
	session := rp.DB.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	// Execute the query with ExecuteRead
	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(context.Background(), statement, mapping)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "There is a problem with executing the statement into db", err)
		}

		var products []domain.Product
		for records.Next(context.Background()) {
			record := records.Record()
			productNode, _ := record.Get("p")
			if productNode != nil {
				node := productNode.(neo4j.Node)

				// Safely retrieve and type-assert each property with additional logging
				productID64, ok := node.Props["productID"].(int64)
				if !ok {
					productIDFloat, floatOk := node.Props["productID"].(float64) // Sometimes integers may be stored as float64
					if floatOk {
						productID64 = int64(productIDFloat)
					} else {
						return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Failed to retrieve 'productID' as int or float64", nil)
					}
				}
				productID := int(productID64)

				name, nameOk := node.Props["name"].(string)
				if !nameOk {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Failed to retrieve 'name' as string", nil)
				}

				description, descOk := node.Props["description"].(string)
				if !descOk {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Failed to retrieve 'description' as string", nil)
				}

				category, categoryOk := node.Props["category"].(string)
				if !categoryOk {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Failed to retrieve 'category' as string", nil)
				}

				cost64, costOk := node.Props["cost"].(int64)
				if !costOk {
					costFloat, costFloatOk := node.Props["cost"].(float64)
					if costFloatOk {
						cost64 = int64(costFloat)
					} else {
						return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Failed to retrieve 'cost' as int or float64", nil)
					}
				}
				cost := int(cost64)

				// Create the product and add it to the list
				product := domain.Product{
					ProductId:   productID,
					Name:        name,
					Description: description,
					Category:    category,
					Cost:        cost,
				}
				products = append(products, product)
			}
		}

		return products, nil
	})

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "There is a problem with session.ExecuteRead function", err)
	}

	return result.([]domain.Product), nil
}
