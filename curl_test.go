package cameradar

import (
	"reflect"
	"testing"

	curl "github.com/ullaakut/go-curl"
)

func TestCurl(t *testing.T) {
	handle := Curl{
		CURL: curl.EasyInit(),
	}

	handle2 := handle.Duphandle()

	if reflect.DeepEqual(handle, handle2) {
		t.Errorf("unexpected identical handle from duphandle: expected %+v got %+v", handle, handle2)
	}
}
