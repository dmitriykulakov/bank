package testPack

import (
	"bank/internal/config"
	"bank/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("config.env"); err != nil {
		log.Print("No .env file found")
	}
}

func TestPost(t *testing.T) {
	var cfg = config.NewConfig()
	var balance int
	resp, _ := http.Get("http://" + cfg.Address + "/api/v1/wallets/e0eebc999c0b4ef8bb6d6bb9bd380a11")

	json.NewDecoder(resp.Body).Decode(&balance)
	id, _ := uuid.Parse("e0eebc999c0b4ef8bb6d6bb9bd380a11")

	tmp := model.PostRequest{OperationType: "DEPOSIT", Client: model.Client{WalletId: id, Amount: 100}}
	for i := 0; i < 1000; i++ {
		operJson, err := json.Marshal(tmp)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(operJson)
		http.Post("http://"+cfg.Address+"/api/v1/wallet", "Apllication/json", req)
	}

	resp, _ = http.Get("http://" + cfg.Address + "/api/v1/wallets/e0eebc999c0b4ef8bb6d6bb9bd380a11")
	var balanceNew int

	json.NewDecoder(resp.Body).Decode(&balanceNew)
	if balanceNew-balance != 100000 {
		t.Errorf("Test DEPOSIT: FAIL")
	} else {
		t.Logf("Test DEPOSIT: OK")
	}

	tmp = model.PostRequest{OperationType: "WITHDRAW", Client: model.Client{WalletId: id, Amount: 100}}
	for i := 0; i < 1000; i++ {
		operJson, err := json.Marshal(tmp)
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(operJson)
		http.Post("http://"+cfg.Address+"/api/v1/wallet", "Apllication/json", req)
	}

	resp, _ = http.Get("http://" + cfg.Address + "/api/v1/wallets/e0eebc999c0b4ef8bb6d6bb9bd380a11")

	json.NewDecoder(resp.Body).Decode(&balance)
	if balanceNew-balance != 100000 {
		t.Errorf("Test WITHDRAW: FAIL")
	} else {
		t.Logf("Test WITHDRAW: OK")
	}
}

func TestBadRequest(t *testing.T) {
	var cfg = config.NewConfig()
	testBadRequest := []struct {
		uuid          string
		operationType string
		amount        int
		want          string
		statusCode    int
	}{
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "DEPOSIT", 500, "id e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11: Balance Changed", 200},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "WITHDRAW", 500, "id e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11: Balance Changed", 200},
		{"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "DEPOSIT", 500, "id e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11: Balance Changed", 200},
		{"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "WITHDRAW", 500, "id e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11: Balance Changed", 200},
		{"e1eebc999c0b4ef8bb6d6bb9bd380a11", "DEPOSIT", 500, "the id \"e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" doesn't exist", 400},
		{"e1eebc999c0b4ef8bb6d6bb9bd380a11", "WITHDRAW", 500, "the id \"e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" doesn't exist", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a111", "DEPOSIT", 500, "the walletId has wrong format", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a111", "WITHDRAW", 500, "the walletId has wrong format", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "DEPOSITT", 500, "the id \"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" operation type \"DEPOSITT\" is wrong", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "WITHDRAWT", 500, "the id \"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" operation type \"WITHDRAWT\" is wrong", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "DEPOSIT", -500, "the id \"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" amount -500 is less than zero", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "WITHDRAW", -500, "the id \"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" amount -500 is less than zero", 400},
		{"e0eebc999c0b4ef8bb6d6bb9bd380a11", "WITHDRAW", 5000000, "the id \"e0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11\" balance is lower than withdraw", 400},
	}
	for i, r := range testBadRequest {
		id, _ := uuid.Parse(r.uuid)
		Json, err := json.Marshal(model.PostRequest{OperationType: r.operationType, Client: model.Client{WalletId: id, Amount: r.amount}})
		if err != nil {
			log.Fatal(err)
		}
		req := bytes.NewReader(Json)

		resp, err := http.Post("http://"+cfg.Address+"/api/v1/wallet", "Apllication/json", req)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != r.statusCode {
			t.Errorf("Test %d: FAIL, expected %d, result %d", i+1, r.statusCode, resp.StatusCode)
		} else {
			t.Logf("Test %d: OK, expected %d, result %d: OK", i+1, r.statusCode, resp.StatusCode)
		}
		var response string
		json.NewDecoder(resp.Body).Decode(&response)
		if response != r.want {
			t.Errorf("Test %d: FAIL, expected %s, result %s", i+1, r.want, response)
		} else {
			t.Logf("Test %d: OK, expected %s, result %s: OK", i+1, r.want, response)
		}
	}
}
