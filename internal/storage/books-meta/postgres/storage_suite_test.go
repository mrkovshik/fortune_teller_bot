package postgres

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
	DSN           string
	stopContainer func()
)

var _ = BeforeSuite(func() {
	ctx := context.Background()

	user, pass, dbname := "testuser", "testpass", "testdb"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       dbname,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).ToNot(HaveOccurred())

	stopContainer = func() { _ = c.Terminate(context.Background()) }

	host, err := c.Host(ctx)
	Expect(err).ToNot(HaveOccurred())
	port, err := c.MappedPort(ctx, "5432/tcp")
	Expect(err).ToNot(HaveOccurred())

	DSN = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pass, host, port.Int(), dbname)
})

var _ = AfterSuite(func() {
	if stopContainer != nil {
		stopContainer()
	}
})
