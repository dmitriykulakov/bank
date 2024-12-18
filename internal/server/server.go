package server

import (
	"bank/internal/config"
	"bank/internal/database"
	"bank/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "bank/docs"

	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Handle(cfg *config.ServerConfig) {
	router := http.NewServeMux()
	router.HandleFunc("GET /api/v1/docs/", httpSwagger.Handler(
		httpSwagger.URL("http://"+cfg.Address+"/api/v1/docs/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	router.HandleFunc("POST /api/v1/wallet", postWallet)
	router.HandleFunc("GET /api/v1/wallets/{WALLET_UUID}", getWallet)
	log.Printf(cfg.Address)
	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(cfg.Address, router))
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
}

// @Summary post_wallet
// @Tags postWallet
// @Accept json
// @Produce json
// @Param input body model.PostRequest true "postWallet"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/wallet [post]
func postWallet(w http.ResponseWriter, r *http.Request) {
	req, err := parseRequest(r)
	if err != nil {
		encode(&w, err.Error(), http.StatusBadRequest)
		return
	}
	database.PostRequestChan <- req
	resp := <-database.PostResponseChan
	if resp.Err != nil {
		encode(&w, resp.Err.Error(), http.StatusBadRequest)
		return
	}
	encode(&w, fmt.Sprintf("id %s: Balance Changed", req.WalletId), http.StatusOK)
}

// @Summary get_wallet
// @Tags getWallet
// @Accept json
// @Produce json
// @Param WALLET_UUID path string true "WALLET UUID"
// @Success 200 {object} model.RespBalance
// @Failure 400 {object} model.Response
// @Failure 404 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/wallets/{WALLET_UUID} [get]
func getWallet(w http.ResponseWriter, r *http.Request) {
	tmp := r.PathValue("WALLET_UUID")
	str := strings.ReplaceAll(tmp, "-", "")
	client, err := uuid.Parse(str)
	if err != nil {
		encode(&w, fmt.Sprintf("%s:%s", err, tmp), http.StatusBadRequest)
		return
	}
	fmt.Println(client)
	database.UuidChan <- client
	resp := <-database.BalanceRespChan
	switch value := resp.(type) {
	case model.Response:
		encode(&w, value.Err.Error(), http.StatusBadRequest)
		log.Printf("get Amount %s: error: %s", tmp, value.Err)
	case model.RespBalance:
		encode(&w, value.Balance, http.StatusOK)
		log.Printf("get Amount %s:%d", tmp, value.Balance)
	}
}

func parseRequest(r *http.Request) (*model.PostRequest, error) {
	var req *model.PostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if req.WalletId == uuid.Nil {
		err = errors.New("the walletId has wrong format")
		return nil, err
	}
	if req.OperationType != "WITHDRAW" && req.OperationType != "DEPOSIT" {
		err = fmt.Errorf("the id \"%s\" operation type \"%s\" is wrong", req.WalletId.String(), req.OperationType)
		return nil, err
	}
	if req.Amount < 0 {
		err = fmt.Errorf("the id \"%s\" amount %d is less than zero", req.WalletId.String(), req.Amount)
		return nil, err
	}
	return req, nil
}

func encode(w *http.ResponseWriter, elem interface{}, status int) {
	(*w).WriteHeader(status)
	if _, ok := elem.(int); !ok {
		log.Print(elem)
	}
	err := json.NewEncoder(*w).Encode(elem)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		log.Print(err)
	}
}
