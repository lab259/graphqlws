package graphqlws_test

import (
	"fmt"
	"github.com/jamillosantos/macchiato"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestGraphWS(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	if os.Getenv("ENV") == "" {
		err := os.Setenv("ENV", "test")
		if err != nil {
			panic(err)
		}
	}
	dir, _ := os.Getwd()
	os.Stdout.WriteString(fmt.Sprintf("CWD: %s\n", dir))
	os.Stdout.WriteString(fmt.Sprintf("Starting with ENV: %s\n", os.Getenv("ENV")))
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "GraphQL WS Test Suite")
}
