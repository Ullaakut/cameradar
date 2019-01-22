//+build CI

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBootstrap(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	s, err := run("go", "run", "bootstrap.go")
	if err != nil {
		t.Fatal(s)
	}
	name := "mage"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	if _, err := os.Stat(filepath.Join(os.Getenv("GOPATH"), "bin", name)); err != nil {
		t.Fatal(err)
	}
}

func run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Env = os.Environ()
	b, err := c.CombinedOutput()
	return string(b), err
}
