package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSingleOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)

	response, err := c.GetSingle(created.Id)
	require.NoError(t, err)
	assert.Equal(t, created.Id, response.Product.Id)
}

func TestGetSingleNotFound(t *testing.T) {
	c := NewClient(BaseAddr)

	_, err := c.GetSingle(uuid.Nil)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.Code)
	assert.Equal(t, "record not found", apiErr.Message)
}
