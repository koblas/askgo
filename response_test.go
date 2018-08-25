package askgo_test

import (
	"testing"

	"github.com/koblas/askgo"
	"github.com/stretchr/testify/require"
)

func Test_Interface(t *testing.T) {
	env := &askgo.ResponseEnvelope{}

	env.WithShouldEndSession(true)

	require.True(t, env.Response.ShouldSessionEnd, "Session End")
}
