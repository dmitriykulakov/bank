package database

import (
	"bank/internal/config"
	"bank/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pg struct {
	pg *gorm.DB
}

func connectToPgDB(cfg *config.DbConfig) *pg {
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
	return &pg{db}
}

func (pgDB *pg) getBalance(client uuid.UUID) interface{} {
	var balance []int
	pgDB.pg.Table("clients").Where("walletId = ?", client).Select("amount").Find(&balance)
	if len(balance) == 1 {
		return model.RespBalance{Balance: balance[0]}
	} else {
		return model.Response{Err: fmt.Errorf("the id \"%s\"doesn't exist", client.String())}
	}
}

func (pgDB *pg) changeWallet(req *model.PostRequest) model.Response {
	var balance []int
	pgDB.pg.Table("clients").Where("walletId = ?", req.WalletId).Select("amount").Find(&balance)
	if len(balance) == 1 {
		if req.OperationType == "DEPOSIT" {
			pgDB.pg.Table("clients").Model(&req.Client).Where("walletId = ?", req.WalletId).Update("amount", balance[0]+req.Amount)
		}
		if req.OperationType == "WITHDRAW" {
			if req.Amount <= balance[0] {
				pgDB.pg.Table("clients").Model(&req.Client).Where("walletId = ?", req.WalletId).Update("amount", balance[0]-req.Amount)
			} else {
				return model.Response{Err: fmt.Errorf("the id \"%s\" balance is lower than withdraw", req.WalletId.String())}
			}
		}
	} else {
		return model.Response{Err: fmt.Errorf("the id \"%s\" doesn't exist", req.WalletId.String())}
	}
	return model.Response{Err: nil}
}
