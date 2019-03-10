// This file depends on functionality not available on Windows, hence we
// must skip it. https://github.com/andelf/go-curl/issues/48

// +build !windows

package curl

/*
#include <stdlib.h>
#include <curl/curl.h>

static CURLMcode curl_multi_setopt_long(CURLM *handle, CURLMoption option, long parameter) {
  return curl_multi_setopt(handle, option, parameter);
}
static CURLMcode curl_multi_setopt_pointer(CURLM *handle, CURLMoption option, void *parameter) {
  return curl_multi_setopt(handle, option, parameter);
}
static CURLMcode curl_multi_fdset_pointer(CURLM *handle,
                            void *read_fd_set,
                            void *write_fd_set,
                            void *exc_fd_set,
                            int *max_fd)
{
  return curl_multi_fdset(handle, read_fd_set, write_fd_set, exc_fd_set, max_fd);
}                            
static CURLMsg *curl_multi_info_read_pointer(CURLM *handle, int *msgs_in_queue)
{
  return curl_multi_info_read(handle, msgs_in_queue);
}                            
*/
import "C"

import (
		"unsafe"
		"syscall"
)

type CurlMultiError C.CURLMcode
type CurlMultiMsg	C.CURLMSG

func (e CurlMultiError) Error() string {
	// ret is const char*, no need to free
	ret := C.curl_multi_strerror(C.CURLMcode(e))
	return C.GoString(ret)
}

func newCurlMultiError(errno C.CURLMcode) error {
	// cannot use C.CURLM_OK here, cause multi.h use a undefined emum num
	if errno == 0 { // if nothing wrong
		return nil
	}
	return CurlMultiError(errno)
}

func newCURLMessage(message *C.CURLMsg) (msg *CURLMessage){
	if message == nil {
		return nil
	}
	msg = new(CURLMessage)
	msg.Msg = CurlMultiMsg(message.msg)
	msg.Easy_handle = &CURL{handle: message.easy_handle}
	msg.Data = message.data
	return msg 
}

type CURLM struct {
	handle unsafe.Pointer
}

var dummy unsafe.Pointer
type CURLMessage struct {
	Msg CurlMultiMsg
	Easy_handle *CURL
	Data [unsafe.Sizeof(dummy)]byte
}

// curl_multi_init - create a multi handle
func MultiInit() *CURLM {
	p := C.curl_multi_init()
	return &CURLM{p}
}

// curl_multi_cleanup - close down a multi session
func (mcurl *CURLM) Cleanup() error {
	p := mcurl.handle
	return newCurlMultiError(C.curl_multi_cleanup(p))
}

// curl_multi_perform - reads/writes available data from each easy handle
func (mcurl *CURLM) Perform() (int, error) {
	p := mcurl.handle
	running_handles := C.int(-1)
	err := newCurlMultiError(C.curl_multi_perform(p, &running_handles))
	return int(running_handles), err
}

// curl_multi_add_handle - add an easy handle to a multi session
func (mcurl *CURLM) AddHandle(easy *CURL) error {
	mp := mcurl.handle
	easy_handle := easy.handle
	return newCurlMultiError(C.curl_multi_add_handle(mp, easy_handle))
}

// curl_multi_remove_handle - remove an easy handle from a multi session
func (mcurl *CURLM) RemoveHandle(easy *CURL) error {
	mp := mcurl.handle
	easy_handle := easy.handle
	return newCurlMultiError(C.curl_multi_remove_handle(mp, easy_handle))
}

func (mcurl *CURLM) Timeout() (int, error) {
	p := mcurl.handle
	timeout := C.long(-1)
	err := newCurlMultiError(C.curl_multi_timeout(p, &timeout))
	return int(timeout), err
}

func (mcurl *CURLM) Setopt(opt int, param interface{}) error {
	p := mcurl.handle
	if param == nil {
		return newCurlMultiError(C.curl_multi_setopt_pointer(p, C.CURLMoption(opt), nil))
	}
	switch {
	//  currently cannot support these option
	//	case MOPT_SOCKETFUNCTION, MOPT_SOCKETDATA, MOPT_TIMERFUNCTION, MOPT_TIMERDATA:
	//		panic("not supported CURLM.Setopt opt")
	case opt >= C.CURLOPTTYPE_LONG:
		val := C.long(0)
		switch t := param.(type) {
		case int:
			val := C.long(t)
			return newCurlMultiError(C.curl_multi_setopt_long(p, C.CURLMoption(opt), val))
		case bool:
			val = C.long(0)
			if t {
				val = C.long(1)
			}
			return newCurlMultiError(C.curl_multi_setopt_long(p, C.CURLMoption(opt), val))
		}
	}
	panic("not supported CURLM.Setopt opt or param")
	return nil
}

func (mcurl *CURLM) Fdset(rset, wset, eset *syscall.FdSet) (int, error) {
	p := mcurl.handle
	read := unsafe.Pointer(rset)
	write := unsafe.Pointer(wset)
	exc := unsafe.Pointer(eset)
	maxfd := C.int(-1)
	err := newCurlMultiError(C.curl_multi_fdset_pointer(p, read, write,
							 exc, &maxfd))
	return int(maxfd), err
}

func (mcurl *CURLM) Info_read() (*CURLMessage, int) {
	p := mcurl.handle
	left := C.int(0)
  	return newCURLMessage(C.curl_multi_info_read_pointer(p, &left)), int(left)
}
