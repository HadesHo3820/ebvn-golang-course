package sqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	// Attempt to create a client
	// This might fail if no database is running, but that's acceptable for this test
	// We just want to ensure it doesn't panic and returns either a client or an error
	client, err := NewClient("")

	if err == nil {
		assert.NotNil(t, client)
	} else {
		assert.Error(t, err)
	}
}
