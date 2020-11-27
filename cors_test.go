package cors

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpec(t *testing.T) {
	cases := []struct {
		name       string
		cors       *Cors
		method     string
		reqHeaders map[string]string
		resHeaders map[string]string
		resCode    int
		resBody    string
	}{
		{
			name:       "no config",
			cors:       &Cors{},
			method:     http.MethodGet,
			reqHeaders: map[string]string{},
			resHeaders: map[string]string{},
			resCode:    http.StatusOK,
			resBody:    "ok",
		},

		{
			name:   "PreflightRequest",
			cors:   &Cors{},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Vary":                         "Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Authorization,x-Ccookie",
			},
			resCode: http.StatusNoContent,
		},

		{
			name: "PreflightRequest",
			cors: &Cors{
				AllowedHeaders: []string{"Authorization"},
			},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Vary":                         "Access-Control-Request-Method",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Authorization",
			},
			resCode: http.StatusNoContent,
		},

		{
			name: "PreflightRequest",
			cors: &Cors{
				AllowedMethods: []string{"POST"},
			},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Vary":                         "Access-Control-Request-Headers, Access-Control-Request-Method",
				"Access-Control-Allow-Methods": "POST",
				"Access-Control-Allow-Headers": "Authorization,x-Ccookie",
			},
			resCode: http.StatusNoContent,
		},

		{
			name: "PreflightRequest",
			cors: &Cors{
				AllowedOrigins: []string{"http://example1.com/foo"},
			},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "http://example1.com/foo",
				"Vary":                         "Access-Control-Request-Method, Access-Control-Request-Headers",
				"Access-Control-Allow-Methods": "GET",
				"Access-Control-Allow-Headers": "Authorization,x-Ccookie",
			},
			resCode: http.StatusNoContent,
		},

		{
			name: "PreflightRequest",
			cors: &Cors{
				AllowedOriginsPatterns: []string{"^http://example1.*"},
			},
			method: http.MethodOptions,
			reqHeaders: map[string]string{
				"Origin":                         "http://example.com/foo",
				"Access-Control-Request-Method":  "GET",
				"Access-Control-Request-Headers": "Authorization,x-Ccookie",
			},
			resHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "",
				"Vary":                         "Origin, Access-Control-Request-Method",
				"Access-Control-Allow-Methods": "",
				"Access-Control-Allow-Headers": "",
			},
			resCode: http.StatusNoContent,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "http://example.com/foo", nil)

			for name, value := range tc.reqHeaders {
				req.Header.Add(name, value)
			}

			t.Run("Handler", func(t *testing.T) {
				res := httptest.NewRecorder()
				getHandlerFun(tc.cors).ServeHTTP(res, req)
				assertResponse(t, res, tc.resCode, tc.resHeaders, tc.resBody)
			})
		})
	}
}

func getHandlerFun(cors *Cors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cors.Handler(w, r) {
			_, _ = w.Write([]byte("ok"))
		}
	})
}

func assertResponse(
	t *testing.T,
	res *httptest.ResponseRecorder,
	httpCode int,
	headers map[string]string,
	body string,
) {
	a := assert.New(t)

	a.Equal(res.Code, httpCode)

	for k, v := range headers {
		a.Equal(res.Header().Get(k), v)
	}

	a.Equal(res.Body.String(), body)
}
