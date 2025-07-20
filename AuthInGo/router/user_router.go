package router

import (
	"AuthInGo/controllers"
	"AuthInGo/middlewares"
	utils "AuthInGo/utils"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	userController *controllers.UserController
}

func NewUserRouter(_userController *controllers.UserController) Router {
	return &UserRouter{
		userController: _userController,
	}
}

func (ur *UserRouter) Register(r chi.Router) {
	r.With(middlewares.AuthMiddleware, utils.RateLimitManagerInstance.UserRateLimitMiddleware, utils.RateLimitManagerInstance.IPRateLimitMiddleware).Get("/profile", ur.userController.GetUserById)
	r.With(utils.RateLimitManagerInstance.IPRateLimitMiddleware, middlewares.UserCreateRequestValidator).Post("/signup", ur.userController.CreateUser)
	r.With(utils.RateLimitManagerInstance.IPRateLimitMiddleware, middlewares.UserLoginRequestValidator).Post("/login", ur.userController.LoginUser)
}
