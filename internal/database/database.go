package database

import (
	"bank/internal/config"
	"bank/internal/model"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var UuidChan = make(chan uuid.UUID)
var BalanceRespChan = make(chan interface{})
var PostRequestChan = make(chan *model.PostRequest)
var PostResponseChan = make(chan model.Response)

func Broadcast(ctx context.Context, wg *sync.WaitGroup, cfg *config.PgConfig) {
	db := connectToDB(cfg)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case client := <-UuidChan:
			BalanceRespChan <- getBalance(client, db)
		case req := <-PostRequestChan:
			PostResponseChan <- operation(req, db)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func connectToDB(cfg *config.PgConfig) *gorm.DB {

	cfgPG := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port)

	db, err := gorm.Open(postgres.Open(cfgPG), &gorm.Config{})
	for i := 0; i < 10 && err != nil; i++ {
		time.Sleep(time.Second * 5)
		db, err = gorm.Open(postgres.Open(cfgPG), &gorm.Config{})
		if err != nil {
			log.Printf("ConnectToDB: error to connect, please wait %v", err)
		}
	}
	if err != nil {
		log.Printf("ConnectToDB: error to connect %v", err)
		return nil
	}
	log.Print("The PSQL is ready")
	return db
}

func getBalance(client uuid.UUID, db *gorm.DB) interface{} {
	var balance []int
	db.Table("clients").Where("walletId = ?", client).Select("amount").Find(&balance)
	if len(balance) == 1 {
		return model.RespBalance{Balance: balance[0]}
	} else {
		return model.Response{Err: fmt.Errorf("the id \"%s\"doesn't exist", client.String())}
	}
}

func operation(req *model.PostRequest, db *gorm.DB) model.Response {
	var balance []int
	db.Table("clients").Where("walletId = ?", req.WalletId).Select("amount").Find(&balance)
	if len(balance) == 1 {
		if req.OperationType == "DEPOSIT" {
			db.Table("clients").Model(&req.Client).Where("walletId = ?", req.WalletId).Update("amount", balance[0]+req.Amount)
		}
		if req.OperationType == "WITHDRAW" {
			if req.Amount <= balance[0] {
				db.Table("clients").Model(&req.Client).Where("walletId = ?", req.WalletId).Update("amount", balance[0]-req.Amount)
			} else {
				return model.Response{Err: fmt.Errorf("the id \"%s\" balance is lower than withdraw", req.WalletId.String())}
			}
		}
	} else {
		return model.Response{Err: fmt.Errorf("the id \"%s\" doesn't exist", req.WalletId.String())}

	}
	return model.Response{Err: nil}
}
