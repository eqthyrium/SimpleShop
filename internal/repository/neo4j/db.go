package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repository struct {
	DB neo4j.DriverWithContext
}

func NewRepository(db neo4j.DriverWithContext) *Repository {
	return &Repository{DB: db}
}
