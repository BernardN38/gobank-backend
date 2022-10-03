package handler

import (
	"context"
	"github.com/bernardn38/gobank/identity-service/helpers"
	"github.com/cristalhq/jwt/v4"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) VerifyJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := CheckForValidCookie(r, h)
		if !ok {
			helpers.ResponseNoPayload(w, 401)
			return
		}

		ctx := context.WithValue(r.Context(), "token", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func CheckForValidCookie(r *http.Request, handler *Handler) (*jwt.RegisteredClaims, bool) {
	cookie, err := r.Cookie("jwtToken")
	if err != nil {
		log.Println(err)
		return nil, false
	}

	cookieFields := strings.Split(cookie.String(), "=")
	if len(cookieFields) != 2 {
		log.Println("Cookie invalid")
		return nil, false
	}
	token := cookieFields[1]
	claims, ok := handler.TokenManager.VerifyToken(token)
	if !ok {
		return nil, false
	}

	return claims, true
}
