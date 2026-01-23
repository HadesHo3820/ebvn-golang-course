package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Test successful client creation
	client, err := NewClient("")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	// We don't verify connection here as we don't have a real Redis instance in unit tests
}
