package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/universeh2h/report/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) Login(c context.Context, req model.Login) (*model.User, error) {
	var name string
	query := `
    SELECT username
    FROM users 
    WHERE username = $1
`
	err := repo.db.QueryRowContext(c, query, req.Username).Scan(&name)
	if err != nil {
		fmt.Printf("errr : %s", err.Error())
		return nil, errors.New("invalid username")
	}

	return &model.User{
		Username: req.Username,
	}, nil
}
