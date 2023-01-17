package demo_test

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gdemo/test"
)

func TestDemo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Demo Suite")
}

var _ = BeforeSuite(func() {
	dir, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
	}

	test.InitTestResource(dir)
})
