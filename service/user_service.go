package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	Route      string
	Controller controller.UserController
}

// CreateUser godoc
// @Summary Create new user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User data"
// @Success 201 {object} user.UserOut
// @Failure 400 {object} inout.BaseResponse
// @Router /api/v1/public/user/create [POST]
func (s UserService) CreateUser(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.CreateUser)
}

// GetUser godoc
// @Summary Get user by ID
// @Tags User
// @Description Get user details by ID
// @Accept json
// @Produce json
// @Success 200 {object} user.UserOut
// @Param id path string true "User ID"
// @Router /api/v1/private/user/details/{id} [GET]
func (s UserService) GetUser(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route+"/:id", s.Controller.GetUser)
}

// ListUsers godoc
// @Summary List all users
// @Tags User
// @Description Get paginated list of users
// @Accept json
// @Produce json
// @Success 200 {object} user.UserListOut
// @Param page query int false "Page number"
// @Router /api/v1/private/user/list [GET]
func (s UserService) ListUsers(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.ListUsers)
}

// UpdateUser godoc
// @Summary Update user
// @Tags User
// @Description Update user profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.UpdateUserRequest true "User data"
// @Success 200 {object} user.UserOut
// @Router /api/v1/private/user/update/{id} [PUT]
func (s UserService) UpdateUser(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route+"/:id", s.Controller.UpdateUser)
}

// DeleteUser godoc
// @Summary Delete user
// @Tags User
// @Description Delete user account
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/private/user/delete/{id} [DELETE]
func (s UserService) DeleteUser(r *gin.RouterGroup, route string) {
	r.DELETE("/"+s.Route+"/"+route+"/:id", s.Controller.DeleteUser)
}

// LoginUser godoc
// @Summary User login
// @Tags User
// @Description Authenticate user and return JWT token
// @Accept json
// @Produce json
// @Param user body user.LoginRequest true "Login credentials"
// @Success 200 {object} user.LoginOut
// @Router /api/v1/public/user/login [POST]
func (s UserService) LoginUser(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.LoginUser)
}

// UpdateUserStatus godoc
// @Summary Update user status
// @Tags User
// @Description Update user status (admin only)
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param status body user.UpdateStatusRequest true "Status data"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/private/user/status/{id} [PUT]
func (s UserService) UpdateUserStatus(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route+"/:id", s.Controller.UpdateUserStatus)
}