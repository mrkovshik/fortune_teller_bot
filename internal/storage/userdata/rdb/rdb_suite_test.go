package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEmbedded(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Embedded suite")
}

var (
	stopContainer func()
	redisHostPort string
)

var _ = BeforeSuite(func() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("6379/tcp"),
		).WithDeadline(30 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).ToNot(HaveOccurred())

	// получить хост:порт, чтобы коннектиться из тестов
	host, err := c.Host(ctx)
	Expect(err).ToNot(HaveOccurred())

	mappedPort, err := c.MappedPort(ctx, "6379")
	Expect(err).ToNot(HaveOccurred())

	redisHostPort = fmt.Sprintf("%s:%s", host, mappedPort.Port())

	stopContainer = func() { _ = c.Terminate(context.Background()) }
})

var _ = AfterSuite(func() {
	if stopContainer != nil {
		stopContainer()
	}
})
