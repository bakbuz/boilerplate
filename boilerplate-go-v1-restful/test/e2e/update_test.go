package e2e_test

import (
	"codegen/utils/random"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)

	req.Name = "e2e_test_rename_" + random.Str(4)

	err = c.Update(created.Id, req)
	require.NoError(t, err)
}

func TestUpdateRecordNotFound(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	req.Name = "e2e_test_rename_" + random.Str(4)

	err := c.Update(uuid.Nil, req)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 404, apiErr.Code)
	assert.Equal(t, "record not found", apiErr.Message)
}

func TestUpdateLongNameBadRequest(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)

	req.Name = "e2e_test_rename_" + random.Str(100)

	err = c.Update(created.Id, req)
	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, 400, apiErr.Code)
	assert.Contains(t, apiErr.Message, "request body invalid")
}
