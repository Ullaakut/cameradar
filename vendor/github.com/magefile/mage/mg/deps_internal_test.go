package mg

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestDepsLogging(t *testing.T) {
	os.Setenv("MAGEFILE_VERBOSE", "1")
	defer os.Unsetenv("MAGEFILE_VERBOSE")
	buf := &bytes.Buffer{}

	defaultLogger := logger
	logger = log.New(buf, "", 0)
	defer func() { logger = defaultLogger }()

	foo()

	if strings.Count(buf.String(), "Running dependency: github.com/magefile/mage/mg.baz") != 1 {
		t.Fatalf("expected one baz to be logged, but got\n%s", buf)
	}
}

func foo() {
	Deps(bar, baz)
}

func bar() {
	Deps(baz)
}

func baz() {}
