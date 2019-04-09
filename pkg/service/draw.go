package service

import (
	"fmt"
	"time"
)

const PoolTransferTypeNone = "NONE"

type Draw struct {
	ProductID        string      `json:"productId"`
	ID               int         `json:"id"`
	Date             time.Time   `json:"date"`
	PrimaryNumbers   []int       `json:"primaryNumbers"`
	SecondaryNumbers []int       `json:"secondaryNumbers"`
	Dividends        []*Dividend `json:"dividends"`
}

func (d *Draw) String() string {
	return fmt.Sprintf("%v", append(d.PrimaryNumbers, d.SecondaryNumbers...))
}

type Dividend struct {
	Division               int     `json:"division"`
	BlocDividend           float32 `json:"blocDividend"`
	BlocNumberOfWinners    int     `json:"blocNumberOfWinners"`
	CompanyID              string  `json:"companyId"`
	CompanyDividend        float32 `json:"companyDividend"`
	CompanyNumberOfWinners int     `json:"companyNumberOfWinners"`
	PoolTransferType       string  `json:"poolTransferType"`
	PoolTransferredTo      int     `json:"poolTransferredTo"`
}
