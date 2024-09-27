package postgresql

import (
	"context"
	"gohelp/internal/models"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, email, password_hash, password) VALUES ($1, $2, $3, $4)",
		user.Username, user.Email, user.PasswordHash, user.Password)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, password_hash, user_role, banned FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Banned)
	return &user, err
}

func (r *UserRepository) GetUserById(ctx context.Context, userID int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, "SELECT id, username, email, password_hash, user_role, banned FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Banned)
	return &user, err
}

func (r *UserRepository) ChangeBanStatus(ctx context.Context, userID int, status bool) error {
	log.Println(userID, status)
	_, err := r.db.Exec("UPDATE users SET banned = $1 WHERE id = $2", status, userID)
	return err
}
