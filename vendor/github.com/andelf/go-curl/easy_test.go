package curl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"sync"
)

func setupTestServer(serverContent string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, serverContent)
	}))
}

func TestEasyInterface(t *testing.T) {
	ts := setupTestServer("")
	defer ts.Close()

	easy := EasyInit()
	defer easy.Cleanup()

	easy.Setopt(OPT_URL, ts.URL)
	if err := easy.Perform(); err != nil {
		t.Fatal(err)
	}
}

func TestCallbackFunction(t *testing.T) {
	serverContent := "A random string"
	ts := setupTestServer(serverContent)
	defer ts.Close()

	easy := EasyInit()
	defer easy.Cleanup()

	easy.Setopt(OPT_URL, ts.URL)
	easy.Setopt(OPT_WRITEFUNCTION, func(buf []byte, userdata interface{}) bool {
		result := string(buf)
		expected := serverContent + "\n"
		if result != expected {
			t.Errorf("output should be %q and is %q.", expected, result)
		}
		return true
	})
	if err := easy.Perform(); err != nil {
		t.Fatal(err)
	}
}

func TestEscape(t *testing.T) {
	easy := EasyInit()
	defer easy.Cleanup()

	payload := `payload={"msg": "First line\nSecond Line"}`
	expected := `payload%3D%7B%22msg%22%3A%20%22First%20line%5CnSecond%20Line%22%7D`
	result := easy.Escape(payload)
	if result != expected {
		t.Errorf("escaped output should be %q and is %q.", expected, result)
	}
}

func TestConcurrentInitAndCleanup(t *testing.T) {
	c := 2
	var wg sync.WaitGroup
	wg.Add(c)
	for i := 0; i < c; i++ {
		go func() {
			wg.Done()
			easy := EasyInit()
			defer easy.Cleanup()
		}()
	}

	wg.Wait()
}
