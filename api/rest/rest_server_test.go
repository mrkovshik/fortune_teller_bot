package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/mrkovshik/fortune_teller_bot/api/rest"
	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	"github.com/mrkovshik/fortune_teller_bot/internal/model"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/book_storage/local"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/state_storage/in_memory"
	"github.com/mrkovshik/fortune_teller_bot/internal/update_processor/basic"
	"github.com/mrkovshik/yandex_diploma/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var (
	cfg       *config.Config
	logger    *zap.Logger
	err       error
	srv       api.Server
	stepStack = model.NewStepStack()
)

const testChatID = 111

var _ = Describe("MessageReplyHandler", Ordered, func() {
	BeforeAll(func() {
		logger, err = zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())
		testBookStorage := local.NewStorage(logger.Sugar())
		testStateStorage := in_memory.NewStateStorage() // TODO: use mock
		stepStack.Push(model.AskingQuestion)
		testStateStorage.Update(testChatID, &model.ChatState{
			StepStack: stepStack,
		})
		proc := basic.NewUpdateProcessor(testBookStorage, testStateStorage, logger.Sugar()) // TODO: use mock
		cfg, err = config.GetConfig()
		Expect(err).NotTo(HaveOccurred())
		srv = rest.NewRestAPIServer(proc, cfg, logger.Sugar())
		ctx := context.Background()
		go func() {
			err := srv.RunServer(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()
		err := waitForServer(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), 3*time.Second)
		Expect(err).NotTo(HaveOccurred())
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
		body, _ := json.Marshal(upd)
		url := fmt.Sprintf("http://%s:%s/telegram", cfg.Host, cfg.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		var reply map[string]interface{}
		err = json.Unmarshal(respBody, &reply)
		Expect(err).NotTo(HaveOccurred())

		Expect(reply).To(HaveKey("text"))
		Expect(len(reply["text"].(string))).To(BeNumerically(">", 20))
	})

	It("Responds request for specific book quote", func() {
		stepStack.Push(model.SelectStartCommand)
		stepStack.Push(model.SelectBook)
		upd := model.Update{
			CallbackQuery: &model.CallbackQuery{
				ID: "321",
				From: model.Chat{
					ID: testChatID,
				},
				Data: "2.fb2",
			},
		}
		body, _ := json.Marshal(upd)
		url := fmt.Sprintf("http://%s:%s/telegram", cfg.Host, cfg.Port)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		var reply map[string]interface{}
		err = json.Unmarshal(respBody, &reply)
		Expect(err).NotTo(HaveOccurred())

		Expect(reply).To(HaveKey("text"))
		Expect(len(reply["text"].(string))).To(BeNumerically(">", 20))
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
