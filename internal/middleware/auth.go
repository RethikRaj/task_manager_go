package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RethikRaj/task_manager_go/internal/ctx"
	"github.com/RethikRaj/task_manager_go/internal/handler"
	"github.com/RethikRaj/task_manager_go/internal/service"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get Authorization header (Format: Bearer <token>)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errResp := handler.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "Missing Authorization Header",
					Code:    "MISSING_AUTHORIZATION_HEADER",
					Success: false,
				}

				handler.SendJSONResponse(w, errResp.Status, errResp)
				return
			}

			// Split authHeader to get Token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errResp := handler.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "Invalid Token format",
					Code:    "INVALID_TOKEN_FORMAT",
					Success: false,
				}

				handler.SendJSONResponse(w, errResp.Status, errResp)
				return
			}

			// Verify token
			token := parts[1]

			claims, err := service.VerifyToken(token, jwtSecret)

			if err != nil {
				errResp := handler.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: err.Error(),
					Code:    "INVALID_TOKEN",
					Success: false,
				}

				handler.SendJSONResponse(w, errResp.Status, errResp)
				return
			}

			user := ctx.ContextUser{
				ID: claims.UserID,
				// Role: claims.Role,
			}

			// Store claims in context
			ctx := context.WithValue(r.Context(), ctx.UserKey, user)
			r = r.WithContext(ctx)

			// Pass to nextHandler with new context
			next.ServeHTTP(w, r)

		})
	}
}
