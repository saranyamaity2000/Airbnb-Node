package middlewares

import (
	"AuthInGo/dto"
	"AuthInGo/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

func UserLoginRequestValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.LoginUserRequestDTO

		// Read and decode the JSON body into the payload
		if err := utils.ReadJsonBody(r, &payload); err != nil {
			utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		// Validate the payload using the Validator instance
		if err := utils.Validator.Struct(payload); err != nil {
			utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		fmt.Println("Payload received for login:", payload)

		ctx := context.WithValue(r.Context(), "payload", payload) // Create a new context with the payload

		next.ServeHTTP(w, r.WithContext(ctx)) // Call the next handler in the chain
	})
}

func UserCreateRequestValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dto.CreateUserRequestDTO

		// Read and decode the JSON body into the payload
		if err := utils.ReadJsonBody(r, &payload); err != nil {
			utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		// Validate the payload using the Validator instance
		if err := utils.Validator.Struct(payload); err != nil {
			utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Validation failed", err)
			return
		}

		ctx := context.WithValue(r.Context(), "payload", payload)

		next.ServeHTTP(w, r.WithContext(ctx)) // Call the next handler in the chain
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AuthMiddleware: Checking for valid user identification")
		bearerToken := r.Header.Get("Authorization")
		fmt.Println("Bearer Token:", bearerToken)
		if !strings.HasPrefix(bearerToken, "Bearer ") {
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("Invalid token format"))
			return
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")
		if token == "" {
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Unauthorized", fmt.Errorf("Token is empty"))
			return
		}
		userClaims, err := utils.ValidateAndExtractUserClaims(token)
		if err != nil {
			fmt.Println(err.Error())
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		ctx := context.WithValue(r.Context(), "userClaims", userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
