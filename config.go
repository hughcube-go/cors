package cors

import "github.com/hughcube-go/utils/msslice"

type Cors struct {
	AllowedHeaders         []string
	AllowedMethods         []string
	AllowedOrigins         []string
	AllowedOriginsPatterns []string
	ExposedHeaders         []string
	MaxAge                 int64
	SupportsCredentials    bool
}

func (c Cors) IsAllowedAllOrigins() bool {
	if 0 >= len(c.AllowedOrigins) && 0 >= len(c.AllowedOriginsPatterns) {
		return true
	}

	if msslice.InArray("*", c.AllowedOrigins) {
		return true
	}

	if msslice.InArray("*", c.AllowedOriginsPatterns) || msslice.InArray(".*", c.AllowedOriginsPatterns) {
		return true
	}

	return false
}

func (c Cors) IsAllowedAllMethods() bool {
	return 0 >= len(c.AllowedMethods) || msslice.InArray("*", c.AllowedMethods)
}

func (c Cors) IsAllowedAllHeaders() bool {
	return 0 >= len(c.AllowedHeaders) || msslice.InArray("*", c.AllowedHeaders)
}
