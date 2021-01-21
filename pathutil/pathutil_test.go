package pathutil

import (
	"os"
	"os/user"
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestExpandUser(t *testing.T) {
	assert := require.New(t)
	var v string
	var err error

	u, _ := user.Current()

	v, err = ExpandUser(`/dev/null`)
	assert.Equal(v, `/dev/null`)
	assert.Nil(err)

	v, err = ExpandUser(`~`)
	assert.Equal(v, u.HomeDir)
	assert.Nil(err)

	v, err = ExpandUser("~" + u.Name)
	assert.Equal(v, u.HomeDir)
	assert.Nil(err)

	v, err = ExpandUser("~/test-123")
	assert.Equal(v, u.HomeDir+"/test-123")
	assert.Nil(err)

	v, err = ExpandUser("~" + u.Name + "/test-123")
	assert.Equal(v, u.HomeDir+"/test-123")
	assert.Nil(err)

	v, err = ExpandUser("~/test-123/~/123")
	assert.Equal(v, u.HomeDir+"/test-123/~/123")
	assert.Nil(err)

	v, err = ExpandUser("~" + u.Name + "/test-123/~" + u.Name + "/123")
	assert.Equal(v, u.HomeDir+"/test-123/~"+u.Name+"/123")
	assert.Nil(err)

	assert.False(IsNonemptyFile(`/nonexistent.txt`))
	assert.False(IsNonemptyDir(`/nonexistent/dir`))

	assert.True(IsNonemptyFile(`/etc/hosts`))
	assert.True(IsNonemptyDir(`/etc`))

	x, err := os.Executable()
	assert.NoError(err)
	assert.True(IsNonemptyExecutableFile(x))
}
