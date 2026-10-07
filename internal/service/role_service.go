package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gin-api/internal/model"
	"gin-api/internal/repository"
)

var (
	ErrRoleNotFound      = errors.New("role tidak ditemukan")
	ErrRoleNameExists    = errors.New("nama role sudah digunakan")
	ErrRoleSystemDefault = errors.New("role bawaan sistem (superadmin/admin/operator/user) tidak dapat dihapus")
	ErrRoleHasUsers      = errors.New("role tidak dapat dihapus karena masih digunakan oleh pengguna aktif")
)

// RoleService mendefinisikan business logic manajemen role
type RoleService interface {
	GetAllRoles(ctx context.Context) ([]model.Role, error)
	GetRoleByID(ctx context.Context, id string) (*model.Role, error)
	CreateRole(ctx context.Context, req model.CreateRoleRequest) (*model.Role, error)
	UpdateRole(ctx context.Context, id string, req model.UpdateRoleRequest) (*model.Role, error)
	DeleteRole(ctx context.Context, id string) error
}

type roleService struct {
	roleRepo repository.RoleRepository
}

// NewRoleService membuat instansiasi baru RoleService
func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return s.roleRepo.FindAll(ctx)
}

func (s *roleService) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *roleService) CreateRole(ctx context.Context, req model.CreateRoleRequest) (*model.Role, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(req.Name))

	// Cek apakah nama role sudah terpakai
	existing, err := s.roleRepo.FindByName(ctx, normalizedName)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleNameExists
	}

	now := time.Now()
	newRole := &model.Role{
		ID:          uuid.New().String(),
		Name:        normalizedName,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.roleRepo.Create(ctx, newRole); err != nil {
		return nil, err
	}

	return newRole, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id string, req model.UpdateRoleRequest) (*model.Role, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}

	normalizedName := strings.ToLower(strings.TrimSpace(req.Name))

	// Jika nama berubah, pastikan tidak konflik dengan role lain
	if normalizedName != role.Name {
		existing, err := s.roleRepo.FindByName(ctx, normalizedName)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrRoleNameExists
		}
	}

	role.Name = normalizedName
	role.Description = strings.TrimSpace(req.Description)
	role.UpdatedAt = time.Now()

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// Lindungi 4 role default sistem agar tidak dapat dihapus
	systemRoles := map[string]bool{
		"superadmin": true,
		"admin":      true,
		"operator":   true,
		"user":       true,
	}
	if systemRoles[role.Name] {
		return ErrRoleSystemDefault
	}

	// Cek apakah masih ada user yang terafiliasi dengan role ini
	count, err := s.roleRepo.CountUsersByRoleID(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal memverifikasi relasi pengguna: %w", err)
	}
	if count > 0 {
		return ErrRoleHasUsers
	}

	return s.roleRepo.Delete(ctx, id)
}
