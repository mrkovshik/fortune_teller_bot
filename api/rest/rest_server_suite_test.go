package rest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEmbedded(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Embedded suite")
}
