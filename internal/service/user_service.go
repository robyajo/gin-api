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
	"gin-api/pkg/utils"
)

var (
	ErrUserNotFound    = errors.New("pengguna tidak ditemukan")
	ErrUserEmailExists = errors.New("email sudah terdaftar dalam sistem")
	ErrInvalidRole     = errors.New("role yang dipilih tidak valid atau tidak terdaftar")
)

// UserService mendefinisikan business logic manajemen pengguna
type UserService interface {
	GetUsers(ctx context.Context, search, roleID string, page, limit int) ([]model.UserResponse, int, error)
	GetUserByID(ctx context.Context, id string) (*model.UserResponse, error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserService membuat instansiasi baru UserService
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *userService) GetUsers(ctx context.Context, search, roleID string, page, limit int) ([]model.UserResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	users, total, err := s.userRepo.FindAll(ctx, search, roleID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]model.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}

	return responses, total, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) CreateUser(ctx context.Context, req model.CreateUserRequest) (*model.UserResponse, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// 1. Verifikasi apakah email sudah terdaftar
	existing, err := s.userRepo.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserEmailExists
	}

	// 2. Verifikasi keberadaan role
	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrInvalidRole
	}

	// 3. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	now := time.Now()
	newUser := &model.User{
		ID:           uuid.New().String(),
		RoleID:       req.RoleID,
		Name:         strings.TrimSpace(req.Name),
		Email:        normalizedEmail,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		RoleName:     role.Name,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	resp := newUser.ToResponse()
	return &resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))

	// 1. Jika email berubah, periksa duplikasi
	if normalizedEmail != user.Email {
		existing, err := s.userRepo.FindByEmail(ctx, normalizedEmail)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrUserEmailExists
		}
	}

	// 2. Jika role berubah, validasi role baru
	if req.RoleID != user.RoleID {
		role, err := s.roleRepo.FindByID(ctx, req.RoleID)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, ErrInvalidRole
		}
		user.RoleName = role.Name
	}

	// 3. Update field
	user.Name = strings.TrimSpace(req.Name)
	user.Email = normalizedEmail
	user.RoleID = req.RoleID
	user.Phone = req.Phone
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// 4. Jika password baru disediakan, lakukan hash ulang
	if req.Password != nil && *req.Password != "" {
		newHash, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("gagal mengenkripsi kata sandi baru: %w", err)
		}
		user.PasswordHash = newHash
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}
