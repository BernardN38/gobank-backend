package handler

import (
	"context"
	"encoding/json"
	"github.com/bernardn38/gobank/transfer-service/helpers"
	"github.com/bernardn38/gobank/transfer-service/sql/transactions"
	"github.com/bernardn38/gobank/transfer-service/token"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	TransactionDb *transactions.Queries
	TokenManager  *token.Manager
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	allTransfers, err := h.TransactionDb.GetAllTransactions(context.Background())
	if err != nil {
		return
	}
	marshal, err := json.Marshal(allTransfers)
	if err != nil {
		return
	}
	helpers.ResponseWithPayload(w, 200, marshal)
}

type Deposit struct {
	Amount      int64     `json:"amount"`
	ToAccount   uuid.UUID `json:"to_account"`
	FromAccount uuid.UUID `json:"from_account"`
	CreatedAt   time.Time `json:"created_at"`
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	log.Println("Initiating transfer")

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var deposit Deposit
	err = json.Unmarshal(reqBody, &deposit)
	if err != nil {
		log.Println(err)
		return
	}
	fromBalance, err := CalculateBalance(deposit.FromAccount, h)
	if err != nil {
		log.Println(err)
		helpers.ResponseNoPayload(w, http.StatusInternalServerError)
		return
	}
	if fromBalance < deposit.Amount {
		log.Println("insufficient funds")
		helpers.ResponseWithPayload(w, http.StatusBadRequest, []byte("Insufficient funds"))
		return
	}

	_, err = h.TransactionDb.CreateTransaction(context.Background(), transactions.CreateTransactionParams{
		Amount:      deposit.Amount,
		FromAccount: deposit.FromAccount,
		ToAccount:   deposit.ToAccount,
		CreatedAt:   time.Now()},
	)
	if err != nil {
		helpers.ResponseNoPayload(w, http.StatusInternalServerError)
		return
	}
	log.Println("transfer successful")
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		helpers.ResponseNoPayload(w, http.StatusInternalServerError)
		return
	}
	balance, err := CalculateBalance(parsedId, h)
	if err != nil {
		helpers.ResponseNoPayload(w, 400)
		return
	}
	jsonResp, err := json.Marshal(balance)
	if err != nil {
		return
	}
	helpers.ResponseWithPayload(w, 200, jsonResp)
}
