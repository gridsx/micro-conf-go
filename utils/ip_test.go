package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	fmt.Println(GetIP())
	assert.True(t, len(GetIP()) > 0)
}
