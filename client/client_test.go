package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinParams(t *testing.T) {
	assert.Equal(t, joinParams([]string{}), "")
	assert.Equal(t, joinParams([]string{"a"}), "")
	assert.Equal(t, joinParams([]string{"k", "v"}), "k=v")
	assert.Equal(t, joinParams([]string{"k", "v", "t"}), "k=v")
	assert.Equal(t, joinParams([]string{"k", "v", "t", "x"}), "k=v&t=x")
}
