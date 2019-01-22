package curl

/*
#cgo freebsd CFLAGS: -I/usr/local/include
#cgo freebsd LDFLAGS: -L/usr/local/lib -lcurl
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

*/
import "C"

import (
	"unsafe"
)

//export goCallHeaderFunction
func goCallHeaderFunction(ptr *C.char, size C.size_t, ctx unsafe.Pointer) uintptr {
	curl := context_map.Get(uintptr(ctx))
	buf := C.GoBytes(unsafe.Pointer(ptr), C.int(size))
	if (*curl.headerFunction)(buf, curl.headerData) {
		return uintptr(size)
	}
	return C.CURL_WRITEFUNC_PAUSE
}

//export goCallWriteFunction
func goCallWriteFunction(ptr *C.char, size C.size_t, ctx unsafe.Pointer) uintptr {
	curl := context_map.Get(uintptr(ctx))
	buf := C.GoBytes(unsafe.Pointer(ptr), C.int(size))
	if (*curl.writeFunction)(buf, curl.writeData) {
		return uintptr(size)
	}
	return C.CURL_WRITEFUNC_PAUSE
}

//export goCallProgressFunction
func goCallProgressFunction(dltotal, dlnow, ultotal, ulnow C.double, ctx unsafe.Pointer) int {
	curl := context_map.Get(uintptr(ctx))
	if (*curl.progressFunction)(float64(dltotal), float64(dlnow),
		float64(ultotal), float64(ulnow),
		curl.progressData) {
		return 0
	}
	return 1
}

//export goCallReadFunction
func goCallReadFunction(ptr *C.char, size C.size_t, ctx unsafe.Pointer) uintptr {
	curl := context_map.Get(uintptr(ctx))
	buf := C.GoBytes(unsafe.Pointer(ptr), C.int(size))
	ret := (*curl.readFunction)(buf, curl.readData)
	str := C.CString(string(buf))
	defer C.free(unsafe.Pointer(str))
	if C.memcpy(unsafe.Pointer(ptr), unsafe.Pointer(str), C.size_t(ret)) == nil {
		panic("read_callback memcpy error!")
	}
	return uintptr(ret)
}
