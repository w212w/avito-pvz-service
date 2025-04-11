package repository

import (
	"avito-pvz-service/internal/models"
	"database/sql"
	"errors"
)

type UserRepoInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepoInterface {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, role, password_hash)
		VALUES ($1, $2, $3, $4)
		`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Role,
		user.PasswordHash,
	)

	if err != nil {
		return errors.New("error to create user")
	}
	return nil

}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	query := `
	SELECT id, email, role, password_hash
	FROM users
	WHERE email = $1
`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}
