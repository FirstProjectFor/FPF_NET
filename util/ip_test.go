package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPublicIP(t *testing.T) {
	assert.False(t, IsPublicIP("127"))
	assert.False(t, IsPublicIP("127.1.1"))
	assert.False(t, IsPublicIP("127.0.0.1"))
	assert.False(t, IsPublicIP("192.168.0.201"))

	assert.True(t, IsPublicIP("103.235.46.40"))
}
