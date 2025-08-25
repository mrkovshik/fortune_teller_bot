package inmemory

import (
	"container/list"
	"errors"

	"github.com/mrkovshik/fortune_teller_bot/internal/storage/steps"
)

var (
	ErrStepNotFound = errors.New("step not found")
	ErrChatNotFound = errors.New("chat not found")
)

// StepStorage is a stack data structure implemented using container/list.
// It stores elements in LIFO (Last In, First Out) order.
type StepStorage map[int64]*list.List

// NewStepStorage creates and returns a new empty StepStack.
func NewStepStorage() *StepStorage {
	s := make(StepStorage)
	return &s
}

// Push adds an element to the top of the stack.
func (s StepStorage) Push(chatID int64, step steps.ChatStep) error {
	stack, exists := s[chatID]
	if !exists {
		stack = list.New()
	}
	stack.PushBack(step)
	return nil
}

// Pop removes and returns the top element of the stack.
// It returns ("", false) if the stack is empty.
func (s StepStorage) Pop(chatID int64) (steps.ChatStep, error) {
	stack, exists := s[chatID]
	if !exists {
		return "", ErrChatNotFound
	}
	el := stack.Back()
	if el == nil {
		return "", ErrStepNotFound
	}
	stack.Remove(el)
	step := el.Value.(steps.ChatStep)
	return step, nil
}

// Peek returns the top element of the stack without removing it.
// It returns ("", false) if the stack is empty.
func (s StepStorage) Peek(chatID int64) (steps.ChatStep, error) {
	stack, exists := s[chatID]
	if !exists {
		return "", ErrChatNotFound
	}
	el := stack.Back()
	if el == nil {
		return "", ErrStepNotFound
	}
	step := el.Value.(steps.ChatStep)
	return step, nil
}

// PeekPrevious returns the second element from the top of the stack without removing it.
// It returns ("", false) if the stack is empty or has only one element.
func (s StepStorage) PeekPrevious(chatID int64) (steps.ChatStep, error) {
	stack, exists := s[chatID]
	if !exists {
		return "", ErrChatNotFound
	}
	if stack.Len() < 2 {
		return "", ErrStepNotFound
	}
	el := stack.Back()
	prev := el.Prev()
	if prev == nil {
		return "", ErrStepNotFound
	}
	step := prev.Value.(steps.ChatStep)
	return step, nil
}

func (s StepStorage) Clear(chatID int64) {
	delete(s, chatID)
}
