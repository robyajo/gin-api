package model

import "time"

// User merepresentasikan entitas akun pengguna dalam database
type User struct {
	ID           string    `json:"id" db:"id"`
	RoleID       string    `json:"role_id" db:"role_id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Phone        *string   `json:"phone,omitempty" db:"phone"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Relasi Role jika di-JOIN
	RoleName        string  `json:"role_name,omitempty" db:"role_name"`
	RoleDescription *string `json:"role_description,omitempty" db:"role_description"`
}

// UserResponse DTO aman yang dikembalikan ke client tanpa password hash
type UserResponse struct {
	ID        string    `json:"id"`
	RoleID    string    `json:"role_id"`
	RoleName  string    `json:"role_name"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     *string   `json:"phone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse mengonversi User entity menjadi UserResponse DTO
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		RoleID:    u.RoleID,
		RoleName:  u.RoleName,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CreateUserRequest payload pendaftaran/pembuatan user baru
type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Email    string  `json:"email" binding:"required,email,max=150"`
	Password string  `json:"password" binding:"required,min=6,max=100"`
	RoleID   string  `json:"role_id" binding:"required,uuid"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,max=20"`
}

// UpdateUserRequest payload pembaruan data user
type UpdateUserRequest struct {
	Name     string  `json:"name" binding:"required,min=2,max=100"`
	Email    string  `json:"email" binding:"required,email,max=150"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6,max=100"`
	RoleID   string  `json:"role_id" binding:"required,uuid"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,max=20"`
	IsActive *bool   `json:"is_active,omitempty"`
}
