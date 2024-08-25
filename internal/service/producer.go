package service

import (
	"context"
)

type Producer interface {
	Produce(ctx context.Context, notifyCh chan struct{}) error
}
