package poker

import (
	"time"

	"go.uber.org/zap"
)

type Poker struct {
	logger *zap.SugaredLogger
	url    string
}

func (p Poker) Poke(ch <-chan time.Time, done chan struct{}) {
	for range ch {
		p.logger.Debug("Starting to poke server")
	}
	done <- struct{}{}
}
