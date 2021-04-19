package configuration_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoConfig Suite")
}
