package curl

/*
#include <curl/curl.h>

static CURLSHcode curl_share_setopt_long(CURLSH *handle, CURLSHoption option, long parameter) {
  return curl_share_setopt(handle, option, parameter);
}
static CURLSHcode curl_share_setopt_pointer(CURLSH *handle, CURLSHoption option, void *parameter) {
  return curl_share_setopt(handle, option, parameter);
}
*/
import "C"

import "unsafe"

// implement os.Error interface
type CurlShareError C.CURLMcode

func (e CurlShareError) Error() string {
	// ret is const char*, no need to free
	ret := C.curl_share_strerror(C.CURLSHcode(e))
	return C.GoString(ret)
}

func newCurlShareError(errno C.CURLSHcode) error {
	if errno == 0 { // if nothing wrong
		return nil
	}
	return CurlShareError(errno)
}

type CURLSH struct {
	handle unsafe.Pointer
}

func ShareInit() *CURLSH {
	p := C.curl_share_init()
	return &CURLSH{p}
}

func (shcurl *CURLSH) Cleanup() error {
	p := shcurl.handle
	return newCurlShareError(C.curl_share_cleanup(p))
}

func (shcurl *CURLSH) Setopt(opt int, param interface{}) error {
	p := shcurl.handle
	if param == nil {
		return newCurlShareError(C.curl_share_setopt_pointer(p, C.CURLSHoption(opt), nil))
	}
	switch opt {
	//	case SHOPT_LOCKFUNC, SHOPT_UNLOCKFUNC, SHOPT_USERDATA:
	//		panic("not supported")
	case SHOPT_SHARE, SHOPT_UNSHARE:
		if val, ok := param.(int); ok {
			return newCurlShareError(C.curl_share_setopt_long(p, C.CURLSHoption(opt), C.long(val)))
		}
	}
	panic("not supported CURLSH.Setopt opt or param")
	return nil
}
