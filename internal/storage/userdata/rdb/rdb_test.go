package rdb

import (
	"context"

	"github.com/mrkovshik/fortune_teller_bot/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ctx         = context.Background()
	dataStorage *UserDataStorage
	err         error
)

const (
	testChatID  = int64(138)
	testBookId1 = int64(25)
	testBookId2 = int64(26)
)

var _ = Describe("Redis DB storage", func() {

	It("Init DB", func() {
		dataStorage, err = NewUserDataStorage(ctx, &config.RDB{
			Addr:         redisHostPort,
			Password:     "",
			DB:           0,
			DialTimeout:  0,
			ReadTimeout:  0,
			WriteTimeout: 0,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(dataStorage).NotTo(BeNil())
	})
	It("Save data", func() {
		Expect(dataStorage.SaveBookID(ctx, testChatID, testBookId1)).To(Succeed())
		Expect(dataStorage.SaveBookID(ctx, testChatID, testBookId2)).To(Succeed())
		Expect(dataStorage.SaveLanguage(ctx, testChatID, config.Russian)).To(Succeed())
	})
	It("Get data", func() {
		bookID, err := dataStorage.GetBookID(ctx, testChatID)
		Expect(err).NotTo(HaveOccurred())
		Expect(bookID).To(Equal(testBookId2))
		lang, err := dataStorage.GetLanguage(ctx, testChatID)
		Expect(err).NotTo(HaveOccurred())
		Expect(lang).To(Equal(config.Russian))
	})
})
