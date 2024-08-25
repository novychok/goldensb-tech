package pkg

import (
	"github.com/novychok/goldensbtech/internal/pkg/log"

	"github.com/google/wire"
	"github.com/novychok/goldensbtech/internal/pkg/nats"
)

var Set = wire.NewSet(
	log.New,
	nats.New,
)
