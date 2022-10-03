package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bernardn38/gobank/identity-service/helpers"
	"github.com/bernardn38/gobank/identity-service/sql/users"
	"github.com/bernardn38/gobank/identity-service/token"
	"github.com/cristalhq/jwt/v4"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Handler struct {
	UsersDb      *users.Queries
	TokenManager *token.Manager
}

func (h *Handler) GetUserView(w http.ResponseWriter, r *http.Request) {
	userIdentity := r.Context().Value("token")

	claims := userIdentity.(*jwt.RegisteredClaims)

	if claims.ID != chi.URLParam(r, "userId") {
		helpers.ResponseNoPayload(w, 401)
		return
	}

	user, err := h.UsersDb.GetUserView(context.Background(), uuid.Must(uuid.FromBytes([]byte(claims.ID))))
	if err != nil {
		log.Println(err)
		helpers.ResponseWithPayload(w, 404, []byte(fmt.Sprintf("Error finding user with id %s", claims.ID)))
		return
	}
	userResp, err := json.Marshal(user)
	if err != nil {
		helpers.ResponseNoPayload(w, 500)
		return
	}
	helpers.ResponseWithPayload(w, 200, userResp)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	listUsers, err := h.UsersDb.ListUsers(context.Background())
	if err != nil {
		log.Println(err)
		helpers.ResponseNoPayload(w, 400)
		return
	}
	marshal, err := json.Marshal(listUsers)
	if err != nil {
		return
	}
	helpers.ResponseWithPayload(w, 200, marshal)
}
