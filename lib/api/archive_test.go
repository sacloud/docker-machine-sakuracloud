package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUbuntuArchiveID(t *testing.T) {
	id, err := client.GetUbuntuArchiveID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	t.Logf("ubuntu archive ID : %s", id)
}
