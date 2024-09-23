package postgresql

import (
	"context"
	"gohelp/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	user.Role = models.CustomerRole
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, email, password_hash, password, user_role) VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Email, user.PasswordHash, user.Password, user.Role)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, password_hash, user_role FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	return user, err
}
