package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gin-api/internal/model"
)

// RoleRepository mendefinisikan interface akses data role
type RoleRepository interface {
	FindAll(ctx context.Context) ([]model.Role, error)
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id string) error
	CountUsersByRoleID(ctx context.Context, roleID string) (int, error)
}

type roleRepository struct {
	db *sqlx.DB
}

// NewRoleRepository membuat instansiasi baru RoleRepository
func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindAll(ctx context.Context) ([]model.Role, error) {
	var roles []model.Role
	query := "SELECT id, name, description, created_at, updated_at FROM roles ORDER BY created_at ASC"
	err := r.db.SelectContext(ctx, &roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	return roles, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name, description, created_at, updated_at FROM roles WHERE id = ?"
	err := r.db.GetContext(ctx, &role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by id: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	query := "SELECT id, name, description, created_at, updated_at FROM roles WHERE name = ?"
	err := r.db.GetContext(ctx, &role, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find role by name: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	query := `INSERT INTO roles (id, name, description, created_at, updated_at)
			  VALUES (:id, :name, :description, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, role)
	if err != nil {
		return fmt.Errorf("failed to insert role: %w", err)
	}
	return nil
}

func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	query := `UPDATE roles 
			  SET name = :name, description = :description, updated_at = :updated_at
			  WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, role)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM roles WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("role not found")
	}
	return nil
}

func (r *roleRepository) CountUsersByRoleID(ctx context.Context, roleID string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE role_id = ?"
	err := r.db.GetContext(ctx, &count, query, roleID)
	if err != nil {
		return 0, fmt.Errorf("failed to count users by role: %w", err)
	}
	return count, nil
}
