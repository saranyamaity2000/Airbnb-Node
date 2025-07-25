package router

import (
	config "AuthInGo/config/proxy"
	"AuthInGo/middlewares"
	"AuthInGo/utils"
	"fmt"

	"github.com/go-chi/chi/v5"
)

type ProxyRouter struct {
}

func (pr *ProxyRouter) Register(r chi.Router) {
	fmt.Println("Registering Proxies in our router")
	// chiRouter.HandleFunc("/fakestoreservice/*", utils.ProxyToService("https://fakestoreapi.in", "/fakestoreservice"))
	// r.With(middlewares.JWTAuthMiddleware).Get("/profile", ur.userController.GetUserById)
	for _, proxy := range config.AvailableProxyServers {
		fmt.Println(proxy)
		r.With(middlewares.JWTAuthMiddleware).HandleFunc(fmt.Sprintf("/%s/*", proxy.Alias), utils.ProxyToService(proxy.BaseURL, fmt.Sprintf("/%s", proxy.Alias)))
	}
}
