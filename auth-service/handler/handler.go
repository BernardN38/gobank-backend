package handler

import (
	"context"
	"github.com/bernardn38/gobank/auth-service/helpers"
	"github.com/bernardn38/gobank/auth-service/sql/users"
	"github.com/bernardn38/gobank/auth-service/token"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	UsersDb      *users.Queries
	TokenManager *token.Manager
}
type RegisterForm struct {
	Username  string `json:"username" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=128"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}
type LoginForm struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (handler *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering User")
	reqBody, _ := io.ReadAll(r.Body)
	form, err := ValidateRegisterForm(reqBody)
	if err != nil {
		helpers.ResponseWithPayload(w, 400, []byte(err.Error()))
		return
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 12)
	if err != nil {
		helpers.ResponseNoPayload(w, 500)
		return
	}

	form.Password = string(encryptedPassword)
	err = CreateUser(handler.UsersDb, form)
	if err != nil {
		log.Println(err)
		helpers.ResponseWithPayload(w, 400, []byte("Error registering user"))
		return
	}
	helpers.ResponseWithPayload(w, 200, []byte("Register Success"))
}
func (handler *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging in User")
	cookie, ok := CheckForValidCookie(r, handler)
	if ok {
		UpdateCookie(w, handler, cookie.ID)
		helpers.ResponseWithPayload(w, 200, []byte("User already logged in, refreshing token."))
		return
	}
	reqBody, _ := io.ReadAll(r.Body)
	log.Println(string(reqBody))
	form, err := ValidateLoginForm(reqBody)
	if err != nil {
		log.Println(err)
		helpers.ResponseWithPayload(w, 400, []byte(err.Error()))
		return
	}
	user, err := handler.UsersDb.GetUserByUsername(context.Background(), form.Username)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		helpers.ResponseNoPayload(w, 401)
		return
	}

	newToken, err := handler.TokenManager.GenerateToken(user.ID.String())
	if err != nil {
		return
	}
	log.Println("Log in successful")
	SetCookie(w, newToken)
	helpers.ResponseWithPayload(w, 200, []byte(newToken.String()))
}
