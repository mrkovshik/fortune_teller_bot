package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/in_memory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"
	"github.com/mrkovshik/yandex_diploma/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("", func() {
	var (
		cfg    *config.Config
		logger *zap.Logger
		err    error
		srv    api.Server
	)
	BeforeEach(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		testBookStorage := local.NewStorage(logger.Sugar())
		testStateStorage := in_memory.NewStateStorage()                                     // TODO: use mock
		proc := basic.NewUpdateProcessor(testBookStorage, testStateStorage, logger.Sugar()) // TODO: use mock
		cfg, err = config.GetConfig()
		Expect(err).NotTo(HaveOccurred())
		srv = rest.NewRestAPIServer(proc, cfg, logger.Sugar())
	})

	It("", func() {
		ctx := context.Background()
		go func() {
			err := srv.RunServer(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()
		err := waitForServer(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), 3*time.Second)
		Expect(err).NotTo(HaveOccurred())
		upd := model.Update{
			Message: &model.Message{
				Chat: model.Chat{
					ID: 111,
				},
				Text: update_processor.GetMagicCommandName,
			},
		}
		body, _ := json.Marshal(upd)
		url := fmt.Sprintf("http://%s:%s/telegram", cfg.Host, cfg.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})

func waitForServer(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return fmt.Errorf("server %s not available after %s", address, timeout)
}
