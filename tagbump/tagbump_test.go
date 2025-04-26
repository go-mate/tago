package tagbump

import (
	"testing"

	"github.com/go-xlan/gitgo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/runpath"
)

func TestGetGitTags(t *testing.T) {
	gcm := gitgo.New(runpath.PARENT.Path()).WithDebug()

	tags, err := gcm.SortedGitTags()
	require.NoError(t, err)
	t.Log(neatjsons.S(tags))
}
