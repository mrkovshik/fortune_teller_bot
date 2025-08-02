package storage

type Storage interface {
	GetRandomSentenceFromBook(bookName string) (string, error)
	ListBooks() ([]string, error)
}
