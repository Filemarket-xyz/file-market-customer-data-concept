package domain

import (
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type Client struct {
	Id        int64
	Address   common.Address
	Agreement bool

	PointBalance decimal.Decimal
}

func (c *Client) ToModel() *models.Client {
	res := &models.Client{
		ID: c.Id,
	}
	return res
}
