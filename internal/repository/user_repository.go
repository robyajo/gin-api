package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"gin-api/internal/model"
)

// UserRepository mendefinisikan interface akses data user
type UserRepository interface {
	FindAll(ctx context.Context, search, roleID string, limit, offset int) ([]model.User, int, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository membuat instansiasi baru UserRepository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context, search, roleID string, limit, offset int) ([]model.User, int, error) {
	var users []model.User
	var total int

	whereClauses := []string{"1=1"}
	args := []interface{}{}

	if search != "" {
		whereClauses = append(whereClauses, "(u.name LIKE ? OR u.email LIKE ?)")
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if roleID != "" {
		whereClauses = append(whereClauses, "u.role_id = ?")
		args = append(args, roleID)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Hitung total baris
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u WHERE %s", whereSQL)
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Ambil data users dengan pagination dan join role
	selectQuery := fmt.Sprintf(`
		SELECT u.id, u.role_id, u.name, u.email, u.password_hash, u.phone, u.is_active, u.created_at, u.updated_at,
		       r.name AS role_name, r.description AS role_description
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?`, whereSQL)

	selectArgs := append(args, limit, offset)
	if err := r.db.SelectContext(ctx, &users, selectQuery, selectArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	query := `SELECT u.id, u.role_id, u.name, u.email, u.password_hash, u.phone, u.is_active, u.created_at, u.updated_at,
	                 r.name AS role_name, r.description AS role_description
	          FROM users u
	          JOIN roles r ON u.role_id = r.id
	          WHERE u.id = ?`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT u.id, u.role_id, u.name, u.email, u.password_hash, u.phone, u.is_active, u.created_at, u.updated_at,
	                 r.name AS role_name, r.description AS role_description
	          FROM users u
	          JOIN roles r ON u.role_id = r.id
	          WHERE u.email = ?`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (id, role_id, name, email, password_hash, phone, is_active, created_at, updated_at)
			  VALUES (:id, :role_id, :name, :email, :password_hash, :phone, :is_active, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users 
			  SET role_id = :role_id, name = :name, email = :email, password_hash = :password_hash,
			      phone = :phone, is_active = :is_active, updated_at = :updated_at
			  WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
