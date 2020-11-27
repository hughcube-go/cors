package cors

import (
	"fmt"
	"github.com/hughcube-go/utils/mshash"
	"github.com/stretchr/testify/assert"
	http2 "github.com/stretchr/testify/http"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_GetRequestMethod(t *testing.T) {
	a := assert.New(t)

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, v := range methods {
		request, err := http.NewRequest(v, "http://example.com/", nil)
		a.Nil(err)
		a.Equal(GetRequestMethod(request), v)
	}
}

func Test_IsRequestMethod(t *testing.T) {
	a := assert.New(t)

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, v := range methods {
		request, err := http.NewRequest(v, "http://example.com/", nil)
		a.Nil(err)
		a.True(IsRequestMethod(request, v))
	}
}

func Test_GetRequestHeader(t *testing.T) {
	a := assert.New(t)

	request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	a.Nil(err)

	for i := 0; i <= 100; i++ {
		name := randName()
		request.Header.Set(name, name)
		a.Equal(GetRequestHeader(request, name), name)
	}

	name := randName()
	a.Equal(GetRequestHeader(request, name), "")
}

func Test_HasRequestHeader(t *testing.T) {
	a := assert.New(t)

	request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	a.Nil(err)

	for i := 0; i <= 100; i++ {
		name := randName()
		request.Header.Set(name, name)
		a.True(HasRequestHeader(request, name))
	}

	name := randName()
	a.False(HasRequestHeader(request, name))
}

func Test_GetRequestHost(t *testing.T) {
	a := assert.New(t)

	for i := 0; i <= 100; i++ {
		host := fmt.Sprintf("%s.com", strconv.FormatInt(time.Now().UnixNano(), 10))

		request, err := http.NewRequest(http.MethodGet, "http://"+host, nil)
		a.Nil(err)
		a.Equal(GetRequestHost(request), host)
	}
}

func Test_GetResponseHeader(t *testing.T) {
	a := assert.New(t)

	writer := new(http2.TestResponseWriter)

	for i := 0; i <= 100; i++ {
		name := randName()
		writer.Header().Set(name, name)
		a.Equal(GetResponseHeader(writer, name), name)
	}

	name := randName()
	a.NotEqual(GetResponseHeader(writer, name), name)
}

func Test_HasResponseHeader(t *testing.T) {
	a := assert.New(t)

	writer := new(http2.TestResponseWriter)

	for i := 0; i <= 100; i++ {
		name := randName()
		writer.Header().Set(name, name)
		a.True(HasResponseHeader(writer, name))
	}

	name := randName()
	a.False(HasResponseHeader(writer, name))
}

func Test_SetResponseHeader(t *testing.T) {
	a := assert.New(t)

	writer := new(http2.TestResponseWriter)

	for i := 0; i <= 100; i++ {
		name := randName()
		SetResponseHeader(writer, name, name)
		a.Equal(writer.Header().Get(name), name)
	}

	name := randName()
	a.NotEqual(writer.Header().Get(name), name)
}

func Test_SetResponseStatusCode(t *testing.T) {
	a := assert.New(t)

	writer := new(http2.TestResponseWriter)

	for i := 0; i <= 100; i++ {
		SetResponseStatusCode(writer, i)
		a.Equal(writer.StatusCode, i)
	}
}

var randValue int64 = 0

func randName() string {
	randValue++
	return mshash.MD5(strconv.FormatInt(randValue, 10))
}
