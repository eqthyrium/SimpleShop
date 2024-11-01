package neo4j

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strconv"
)

func (rp *Repository) RetrieveProducts(role string, userId int) ([]domain.Product, error) {
	var statement string
	mapping := make(map[string]interface{})
	if role == "User" {
		statement = `		
			MATCH (u:User {UserID: $userId}), (p:Product)
			WHERE NOT (u)-[:PURCHASED]->(p)
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

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(context.Background(), statement, mapping)
		if err != nil {
			return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Error executing statement", err)
		}

		var products []domain.Product
		for records.Next(context.Background()) {
			record := records.Record()
			productNode, _ := record.Get("p")
			if productNode != nil {
				node := productNode.(neo4j.Node)

				// Access and convert properties safely
				productID, err := getIntProperty(node.Props, "productID")
				if err != nil {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", err.Error(), nil)
				}

				name, err := getStringProperty(node.Props, "name")
				if err != nil {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", err.Error(), nil)
				}

				category, err := getStringProperty(node.Props, "category")
				if err != nil {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", err.Error(), nil)
				}
				description, err := getStringProperty(node.Props, "description")
				if err != nil {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", err.Error(), nil)
				}

				cost, err := getIntProperty(node.Props, "cost")
				if err != nil {
					return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", err.Error(), nil)
				}

				product := domain.Product{
					ProductId:   productID,
					Name:        name,
					Category:    category,
					Description: description,
					Cost:        cost,
				}
				products = append(products, product)
			}
		}

		return products, nil
	})

	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "retrieveProductsOperation", "Error in session.ExecuteRead", err)
	}

	return result.([]domain.Product), nil
}

// Helper functions to safely get properties
func getIntProperty(props map[string]interface{}, key string) (int, error) {
	if val, ok := props[key]; ok {
		switch v := val.(type) {
		case int64:
			return int(v), nil
		case float64:
			return int(v), nil
		case string:
			parsedVal, err := strconv.Atoi(v)
			if err != nil {
				return 0, fmt.Errorf("Cannot convert %s to int", key)
			}
			return parsedVal, nil
		default:
			return 0, fmt.Errorf("Unknown type for %s", key)
		}
	}
	return 0, fmt.Errorf("Property %s not found", key)
}

func getStringProperty(props map[string]interface{}, key string) (string, error) {
	if val, ok := props[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal, nil
		}
		return "", fmt.Errorf("Property %s is not a string", key)
	}
	return "", fmt.Errorf("Property %s not found", key)
}

func (rp *Repository) RetrieveBehaviourBasedProduct(userId int) ([]domain.Product, error) {
	statement :=
		`MATCH (u:User {UserID: $userId})-[:LIKED|PURCHASED]->(p:Product)
		WITH u, collect(DISTINCT p.category) AS likedPurchasedCategories
		MATCH (otherProducts:Product)
		WHERE otherProducts.category IN likedPurchasedCategories AND NOT (u)-[:PURCHASED]->(otherProducts)
		RETURN otherProducts AS p
	`
	mapping := map[string]interface{}{"userId": userId}

	products, err := retrieveProductsOperation(rp, statement, mapping)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "RetrieveBehaviourBasedProduct", "There is a problem with executing the statement into db", err)
	}

	return products, nil
}

// Find similar users to user1 and recommend products they liked but user1 hasn't
func (rp *Repository) RetrieveCollaborativeProduct(userId int) ([]domain.Product, error) {
	fmt.Println("The RetrieveCollaborativeProduct function is started")
	statement := `
        MATCH (u1:User {UserID: $userId})-[:LIKED|PURCHASED]->(p:Product)<-[:LIKED|PURCHASED]-(u2:User)
        WITH u1, u2, COUNT(p) AS commonInteraction
        MATCH (u1)-[:LIKED|PURCHASED]->(p1:Product)
        WITH u1, u2, commonInteraction, COUNT(p1) AS totalU1Interaction
        MATCH (u2)-[:LIKED|PURCHASED]->(p2:Product)
        WITH u1, u2, commonInteraction, totalU1Interaction, COUNT(p2) AS totalU2Interaction
        WITH u1, u2, commonInteraction, totalU1Interaction, totalU2Interaction, 
             (1.0 * commonInteraction) / (totalU1Interaction + totalU2Interaction - commonInteraction) AS jaccardSimilarity
        WHERE jaccardSimilarity > 0.1  // Threshold to filter similar users
        WITH u1, u2, jaccardSimilarity

        // Subquery to recommend products based on similar users' likes
        CALL {
            WITH u1, u2, jaccardSimilarity
            MATCH (u2)-[:LIKED|PURCHASED]->(p2:Product)
            WHERE NOT (u1)-[:LIKED|PURCHASED]->(p2) // Exclude products user1 has already liked
            RETURN p2, (1.0 * jaccardSimilarity) AS recommendationScore
        }
        RETURN p2 AS p, SUM(recommendationScore) AS totalScore
        ORDER BY totalScore DESC;
    `
	mapping := map[string]interface{}{"userId": userId}

	products, err := retrieveProductsOperation(rp, statement, mapping)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "RetrieveCollaborativeProduct", "There is a problem with executing the statement into db", err)
	}

	fmt.Println("The collaborative products:", products)
	return products, nil
}

func (rp *Repository) RetrievePurchasedProduct(userId int) ([]domain.Product, error) {

	statement := `
		MATCH (u:User {UserID: $userId})-[:PURCHASED]->(p:Product)
		RETURN p
	`
	mapping := map[string]interface{}{"userId": userId}

	products, err := retrieveProductsOperation(rp, statement, mapping)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "RetrievePurchasedProduct", "There is a problem with executing the statement into db", err)
	}
	return products, nil

}

func (rp *Repository) RetrieveLikedProduct(userId int) ([]domain.Product, error) {

	statement := `
		MATCH (u:User {UserID: $userId})-[:LIKED]->(p:Product)
		RETURN p
	`
	mapping := map[string]interface{}{"userId": userId}

	products, err := retrieveProductsOperation(rp, statement, mapping)
	if err != nil {
		return nil, logger.ErrorWrapper("Repository", "RetrievePurchasedProduct", "There is a problem with executing the statement into db", err)
	}
	return products, nil

}
