package mg

import "testing"

func TestFatalExit(t *testing.T) {
	expected := 99
	code := ExitStatus(Fatal(expected))
	if code != expected {
		t.Fatalf("Expected code %v but got %v", expected, code)
	}
}

func TestFatalfExit(t *testing.T) {
	expected := 99
	code := ExitStatus(Fatalf(expected, "boo!"))
	if code != expected {
		t.Fatalf("Expected code %v but got %v", expected, code)
	}
}
