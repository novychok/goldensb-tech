package service

import (
	"github.com/google/wire"
	"github.com/novychok/goldensbtech/internal/service/producer"
	"github.com/novychok/goldensbtech/internal/service/realtime"
)

var Set = wire.NewSet(
	producer.New,
	realtime.New,
)
