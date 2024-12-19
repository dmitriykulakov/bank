package database

import (
	"bank/internal/config"
	"bank/internal/model"
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

var UuidChan = make(chan uuid.UUID)
var BalanceRespChan = make(chan interface{})
var PostRequestChan = make(chan *model.PostRequest)
var PostResponseChan = make(chan model.Response)

type db interface {
	getBalance(client uuid.UUID) interface{}
	changeWallet(req *model.PostRequest) model.Response
}

func Broadcast(ctx context.Context, wg *sync.WaitGroup, cfg *config.DbConfig) {
	var db db = connectToPgDB(cfg)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case client := <-UuidChan:
			BalanceRespChan <- db.getBalance(client)
		case req := <-PostRequestChan:
			PostResponseChan <- db.changeWallet(req)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
