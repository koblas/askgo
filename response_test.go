package askgo_test

import (
	"testing"

	"github.com/koblas/askgo"
)

func Test_Interface(t *testing.T) {
	env := &askgo.ResponseEnvelope{}

	env.WithShouldEndSession(true)

	if !env.Response.ShouldSessionEnd {
		t.Error("Expected true")
	}
}
