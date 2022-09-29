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
	GetOrderStatus(orderID string) (*order.ProcessingOrder, error)
	Run()
}

type accrualProcessor struct {
	address       string
	endPointTempl string
	rateLimit     int
	poolSize      int
	limiter       *rate.Limiter
	db            storage.Storager

	OrderQueue chan string
}

func NewAccrualProcessor(db storage.Storager, addr string, poolSize int) AccrualProcessor {
	ap := &accrualProcessor{
		db:            db,
		address:       addr,
		endPointTempl: "%s/api/orders/%s",
		poolSize:      poolSize,
		OrderQueue:    make(chan string, poolSize),
	}
	for i := 0; i < poolSize; i++ {
		go ap.queueWorker()
	}
	return ap
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
		orderStatus, err := a.GetOrderStatus(orderID)
		if err != nil {
			log.Println(err)
		}
		if orderStatus != nil {
			if err = a.db.UpdateOrderState(orderStatus); err != nil {
				log.Println(err)
			}
		}
		wg.Done()
	}
}

var wg sync.WaitGroup

func (a *accrualProcessor) Run() {
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
}

func (a *accrualProcessor) GetOrderStatus(orderID string) (*order.ProcessingOrder, error) {
	res, err := http.Get(fmt.Sprintf(a.endPointTempl, a.address, orderID))
	if err != nil {
		log.Println("request error", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		var resBody []byte
		resBody, err = io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var rl int
		_, err = fmt.Sscanf(tooManyRequestTemplate, string(resBody), &rl)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		a.setLimit(rl)
	}
	if res.StatusCode != http.StatusOK {
		return nil, err
	}
	var aos order.ProcessingOrder
	if err = json.NewDecoder(res.Body).Decode(&aos); err != nil {
		log.Println(err)
		return nil, err
	}
	return &aos, nil
}

func (a *accrualProcessor) setLimit(n int) {
	if n <= 0 {
		a.limiter = nil
		return
	}
	a.limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(n)), n)
}
