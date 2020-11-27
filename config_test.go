package cors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsAllowedAllOrigins(t *testing.T) {
	a := assert.New(t)

	c := &Cors{}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedHeaders: []string{"*"}}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedOriginsPatterns: []string{".*"}}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedOriginsPatterns: []string{".*"}}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedHeaders: []string{"Host"}}
	a.False(c.IsAllowedAllHeaders())

	c = &Cors{AllowedOriginsPatterns: []string{"Host"}}
	a.False(c.IsAllowedAllOrigins())
}

func Test_IsAllowedAllMethods(t *testing.T) {
	a := assert.New(t)

	c := &Cors{}
	a.True(c.IsAllowedAllMethods())

	c = &Cors{AllowedMethods: []string{"*"}}
	a.True(c.IsAllowedAllMethods())

	c = &Cors{AllowedMethods: []string{"GET"}}
	a.False(c.IsAllowedAllMethods())
}

func Test_IsAllowedAllHeaders(t *testing.T) {
	a := assert.New(t)

	c := &Cors{}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedHeaders: []string{"*"}}
	a.True(c.IsAllowedAllHeaders())

	c = &Cors{AllowedHeaders: []string{"Host"}}
	a.False(c.IsAllowedAllHeaders())
}
