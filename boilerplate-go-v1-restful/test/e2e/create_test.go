package e2e_test

import (
	"codegen/utils/random"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)
}

func TestCreateNOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	req.Name = ""

	_, err := c.Create(req)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "request body invalid")
}

func TestCreateLongNameBadRequest(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	req.Name = "e2e_test_" + random.Str(100)

	_, err := c.Create(req)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "request body invalid")
}
