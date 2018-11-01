package httputil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryStringModify(t *testing.T) {
	assert := require.New(t)

	u, err := url.Parse(`https://example.com`)
	assert.NoError(err)
	assert.NotNil(u)

	SetQ(u, `test`, false)
	SetQ(u, `test`, true)

	AddQ(u, `test2`, 1)
	AddQ(u, `test2`, 3)

	SetQ(u, `nope`, true)
	DelQ(u, `nope`)

	assert.Equal(u.String(), `https://example.com?test=true&test2=1&test2=3`)
}