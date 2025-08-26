package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/updateprocessor/basic"
	mock "github.com/mrkovshik/fortune_teller_bot/mocks"

	"github.com/mrkovshik/yandex_diploma/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var (
	cfg             *config.Config
	logger          *zap.Logger
	err             error
	srv             api.Server
	ctrl            *gomock.Controller
	updateProcessor *mock.MockUpdateProcessor
)

const testChatID = 111

var _ = Describe("MessageReplyHandler", Ordered, func() {
	BeforeAll(func() {
		ctrl = gomock.NewController(GinkgoT())
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		updateProcessor = mock.NewMockUpdateProcessor(ctrl)
		cfg, err = config.GetConfig()
		Expect(err).NotTo(HaveOccurred())
		srv = rest.NewRestAPIServer(updateProcessor, cfg, logger.Sugar())
		ctx := context.Background()
		go func() {
			err := srv.RunServer(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()
		err := waitForServer(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), 3*time.Second)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterAll(func() {
		DeferCleanup(ctrl.Finish)
	})

	It("Responds start command with menu", func() {

		upd := model.Update{
			Message: &model.Message{
				Chat: model.Chat{
					ID: testChatID,
				},
				Text: "/start",
			},
		}
		payload := map[string]interface{}{
			"chat_id":      testChatID,
			"text":         "Что бы вы хотели сделать?",
			"reply_markup": basic.StartMenu,
		}
		updateProcessor.EXPECT().ProcessMessage(upd.Message).Return(payload, nil)
		body, _ := json.Marshal(upd)
		url := fmt.Sprintf("http://%s:%s/telegram", cfg.Host, cfg.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("Responds request for random book quote", func() {
		upd := model.Update{
			Message: &model.Message{
				Chat: model.Chat{
					ID: testChatID,
				},
				Text: "Some random text",
			},
		}
		payload := map[string]interface{}{
			"chat_id": testChatID,
			"text":    "Some random quote",
		}
		updateProcessor.EXPECT().ProcessMessage(upd.Message).Return(payload, nil)
		body, _ := json.Marshal(upd)
		url := fmt.Sprintf("http://%s:%s/telegram", cfg.Host, cfg.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		defer resp.Body.Close()
	})

	It("Responds request for specific book quote", func() {
		upd := model.Update{
			CallbackQuery: &model.CallbackQuery{
				ID: "321",
				From: model.Chat{
					ID: testChatID,
				},
				Data: "2.fb2",
			},
		}
		payload := map[string]interface{}{
			"chat_id": testChatID,
			"text":    "Some random quote",
		}
		updateProcessor.EXPECT().ProcessCallback(upd.CallbackQuery).Return(payload, nil)
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
