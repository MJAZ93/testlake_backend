package app

import (
	"testlake/controller"
	"testlake/service"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.RouterGroup) {
	userService := service.UserService{
		Route:      "user",
		Controller: controller.UserController{},
	}
	
	userService.CreateUser(r, "create")
	userService.LoginUser(r, "login")
}

func PrivateRoutes(r *gin.RouterGroup) {
	userService := service.UserService{
		Route:      "user",
		Controller: controller.UserController{},
	}
	
	userService.GetUser(r, "details")
	userService.ListUsers(r, "list")
	userService.UpdateUser(r, "update")
	userService.DeleteUser(r, "delete")
	userService.UpdateUserStatus(r, "status")
}