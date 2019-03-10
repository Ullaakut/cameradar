package curl

/*
#include <stdlib.h>
#include <curl/curl.h>
#include "callback.h"
#include "compat.h"

static CURLcode curl_easy_setopt_long(CURL *handle, CURLoption option, long parameter) {
  return curl_easy_setopt(handle, option, parameter);
}
static CURLcode curl_easy_setopt_string(CURL *handle, CURLoption option, char *parameter) {
  return curl_easy_setopt(handle, option, parameter);
}
static CURLcode curl_easy_setopt_slist(CURL *handle, CURLoption option, struct curl_slist *parameter) {
  return curl_easy_setopt(handle, option, parameter);
}
static CURLcode curl_easy_setopt_pointer(CURL *handle, CURLoption option, void *parameter) {
  return curl_easy_setopt(handle, option, parameter);
}
static CURLcode curl_easy_setopt_off_t(CURL *handle, CURLoption option, off_t parameter) {
  return curl_easy_setopt(handle, option, parameter);
}

static CURLcode curl_easy_getinfo_string(CURL *curl, CURLINFO info, char **p) {
 return curl_easy_getinfo(curl, info, p);
}
static CURLcode curl_easy_getinfo_long(CURL *curl, CURLINFO info, long *p) {
 return curl_easy_getinfo(curl, info, p);
}
static CURLcode curl_easy_getinfo_double(CURL *curl, CURLINFO info, double *p) {
 return curl_easy_getinfo(curl, info, p);
}
static CURLcode curl_easy_getinfo_slist(CURL *curl, CURLINFO info, struct curl_slist **p) {
 return curl_easy_getinfo(curl, info, p);
}

static CURLFORMcode curl_formadd_name_content_length(
    struct curl_httppost **httppost, struct curl_httppost **last_post, char *name, char *content, int length) {
    return curl_formadd(httppost, last_post,
                        CURLFORM_COPYNAME, name,
                        CURLFORM_COPYCONTENTS, content,
                        CURLFORM_CONTENTSLENGTH, length, CURLFORM_END);
}
static CURLFORMcode curl_formadd_name_content_length_type(
    struct curl_httppost **httppost, struct curl_httppost **last_post, char *name, char *content, int length, char *type) {
    return curl_formadd(httppost, last_post,
                        CURLFORM_COPYNAME, name,
                        CURLFORM_COPYCONTENTS, content,
                        CURLFORM_CONTENTSLENGTH, length,
                        CURLFORM_CONTENTTYPE, type, CURLFORM_END);
}
static CURLFORMcode curl_formadd_name_file_type(
    struct curl_httppost **httppost, struct curl_httppost **last_post, char *name, char *filename, char *type) {
    return curl_formadd(httppost, last_post,
                        CURLFORM_COPYNAME, name,
                        CURLFORM_FILE, filename,
                        CURLFORM_CONTENTTYPE, type, CURLFORM_END);
}
 // TODO: support multi file

*/
import "C"

import (
	"fmt"
	"mime"
	"path"
	"unsafe"
	"sync"
)

type CurlInfo C.CURLINFO
type CurlError C.CURLcode

type CurlString *C.char

func NewCurlString(s string) CurlString {
	return CurlString(unsafe.Pointer(C.CString(s)))
}

func FreeCurlString(s CurlString) {
	C.free(unsafe.Pointer(s))
}

func (e CurlError) Error() string {
	// ret is const char*, no need to free
	ret := C.curl_easy_strerror(C.CURLcode(e))
	return fmt.Sprintf("curl: %s", C.GoString(ret))
}

func newCurlError(errno C.CURLcode) error {
	if errno == C.CURLE_OK { // if nothing wrong
		return nil
	}
	return CurlError(errno)
}

// curl_easy interface
type CURL struct {
	handle unsafe.Pointer
	// callback functions, bool ret means ok or not
	headerFunction, writeFunction *func([]byte, interface{}) bool
	readFunction                  *func([]byte, interface{}) int // return num of bytes writed to buf
	progressFunction              *func(float64, float64, float64, float64, interface{}) bool
	fnmatchFunction               *func(string, string, interface{}) int
	// callback datas
	headerData, writeData, readData, progressData, fnmatchData interface{}
	// list of C allocs
	mallocAllocs []*C.char
}

// concurrent safe context map
type contextMap struct {
	items map[uintptr]*CURL
	sync.RWMutex
}

func (c *contextMap) Set(k uintptr, v *CURL) {
	c.Lock()
	defer c.Unlock()

	c.items[k] = v
}

func (c *contextMap) Get(k uintptr) *CURL {
	c.RLock()
	defer c.RUnlock()

	return c.items[k]
}

func (c *contextMap) Delete(k uintptr) {
	c.Lock()
	defer c.Unlock()

	delete(c.items, k)
}

var context_map = &contextMap {
	items: make(map[uintptr]*CURL),
}

// curl_easy_init - Start a libcurl easy session
func EasyInit() *CURL {
	p := C.curl_easy_init()
	c := &CURL{handle: p, mallocAllocs: make([]*C.char, 0)} // other field defaults to nil
	context_map.Set(uintptr(p), c)
	return c
}

// curl_easy_duphandle - Clone a libcurl session handle
func (curl *CURL) Duphandle() *CURL {
	p := C.curl_easy_duphandle(curl.handle)
	c := &CURL{handle: p}
	context_map.Set(uintptr(p), c)
	return c
}

// curl_easy_cleanup - End a libcurl easy session
func (curl *CURL) Cleanup() {
	p := curl.handle
	C.curl_easy_cleanup(p)
	curl.MallocFreeAfter(0)
	context_map.Delete(uintptr(p))
}

// curl_easy_setopt - set options for a curl easy handle
// WARNING: a function pointer is &fun, but function addr is reflect.ValueOf(fun).Pointer()
func (curl *CURL) Setopt(opt int, param interface{}) error {
	p := curl.handle
	if param == nil {
		// NOTE: some option will crash program when got a nil param
		return newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), nil))
	}
	switch {
	// not really set
	case opt == OPT_READDATA: // OPT_INFILE
		curl.readData = param
		return nil
	case opt == OPT_PROGRESSDATA:
		curl.progressData = param
		return nil
	case opt == OPT_HEADERDATA: // also known as OPT_WRITEHEADER
		curl.headerData = param
		return nil
	case opt == OPT_WRITEDATA: // OPT_FILE
		curl.writeData = param
		return nil

	case opt == OPT_READFUNCTION:
		fun := param.(func([]byte, interface{}) int)
		curl.readFunction = &fun

		ptr := C.return_read_function()
		if err := newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), ptr)); err == nil {
			return newCurlError(C.curl_easy_setopt_pointer(p, OPT_READDATA, unsafe.Pointer(curl.handle)))
		} else {
			return err
		}

	case opt == OPT_PROGRESSFUNCTION:
		fun := param.(func(float64, float64, float64, float64, interface{}) bool)
		curl.progressFunction = &fun

		ptr := C.return_progress_function()
		if err := newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), ptr)); err == nil {
			return newCurlError(C.curl_easy_setopt_pointer(p, OPT_PROGRESSDATA, unsafe.Pointer(curl.handle)))
		} else {
			return err
		}

	case opt == OPT_HEADERFUNCTION:
		fun := param.(func([]byte, interface{}) bool)
		curl.headerFunction = &fun

		ptr := C.return_header_function()
		if err := newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), ptr)); err == nil {
			return newCurlError(C.curl_easy_setopt_pointer(p, OPT_HEADERDATA, unsafe.Pointer(curl.handle)))
		} else {
			return err
		}

	case opt == OPT_WRITEFUNCTION:
		fun := param.(func([]byte, interface{}) bool)
		curl.writeFunction = &fun

		ptr := C.return_write_function()
		if err := newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), ptr)); err == nil {
			return newCurlError(C.curl_easy_setopt_pointer(p, OPT_WRITEDATA, unsafe.Pointer(curl.handle)))
		} else {
			return err
		}

	// for OPT_HTTPPOST, use struct Form
	case opt == OPT_HTTPPOST:
		post := param.(*Form)
		ptr := post.head
		return newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt), unsafe.Pointer(ptr)))

	case opt >= C.CURLOPTTYPE_OFF_T:
		val := C.off_t(0)
		switch t := param.(type) {
		case int:
			val = C.off_t(t)
		case uint64:
			val = C.off_t(t)
		default:
			panic("OFF_T conversion not supported")
		}
		return newCurlError(C.curl_easy_setopt_off_t(p, C.CURLoption(opt), val))

	case opt >= C.CURLOPTTYPE_FUNCTIONPOINT:
		// function pointer
		panic("function pointer not implemented yet!")

	case opt >= C.CURLOPTTYPE_OBJECTPOINT:
		switch t := param.(type) {
		case string:
			ptr := C.CString(t)
			curl.mallocAddPtr(ptr)
			return newCurlError(C.curl_easy_setopt_string(p, C.CURLoption(opt), ptr))
		case CurlString:
			ptr := (*C.char)(t)
			return newCurlError(C.curl_easy_setopt_string(p, C.CURLoption(opt), ptr))
		case []string:
			if len(t) > 0 {
				ptr := C.CString(t[0])
				curl.mallocAddPtr(ptr)
				a_slist := C.curl_slist_append(nil, ptr)
				for _, s := range t[1:] {
					ptr := C.CString(s)
					curl.mallocAddPtr(ptr)
					a_slist = C.curl_slist_append(a_slist, ptr)
				}
				return newCurlError(C.curl_easy_setopt_slist(p, C.CURLoption(opt), a_slist))
			} else {
				return newCurlError(C.curl_easy_setopt_slist(p, C.CURLoption(opt), nil))
			}
		case []CurlString:
			if len(t) > 0 {
				ptr := (*C.char)(t[0])
				a_slist := C.curl_slist_append(nil, ptr)
				for _, s := range t[1:] {
					ptr := (*C.char)(s)
					a_slist = C.curl_slist_append(a_slist, ptr)
				}
				return newCurlError(C.curl_easy_setopt_slist(p, C.CURLoption(opt), a_slist))
			} else {
				return newCurlError(C.curl_easy_setopt_slist(p, C.CURLoption(opt), nil))
			}
		default:
			// It panics if v's Kind is not Chan, Func, Map, Ptr, Slice, or UnsafePointer.
			// val := reflect.ValueOf(param)
			//fmt.Printf("DEBUG(Setopt): param=%x\n", val.Pointer())
			//println("DEBUG can addr =", val.Pointer(), "opt=", opt)
			// pass a pointer to GoInterface
			return newCurlError(C.curl_easy_setopt_pointer(p, C.CURLoption(opt),
				unsafe.Pointer(&param)))
		}
	case opt >= C.CURLOPTTYPE_LONG:
		val := C.long(0)
		switch t := param.(type) {
		case int:
			val = C.long(t)
		case bool:
			if t {
				val = 1
			}
		case int64:
			val = C.long(t)
		case int32:
			val = C.long(t)
		default:
			panic("not supported converstion to c long")
		}
		return newCurlError(C.curl_easy_setopt_long(p, C.CURLoption(opt), val))
	}
	panic("opt param error!")
}

// curl_easy_send - sends raw data over an "easy" connection
func (curl *CURL) Send(buffer []byte) (int, error) {
	p := curl.handle
	buflen := len(buffer)
	n := C.size_t(0)
	ret := C.curl_easy_send(p, unsafe.Pointer(&buffer[0]), C.size_t(buflen), &n)
	return int(n), newCurlError(ret)
}

// curl_easy_recv - receives raw data on an "easy" connection
func (curl *CURL) Recv(buffer []byte) (int, error) {
	p := curl.handle
	buflen := len(buffer)
	buf := C.CString(string(buffer))
	n := C.size_t(0)
	ret := C.curl_easy_recv(p, unsafe.Pointer(buf), C.size_t(buflen), &n)
	return copy(buffer, C.GoStringN(buf, C.int(n))), newCurlError(ret)
}

// curl_easy_perform - Perform a file transfer
func (curl *CURL) Perform() error {
	p := curl.handle
	return newCurlError(C.curl_easy_perform(p))
}

// curl_easy_pause - pause and unpause a connection
func (curl *CURL) Pause(bitmask int) error {
	p := curl.handle
	return newCurlError(C.curl_easy_pause(p, C.int(bitmask)))
}

// curl_easy_reset - reset all options of a libcurl session handle
func (curl *CURL) Reset() {
	p := curl.handle
	C.curl_easy_reset(p)
}

// curl_easy_escape - URL encodes the given string
func (curl *CURL) Escape(url string) string {
	p := curl.handle
	oldUrl := C.CString(url)
	defer C.free(unsafe.Pointer(oldUrl))
	newUrl := C.curl_easy_escape(p, oldUrl, 0)
	defer C.curl_free(unsafe.Pointer(newUrl))
	return C.GoString(newUrl)
}

// curl_easy_unescape - URL decodes the given string
func (curl *CURL) Unescape(url string) string {
	p := curl.handle
	oldUrl := C.CString(url)
	outlength := C.int(0)
	defer C.free(unsafe.Pointer(oldUrl))
	// If outlength is non-NULL, the function will write the length of the
	// returned string in  the  integer  it  points  to.  This allows an
	// escaped string containing %00 to still get used properly after unescaping.
	newUrl := C.curl_easy_unescape(p, oldUrl, 0, &outlength)
	defer C.curl_free(unsafe.Pointer(newUrl))
	return C.GoStringN(newUrl, outlength)
}

// curl_easy_getinfo - extract information from a curl handle
func (curl *CURL) Getinfo(info CurlInfo) (ret interface{}, err error) {
	p := curl.handle
	cInfo := C.CURLINFO(info)
	switch cInfo & C.CURLINFO_TYPEMASK {
	case C.CURLINFO_STRING:
		a_string := C.CString("")
		defer C.free(unsafe.Pointer(a_string))
		err := newCurlError(C.curl_easy_getinfo_string(p, cInfo, &a_string))
		ret := C.GoString(a_string)
		debugf("Getinfo %s", ret)
		return ret, err
	case C.CURLINFO_LONG:
		a_long := C.long(-1)
		err := newCurlError(C.curl_easy_getinfo_long(p, cInfo, &a_long))
		ret := int(a_long)
		debugf("Getinfo %s", ret)
		return ret, err
	case C.CURLINFO_DOUBLE:
		a_double := C.double(0.0)
		err := newCurlError(C.curl_easy_getinfo_double(p, cInfo, &a_double))
		ret := float64(a_double)
		debugf("Getinfo %s", ret)
		return ret, err
	case C.CURLINFO_SLIST:
		a_ptr_slist := new(_Ctype_struct_curl_slist)
		err := newCurlError(C.curl_easy_getinfo_slist(p, cInfo, &a_ptr_slist))
		ret := []string{}
		for a_ptr_slist != nil {
			debugf("Getinfo %s %v", C.GoString(a_ptr_slist.data), a_ptr_slist.next)
			ret = append(ret, C.GoString(a_ptr_slist.data))
			a_ptr_slist = a_ptr_slist.next
		}
		return ret, err
	default:
		panic("error calling Getinfo\n")
	}
	panic("not implemented yet!")
	return nil, nil
}

func (curl *CURL) GetHandle() unsafe.Pointer {
	return curl.handle
}

func (curl *CURL) MallocGetPos() int {
	return len(curl.mallocAllocs)
}

func (curl *CURL) MallocFreeAfter(from int) {
	l := len(curl.mallocAllocs)
	for idx := from; idx < l; idx++ {
		C.free(unsafe.Pointer(curl.mallocAllocs[idx]))
		curl.mallocAllocs[idx] = nil
	}
	curl.mallocAllocs = curl.mallocAllocs[0:from]
}

func (curl *CURL) mallocAddPtr(ptr *C.char) {
	curl.mallocAllocs = append(curl.mallocAllocs, ptr)
}

// A multipart/formdata HTTP POST form
type Form struct {
	head, last *C.struct_curl_httppost
}

func NewForm() *Form {
	return &Form{}
}

func (form *Form) Add(name string, content interface{}) error {
	head, last := form.head, form.last
	namestr := C.CString(name)
	defer C.free(unsafe.Pointer(namestr))
	var (
		buffer *C.char
		length C.int
	)
	switch t := content.(type) {
	case string:
		buffer = C.CString(t)
		length = C.int(len(t))
	case []byte:
		buffer = C.CString(string(t))
		length = C.int(len(t))
	default:
		panic("not implemented")
	}
	defer C.free(unsafe.Pointer(buffer))
	C.curl_formadd_name_content_length(&head, &last, namestr, buffer, length)
	form.head, form.last = head, last
	return nil
}

func (form *Form) AddWithType(name string, content interface{}, content_type string) error {
	head, last := form.head, form.last
	namestr := C.CString(name)
	typestr := C.CString(content_type)
	defer C.free(unsafe.Pointer(namestr))
	defer C.free(unsafe.Pointer(typestr))
	var (
		buffer *C.char
		length C.int
	)
	switch t := content.(type) {
	case string:
		buffer = C.CString(t)
		length = C.int(len(t))
	case []byte:
		buffer = C.CString(string(t))
		length = C.int(len(t))
	default:
		panic("not implemented")
	}
	defer C.free(unsafe.Pointer(buffer))
	C.curl_formadd_name_content_length_type(&head, &last, namestr, buffer, length, typestr)
	form.head, form.last = head, last
	return nil
}

func (form *Form) AddFile(name, filename string) error {
	head, last := form.head, form.last
	namestr := C.CString(name)
	pathstr := C.CString(filename)
	typestr := C.CString(guessType(filename))
	defer C.free(unsafe.Pointer(namestr))
	defer C.free(unsafe.Pointer(pathstr))
	defer C.free(unsafe.Pointer(typestr))
	C.curl_formadd_name_file_type(&head, &last, namestr, pathstr, typestr)
	form.head, form.last = head, last
	return nil
}

func (form *Form) AddFromFile(name, filename string) {
}

func guessType(filename string) string {
	ext := path.Ext(filename)
	file_type := mime.TypeByExtension(ext)
	if file_type == "" {
		return "application/octet-stream"
	}
	return file_type
}
