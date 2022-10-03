package handler

import (
	"context"
	"database/sql"
	"github.com/bernardn38/gobank/transfer-service/helpers"
	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	rawToken := r.Header.Get("Authorization")
	if len(rawToken) < 1 {
		log.Println("no token present")
		return nil, false
	}
	token := strings.TrimPrefix(rawToken, "Bearer ")
	claims, ok := handler.TokenManager.VerifyToken(token)
	if !ok {
		log.Println("token can not be verified")
		return nil, false
	}
	if !claims.IsValidAt(time.Now()) {
		return nil, false
	}
	return claims, true
}

func CalculateBalance(userId uuid.UUID, h *Handler) (int64, error) {
	var balance int64

	deposits, err := h.TransactionDb.GetDeposits(context.Background(), userId)
	if err != nil {
		log.Println(err)
		if err != sql.ErrNoRows {
			log.Println(err)
			return 0, err
		}
	}
	withdrawals, err := h.TransactionDb.GetWithdrawals(context.Background(), userId)
	if err != nil {
		log.Println(err)
		if err != sql.ErrNoRows {
			log.Println(err)
			return 0, err
		}
	}
	depositBytes := deposits.([]uint8)
	withdrawalBytes := withdrawals.([]uint8)

	depositTotal, err := strconv.Atoi(string(depositBytes))
	if err != nil {
		log.Println(err)
	}
	withdrawalTotal, err := strconv.Atoi(string(withdrawalBytes))
	if err != nil {
		log.Println(err)
	}
	balance += int64(depositTotal)
	balance -= int64(withdrawalTotal)
	return balance, nil
}
