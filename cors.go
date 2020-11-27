package cors

import (
	"github.com/hughcube-go/utils/msslice"
	"net/http"
	url "net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	HeaderKeyVary   = "Vary"
	HeaderKeyOrigin = "Origin"

	HeaderKeyAccessControlRequestMethod  = "Access-Control-Request-Method"
	HeaderKeyAccessControlAllowOrigin    = "Access-Control-Allow-Origin"
	HeaderKeyAccessControlMaxAge         = "Access-Control-Max-Age"
	HeaderKeyAccessControlExposeHeaders  = "Access-Control-Expose-Headers"
	HeaderKeyAccessControlRequestHeaders = "Access-Control-Request-Headers"

	HeaderKeyAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderKeyAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderKeyAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

func (c *Cors) Handler(writer http.ResponseWriter, request *http.Request) bool {
	if c.IsPreflightRequest(request) {
		c.HandlePreflightRequest(request, writer)
		c.VaryHeader(writer, HeaderKeyAccessControlRequestMethod)
		return false
	}

	if IsRequestMethod(request, http.MethodOptions) {
		c.VaryHeader(writer, HeaderKeyAccessControlRequestMethod)
	}
	c.AddActualRequestHeaders(request, writer)
	return true
}

func (c *Cors) IsActualRequestAllowed(request *http.Request) bool {
	return c.isOriginAllowed(request)
}

func (c *Cors) IsCorsRequest(request *http.Request) bool {
	return HasRequestHeader(request, HeaderKeyOrigin) &&
		!c.isSameHost(request)
}

func (c *Cors) IsPreflightRequest(request *http.Request) bool {
	return IsRequestMethod(request, http.MethodOptions) &&
		HasRequestHeader(request, HeaderKeyAccessControlRequestMethod)
}

func (c *Cors) VaryHeader(writer http.ResponseWriter, name string) {
	if !HasResponseHeader(writer, HeaderKeyVary) {
		SetResponseHeader(writer, HeaderKeyVary, name)
	} else if !msslice.InArray(name, strings.Split(GetResponseHeader(writer, HeaderKeyVary), ", ")) {
		SetResponseHeader(writer, HeaderKeyVary, GetResponseHeader(writer, HeaderKeyVary)+", "+name)
	}
}

func (c *Cors) HandlePreflightRequest(request *http.Request, writer http.ResponseWriter) {
	SetResponseStatusCode(writer, http.StatusNoContent)
	c.AddPreflightRequestHeaders(request, writer)
}

func (c *Cors) AddActualRequestHeaders(request *http.Request, writer http.ResponseWriter) {
	c.configureAllowedOrigin(request, writer)

	if HasResponseHeader(writer, HeaderKeyAccessControlAllowOrigin) {
		c.configureAllowCredentials(request, writer)
		c.configureExposedHeaders(request, writer)
	}
}

func (c *Cors) AddPreflightRequestHeaders(request *http.Request, writer http.ResponseWriter) {
	c.configureAllowedOrigin(request, writer)

	if HasResponseHeader(writer, HeaderKeyAccessControlAllowOrigin) {
		c.configureAllowCredentials(request, writer)
		c.configureAllowedMethods(request, writer)
		c.configureAllowedHeaders(request, writer)
		c.configureMaxAge(request, writer)
	}
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func (c *Cors) isOriginAllowed(request *http.Request) bool {
	if c.IsAllowedAllOrigins() {
		return true
	}

	if !HasRequestHeader(request, HeaderKeyOrigin) {
		return false
	}

	origin := GetRequestHeader(request, HeaderKeyOrigin)

	for _, o := range c.AllowedOrigins {
		if o == origin {
			return true
		}
	}

	for _, p := range c.AllowedOriginsPatterns {
		isMatch, matchError := regexp.MatchString(p, origin)
		if nil == matchError && isMatch {
			return true
		}
	}

	return false
}

func (c *Cors) configureAllowedOrigin(request *http.Request, writer http.ResponseWriter) {
	if c.IsAllowedAllOrigins() && !c.SupportsCredentials {
		// Safe+cacheable, allow everything
		SetResponseHeader(writer, HeaderKeyAccessControlAllowOrigin, "*")
	} else if c.isSingleOriginAllowed() {
		// Single origins can be safely set
		SetResponseHeader(writer, HeaderKeyAccessControlAllowOrigin, c.AllowedOrigins[0])
	} else {
		// For dynamic headers, check the origin first
		origin := GetRequestHeader(request, HeaderKeyOrigin)
		if c.isOriginAllowed(request) && 0 < len(origin) {
			SetResponseHeader(writer, HeaderKeyAccessControlAllowOrigin, origin)
		}

		c.VaryHeader(writer, HeaderKeyOrigin)
	}
}

func (c *Cors) isSingleOriginAllowed() bool {
	if c.IsAllowedAllOrigins() || 0 < len(c.AllowedOriginsPatterns) {
		return false
	}

	return 1 == len(c.AllowedOrigins)
}

func (c *Cors) configureAllowedMethods(request *http.Request, writer http.ResponseWriter) {
	var allowMethods = c.AllowedMethods

	if c.IsAllowedAllMethods() {
		allowMethods = strings.Split(GetRequestHeader(request, HeaderKeyAccessControlRequestMethod), ", ")
		c.VaryHeader(writer, HeaderKeyAccessControlRequestMethod)
	}

	for k, v := range allowMethods {
		allowMethods[k] = strings.ToTitle(v)
	}

	if 0 < len(allowMethods) {
		SetResponseHeader(writer, HeaderKeyAccessControlAllowMethods, strings.Join(allowMethods, ", "))
	}
}

func (c *Cors) configureAllowedHeaders(request *http.Request, writer http.ResponseWriter) {
	var allowHeaders = c.AllowedHeaders

	if c.IsAllowedAllHeaders() {
		allowHeaders = strings.Split(GetRequestHeader(request, HeaderKeyAccessControlRequestHeaders), ", ")
		c.VaryHeader(writer, HeaderKeyAccessControlRequestHeaders)
	}

	if 0 < len(allowHeaders) {
		SetResponseHeader(writer, HeaderKeyAccessControlAllowHeaders, strings.Join(allowHeaders, ", "))
	}
}

func (c Cors) configureAllowCredentials(request *http.Request, writer http.ResponseWriter) {
	if c.SupportsCredentials {
		SetResponseHeader(writer, HeaderKeyAccessControlAllowCredentials, "true")
	}
}

func (c *Cors) configureExposedHeaders(request *http.Request, writer http.ResponseWriter) {
	if 0 < len(c.ExposedHeaders) {
		SetResponseHeader(writer, HeaderKeyAccessControlExposeHeaders, strings.Join(c.ExposedHeaders, ", "))
	}
}

func (c *Cors) configureMaxAge(request *http.Request, writer http.ResponseWriter) {
	if 0 != c.MaxAge {
		SetResponseHeader(writer, HeaderKeyAccessControlMaxAge, strconv.FormatInt(c.MaxAge, 10))
	}
}

func (c *Cors) isSameHost(request *http.Request) bool {
	origin := GetRequestHeader(request, HeaderKeyOrigin)

	originUrl, err := url.Parse(origin)
	if err != nil {
		return false
	}

	return originUrl.Host == GetRequestHost(request)
}
