package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Krynegal/gophermart.git/internal/order"
	"github.com/Krynegal/gophermart.git/internal/storage"
	"golang.org/x/time/rate"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const tooManyRequestTemplate = "No more than %d requests per minute allowed"

type AccrualProcessor interface {
	GetOrderStatus(orderID string) *order.ProcessingOrder
	Run()
}

type accrualProcessor struct {
	address   string
	rateLimit int
	poolSize  int
	limiter   *rate.Limiter
	db        storage.Storager

	OrderQueue chan string
}

func NewAccrualProcessor(db storage.Storager, addr string, poolSize int) AccrualProcessor {
	ap := &accrualProcessor{
		db:         db,
		address:    addr,
		poolSize:   poolSize,
		OrderQueue: make(chan string, poolSize),
	}
	for i := 0; i < poolSize; i++ {
		go ap.queueWorker()
	}
	return ap
}

var wg sync.WaitGroup

func (a *accrualProcessor) Run() {
	go func() {
		for {
			orderList, err := a.db.GetOrdersForProcessing(a.poolSize)
			if err != nil || len(orderList) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}
			wg.Add(len(orderList))
			for _, orderID := range orderList {
				a.OrderQueue <- orderID
			}
			wg.Wait()
		}
	}()
}

func (a *accrualProcessor) GetOrderStatus(orderID string) *order.ProcessingOrder {
	res, err := http.Get(fmt.Sprintf("%s/api/orders/%s", a.address, orderID))
	if err != nil {
		log.Println("request error", err)
		return nil
	}
	defer func() {
		err = res.Body.Close()
	}()

	if res.StatusCode == http.StatusTooManyRequests {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return nil
		}
		var rl int
		_, err = fmt.Sscanf(tooManyRequestTemplate, string(resBody), &rl)
		if err != nil {
			log.Println(err)
			return nil
		}
		a.setLimit(rl)
	}
	if res.StatusCode != http.StatusOK {
		return nil
	}
	var aos order.ProcessingOrder
	if err = json.NewDecoder(res.Body).Decode(&aos); err != nil {
		log.Println(err)
		return nil
	}
	return &aos
}

func (a *accrualProcessor) setLimit(n int) {
	if n <= 0 {
		a.limiter = nil
		return
	}
	a.limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(n)), n)
}

func (a *accrualProcessor) queueWorker() {
	for orderID := range a.OrderQueue {
		if a.limiter != nil && !a.limiter.Allow() {
			err := a.limiter.Wait(context.Background())
			if err != nil {
				log.Println(err)
				wg.Done()
				return
			}
		}
		orderStatus := a.GetOrderStatus(orderID)
		if orderStatus != nil {
			if err := a.db.UpdateOrderState(orderStatus); err != nil {
				log.Println(err)
			}
		}
		wg.Done()
	}
}
