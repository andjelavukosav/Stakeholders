package repo

import (
	"context"
	"database-example/model"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserRepository struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func NewUserRepository(logger *log.Logger) (*UserRepository, error) {
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	return &UserRepository{
		driver: driver,
		logger: logger,
	}, nil
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx,
			"CREATE (u:User {id:$id, username:$username, password:$password, email:$email, role:$role,  isBlocked: $isBlocked}) RETURN u",
			map[string]any{
				"id":        user.ID.String(),
				"username":  user.Username,
				"password":  user.Password,
				"email":     user.Email,
				"role":      user.Role,
				"isBlocked": false,
			})
	})
	if err != nil {
		repo.logger.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx,
			"MATCH (u:User {username:$username}) RETURN u.id, u.username, u.password, u.email, u.role,u.isBlocked",
			map[string]any{"username": username})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			record := result.Record()
			return &model.User{
				ID:        uuid.MustParse(record.Values[0].(string)),
				Username:  record.Values[1].(string),
				Password:  record.Values[2].(string),
				Email:     record.Values[3].(string),
				Role:      record.Values[4].(string),
				IsBlocked: record.Values[5].(bool),
			}, nil
		}

		// Ako ne postoji user
		return nil, errors.New("user not found")
	})
	if err != nil {
		return nil, err
	}

	user, ok := res.(*model.User)
	if !ok || user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r *UserRepository) DriverClose(ctx context.Context) error {
	if r.driver != nil {
		return r.driver.Close(ctx)
	}
	return nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "MATCH (u:User) RETURN u.id, u.username, u.email, u.role,u.isBlocked", nil)
		if err != nil {
			return nil, err
		}

		var users []*model.User
		for result.Next(ctx) {
			record := result.Record()
			users = append(users, &model.User{
				ID:        uuid.MustParse(record.Values[0].(string)),
				Username:  record.Values[1].(string),
				Email:     record.Values[2].(string),
				Role:      record.Values[3].(string),
				IsBlocked: record.Values[4].(bool),
				Password:  "",
			})
		}

		if err := result.Err(); err != nil {
			return nil, err
		}

		return users, nil
	})

	if err != nil {
		repo.logger.Println("Error getting all users:", err)
		return nil, err
	}

	users, ok := res.([]*model.User)
	if !ok {
		return nil, errors.New("failed to cast users result")
	}

	return users, nil
}

func (repo *UserRepository) SetUserBlocked(ctx context.Context, userID string, blocked bool) error {
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx,
			`MATCH (u:User {id:$id})
			 SET u.isBlocked = $blocked
			 RETURN u`,
			map[string]any{
				"id":      userID,
				"blocked": blocked,
			})
		return nil, err
	})

	if err != nil {
		repo.logger.Println("Error updating user block status:", err)
		return err
	}

	return nil

}

