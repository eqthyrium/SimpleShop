package neo4j

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (rp *Repository) CreateUser(user *domain.User) error {

	session := rp.DB.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.Background())

	statement := `
        CREATE (u:User {
			UserID : $userID,
            Nickname: $nickname,
            MemberIdentity: $memberIdentity,
            Password: $password,
            Role: $role
        })
    `

	// Prepare parameters to be injected into the Cypher query
	params := map[string]interface{}{
		"userID":         user.UserId,
		"nickname":       user.Nickname,
		"memberIdentity": user.MemberIdentity,
		"password":       user.Password,
		"role":           user.Role,
	}

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(context.Background(), statement, params)
		return nil, err
	})

	if err != nil {
		return logger.ErrorWrapper("Repository", "CreateUser", "The problem within the process of creation of the user in Neo4j", err)
	}

	return nil
}

func (rp *Repository) GetUserByEmail(memberIdentity string) (*domain.User, error) {

	session := rp.DB.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	statement := `
        MATCH (u:User { MemberIdentity: $memberIdentity })
        RETURN u.UserID AS UserID, u.MemberIdentity AS MemberIdentity, u.Password AS Password, u.Role AS Role
    `
	params := map[string]interface{}{
		"memberIdentity": memberIdentity,
	}

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		res, err := tx.Run(context.Background(), statement, params)
		if err != nil {
			return nil, err
		}

		if res.Next(context.Background()) {
			record := res.Record()
			user := &domain.User{}

			if value, ok := record.Get("UserID"); ok && value != nil {
				user.UserId = int(value.(int64)) // Adjust type if necessary
			}
			if value, ok := record.Get("MemberIdentity"); ok && value != nil {
				user.MemberIdentity = value.(string)
			}
			if value, ok := record.Get("Password"); ok && value != nil {
				user.Password = value.(string)
			}
			if value, ok := record.Get("Role"); ok && value != nil {
				user.Role = value.(string)
			}

			return user, nil
		}

		if err = res.Err(); err != nil {
			return nil, err
		}

		return nil, domain.ErrNoRecord
	})

	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "There is no such user in the db", domain.ErrUserNotFound)
		} else {
			return nil, logger.ErrorWrapper("Repository", "GetUserByEmail", "Problem during user retrieval from the database", err)
		}
	}

	return result.(*domain.User), nil
}

func (rp *Repository) GetLastUserId() (int, error) {
	session := rp.DB.NewSession(context.Background(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(context.Background())

	// Cypher query to get the maximum UserID
	statement := `
        MATCH (u:User)
        RETURN COALESCE(MAX(u.UserID), 0) AS maxUserID
    `

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		records, err := tx.Run(context.Background(), statement, nil)
		if err != nil {
			return nil, err
		}

		if records.Next(context.Background()) {
			record := records.Record()
			maxUserIDInterface, found := record.Get("maxUserID")
			if !found || maxUserIDInterface == nil {
				return 0, nil
			}

			// Type assertion to int64, which is the Neo4j integer type in Go
			maxUserID, ok := maxUserIDInterface.(int64)
			if !ok {
				return 0, fmt.Errorf("maxUserID is not of type int64")
			}
			return int(maxUserID), nil
		}

		return 0, records.Err()
	})

	if err != nil {
		return 0, logger.ErrorWrapper("Repository", "GetLastUserId", "Error getting last UserID from Neo4j", err)
	}

	return result.(int), nil
}
