package client

import (
	"context"

	"github.com/leehaowei/tolling-micro-service/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}
