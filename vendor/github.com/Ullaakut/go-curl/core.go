// libcurl go bingding
package curl

/*
#cgo linux pkg-config: libcurl
#cgo darwin LDFLAGS: -lcurl
#cgo windows LDFLAGS: -lcurl
#include <stdlib.h>
#include <curl/curl.h>

static char *string_array_index(char **p, int i) {
  return p[i];
}
*/
import "C"

import (
	"time"
	"unsafe"
)

// curl_global_init - Global libcurl initialisation
func GlobalInit(flags int) error {
	return newCurlError(C.curl_global_init(C.long(flags)))
}

// curl_global_cleanup - global libcurl cleanup
func GlobalCleanup() {
	C.curl_global_cleanup()
}

type VersionInfoData struct {
	Age C.CURLversion
	// age >= 0
	Version       string
	VersionNum    uint
	Host          string
	Features      int
	SslVersion    string
	SslVersionNum int
	LibzVersion   string
	Protocols     []string
	// age >= 1
	Ares    string
	AresNum int
	// age >= 2
	Libidn string
	// age >= 3
	IconvVerNum   int
	LibsshVersion string
}

// curl_version - returns the libcurl version string
func Version() string {
	return C.GoString(C.curl_version())
}

// curl_version_info - returns run-time libcurl version info
func VersionInfo(ver C.CURLversion) *VersionInfoData {
	data := C.curl_version_info(ver)
	ret := new(VersionInfoData)
	ret.Age = data.age
	switch age := ret.Age; {
	case age >= 0:
		ret.Version = string(C.GoString(data.version))
		ret.VersionNum = uint(data.version_num)
		ret.Host = C.GoString(data.host)
		ret.Features = int(data.features)
		ret.SslVersion = C.GoString(data.ssl_version)
		ret.SslVersionNum = int(data.ssl_version_num)
		ret.LibzVersion = C.GoString(data.libz_version)
		// ugly but works
		ret.Protocols = []string{}
		for i := C.int(0); C.string_array_index(data.protocols, i) != nil; i++ {
			p := C.string_array_index(data.protocols, i)
			ret.Protocols = append(ret.Protocols, C.GoString(p))
		}
		fallthrough
	case age >= 1:
		ret.Ares = C.GoString(data.ares)
		ret.AresNum = int(data.ares_num)
		fallthrough
	case age >= 2:
		ret.Libidn = C.GoString(data.libidn)
		fallthrough
	case age >= 3:
		ret.IconvVerNum = int(data.iconv_ver_num)
		ret.LibsshVersion = C.GoString(data.libssh_version)
	}
	return ret
}

// curl_getdate - Convert a date string to number of seconds since January 1, 1970
// In golang, we convert it to a *time.Time
func Getdate(date string) *time.Time {
	datestr := C.CString(date)
	defer C.free(unsafe.Pointer(datestr))
	t := C.curl_getdate(datestr, nil)
	if t == -1 {
		return nil
	}
	unix := time.Unix(int64(t), 0).UTC()
	return &unix

	/*
	   // curl_getenv - return value for environment name
	   func Getenv(name string) string {
	           namestr := C.CString(name)
	           defer C.free(unsafe.Pointer(namestr))
	           ret := C.curl_getenv(unsafe.Pointer(namestr))
	           defer C.free(unsafe.Pointer(ret))

	           return C.GoString(ret)
	   }
	*/
}

// TODO: curl_global_init_mem
