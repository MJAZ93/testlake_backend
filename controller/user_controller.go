package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"testlake/dao"
	"testlake/inout"
	"testlake/inout/user"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct{}

func (controller UserController) CreateUser(context *gin.Context) {
	var request user.CreateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()

	emailExists, err := userDao.EmailExists(request.Email)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if emailExists {
		utils.ReportBadRequest(context, "Email already exists")
		return
	}

	usernameExists, err := userDao.UsernameExists(request.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if usernameExists {
		utils.ReportBadRequest(context, "Username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to process password")
		return
	}

	newUser := &model.User{
		Email:        request.Email,
		Username:     request.Username,
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		AuthProvider: request.AuthProvider,
		PasswordHash: &hashedPassword,
		Status:       model.UserStatusActive,
	}

	if err := userDao.Create(newUser); err != nil {
		utils.ReportInternalServerError(context, "Failed to create user")
		return
	}

	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.FromModel(newUser),
	}

	context.JSON(http.StatusCreated, response)
}

func (controller UserController) GetUser(context *gin.Context) {
	idParam := context.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID format")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.FromModel(foundUser),
	}

	context.JSON(http.StatusOK, response)
}

func (controller UserController) ListUsers(context *gin.Context) {
	pageStr := context.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		utils.ReportBadRequest(context, "Invalid page number")
		return
	}

	userDao := dao.NewUserDao()
	users, total, err := userDao.GetAll(page)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(userDao.Limit)))

	response := user.UserListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: user.FromModelList(users),
		Meta: inout.PaginationMeta{
			Page:       page,
			Limit:      userDao.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	context.JSON(http.StatusOK, response)
}

func (controller UserController) UpdateUser(context *gin.Context) {
	idParam := context.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID format")
		return
	}

	currentUserID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if currentUserID != id {
		utils.ReportForbidden(context, "Cannot update other user's profile")
		return
	}

	var request user.UpdateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	existingUser, err := userDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	existingUser.FirstName = request.FirstName
	existingUser.LastName = request.LastName
	existingUser.AvatarURL = request.AvatarURL

	if err := userDao.Update(existingUser); err != nil {
		utils.ReportInternalServerError(context, "Failed to update user")
		return
	}

	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.FromModel(existingUser),
	}

	context.JSON(http.StatusOK, response)
}

func (controller UserController) DeleteUser(context *gin.Context) {
	idParam := context.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID format")
		return
	}

	currentUserID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if currentUserID != id {
		utils.ReportForbidden(context, "Cannot delete other user's account")
		return
	}

	userDao := dao.NewUserDao()
	_, err = userDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	if err := userDao.Delete(id); err != nil {
		utils.ReportInternalServerError(context, "Failed to delete user")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "User deleted successfully",
	}

	context.JSON(http.StatusOK, response)
}

func (controller UserController) LoginUser(context *gin.Context) {
	var request user.LoginRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportUnauthorized(context, "Invalid credentials")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	if foundUser.PasswordHash == nil || !utils.CheckPasswordHash(request.Password, *foundUser.PasswordHash) {
		utils.ReportUnauthorized(context, "Invalid credentials")
		return
	}

	if foundUser.Status != model.UserStatusActive {
		utils.ReportForbidden(context, "Account is not active")
		return
	}

	token, err := utils.GenerateJWT(foundUser.ID, foundUser.Email, foundUser.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to generate token")
		return
	}

	userDao.UpdateLastLogin(foundUser.ID)

	response := user.LoginOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
	}
	response.Data.User = user.FromModel(foundUser)
	response.Data.Token = token

	context.JSON(http.StatusOK, response)
}

func (controller UserController) UpdateUserStatus(context *gin.Context) {
	idParam := context.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID format")
		return
	}

	var request user.UpdateStatusRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	_, err = userDao.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	if err := userDao.UpdateStatus(id, request.Status); err != nil {
		utils.ReportInternalServerError(context, "Failed to update user status")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "User status updated successfully",
	}

	context.JSON(http.StatusOK, response)
}