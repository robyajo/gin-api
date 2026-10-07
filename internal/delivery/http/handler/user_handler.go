package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gin-api/internal/model"
	"gin-api/internal/service"
)

// UserHandler mengelola permintaan HTTP terkait User
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler membuat instansiasi baru UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers mengambil daftar pengguna dengan pagination dan filter
func (h *UserHandler) GetUsers(c *gin.Context) {
	search := c.Query("search")
	roleID := c.Query("role_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userService.GetUsers(c.Request.Context(), search, roleID, page, limit)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pengguna", err.Error())
		return
	}

	PaginatedResponse(c, http.StatusOK, "Daftar pengguna berhasil diambil", users, page, limit, total)
}

// GetUserByID mengambil profil satu pengguna berdasarkan ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan", nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail pengguna", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Detail pengguna berhasil diambil", user)
}

// CreateUser membuat pengguna baru
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Validasi pembuatan pengguna gagal", err.Error())
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUserEmailExists) {
			ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrInvalidRole) {
			ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan pengguna", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, "Pengguna berhasil dibuat", user)
}

// UpdateUser memperbarui data pengguna
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Validasi pembaruan pengguna gagal", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan", nil)
			return
		}
		if errors.Is(err, service.ErrUserEmailExists) {
			ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrInvalidRole) {
			ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui pengguna", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Pengguna berhasil diperbarui", user)
}

// DeleteUser menghapus pengguna berdasarkan ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan", nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus pengguna", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Pengguna berhasil dihapus", nil)
}
