package neo4j

import (
	"SimpleShop/pkg/logger"
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (rp *Repository) Purchase(userId, productId int) error {
	// Begin a Neo4j transaction
	session := rp.DB.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	// Define the query to create the "purchase" relationship
	query := `
		MATCH (u:User {id: $userId}), (p:Product {id: $productId})
		MERGE (u)-[:PURCHASED]->(p)
		RETURN u, p
	`

	// Run the query within a session
	_, err := session.Run(context.Background(), query, map[string]interface{}{
		"userId":    userId,
		"productId": productId,
	})

	// Return any error that occurred
	if err != nil {
		return logger.ErrorWrapper("Repository", "Purchase", "failed to create purchase relationship", err)
	}

	return nil
}

func (repo *Repository) Like(userId, productId int) error {
	return nil
}
