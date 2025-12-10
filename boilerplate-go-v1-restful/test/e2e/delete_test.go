package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)

	err = c.Delete(created.Id)
	require.NoError(t, err)
}

func TestDeleteNOK(t *testing.T) {
	c := NewClient(BaseAddr)

	err := c.Delete(uuid.Nil)
	require.NoError(t, err)
}
