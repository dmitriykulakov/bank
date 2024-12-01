package model

import (
	"github.com/google/uuid"
)

type PostRequest struct {
	OperationType string `json:"operationType"`
	Client
}

type Client struct {
	WalletId uuid.UUID `json:"walletId"`
	Amount   int       `json:"amount"`
}

type RespBalance struct {
	Balance int
	Err     error
}

type Response struct {
	Err error
}
