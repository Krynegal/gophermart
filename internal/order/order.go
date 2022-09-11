package order

import (
	"github.com/Krynegal/gophermart.git/internal/status"
	"time"
)

type AccrualOrder struct {
	UserID     int           `json:"-"`
	Number     uint64        `json:"number,string"`
	Status     status.Status `json:"status"`
	Accrual    float32       `json:"accrual,omitempty"`
	UploadedAt time.Time     `json:"uploaded_at"`
}

type WithdrawOrder struct {
	UserID      int       `json:"-"`
	Order       uint64    `json:"order,string"`
	Sum         float32   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}

type ProcessingOrder struct {
	Order   uint64  `json:"order,string"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
