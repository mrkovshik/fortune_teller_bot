package inmemory_test

import (
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps/inmemory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testChatID = 111

var _ = Describe("", Ordered, func() {
	var (
		err     error
		storage *inmemory.StepStorage
	)
	BeforeAll(func() {
		storage = inmemory.NewStepStorage()
	})

	It("Initial steps", func() {
		err = storage.Push(testChatID, steps.SelectStartCommand)
		Expect(err).NotTo(HaveOccurred())
		command, err := storage.Peek(testChatID)
		Expect(err).NotTo(HaveOccurred())
		Expect(command).To(BeEquivalentTo(steps.SelectStartCommand))
	})

	It("Takes Previous step", func() {
		err = storage.Push(testChatID, steps.SelectBook)
		Expect(err).NotTo(HaveOccurred())
		command, err := storage.PeekPrevious(testChatID)
		Expect(err).NotTo(HaveOccurred())
		Expect(command).To(BeEquivalentTo(steps.SelectStartCommand))
	})

	It("Pops step", func() {
		command, err := storage.Pop(testChatID)
		Expect(err).NotTo(HaveOccurred())
		command, err = storage.Peek(testChatID)
		Expect(err).NotTo(HaveOccurred())
		Expect(command).To(BeEquivalentTo(steps.SelectStartCommand))
		command, err = storage.PeekPrevious(testChatID)
		Expect(err).To(HaveOccurred())
		Expect(err).To(BeEquivalentTo(steps.ErrStepNotFound))

	})
})
