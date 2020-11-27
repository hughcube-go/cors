package cors

import (
	"net/http"
	"strings"
)

// get request method
func GetRequestMethod(request *http.Request) string {
	return strings.ToTitle(request.Method)
}

// check is method request
func IsRequestMethod(request *http.Request, value string) bool {
	return strings.ToTitle(value) == GetRequestMethod(request)
}

// get request header
func GetRequestHeader(request *http.Request, name string) string {
	return request.Header.Get(name)
}

// check has request header
func HasRequestHeader(request *http.Request, name string) bool {
	header := GetRequestHeader(request, name)
	return 0 < len(header)
}

// get request header Host
func GetRequestHost(request *http.Request) string {
	return request.Host
}

/////////////////////////////////////////////////
/////////////////////////////////////////////////
/////////////////////////////////////////////////

// get response header
func GetResponseHeader(writer http.ResponseWriter, name string) string {
	return writer.Header().Get(name)
}

// check has response header
func HasResponseHeader(writer http.ResponseWriter, name string) bool {
	header := GetResponseHeader(writer, name)
	return 0 < len(header)
}

// set response header
func SetResponseHeader(writer http.ResponseWriter, name string, value string) {
	writer.Header().Set(name, value)
}

// set http code
func SetResponseStatusCode(writer http.ResponseWriter, code int) {
	writer.WriteHeader(code)
}
