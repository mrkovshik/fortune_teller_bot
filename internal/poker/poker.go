package poker

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Poker struct {
	logger *zap.SugaredLogger
	url    string
}

func NewPoker(logger *zap.SugaredLogger, url string) *Poker {
	return &Poker{
		logger: logger,
		url:    url,
	}
}

func (p Poker) Poke(ch <-chan time.Time, done chan struct{}) {
	p.logger.Debug("Starting poke server")
	for range ch {
		p.logger.Debug("Poking url ", p.url)
		res, err := http.Get(p.url)
		if err != nil {
			p.logger.Error(err)
			continue
		}
		p.logger.Debug("Responded  ", res.StatusCode)
	}
	done <- struct{}{}
}
