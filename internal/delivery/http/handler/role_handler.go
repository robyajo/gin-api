package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gin-api/internal/model"
	"gin-api/internal/service"
)

// RoleHandler mengelola permintaan HTTP terkait Role
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler membuat instansiasi baru RoleHandler
func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetRoles mengambil seluruh daftar role
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles(c.Request.Context())
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar role", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Daftar role berhasil diambil", roles)
}

// GetRoleByID mengambil detail satu role berdasarkan ID
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	role, err := h.roleService.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan", nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail role", err.Error())
		return
	}
	SuccessResponse(c, http.StatusOK, "Detail role berhasil diambil", role)
}

// CreateRole membuat role baru
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req model.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Validasi input role gagal", err.Error())
		return
	}

	role, err := h.roleService.CreateRole(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrRoleNameExists) {
			ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan role", err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, "Role berhasil ditambahkan", role)
}

// UpdateRole memperbarui data role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Validasi input update role gagal", err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan", nil)
			return
		}
		if errors.Is(err, service.ErrRoleNameExists) {
			ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui role", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Role berhasil diperbarui", role)
}

// DeleteRole menghapus role berdasarkan ID
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	err := h.roleService.DeleteRole(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan", nil)
			return
		}
		if errors.Is(err, service.ErrRoleSystemDefault) {
			ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		if errors.Is(err, service.ErrRoleHasUsers) {
			ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus role", err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Role berhasil dihapus", nil)
}
