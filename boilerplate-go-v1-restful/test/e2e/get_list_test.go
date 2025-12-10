package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetListOK(t *testing.T) {
	c := NewClient(BaseAddr)

	req := sampleReq()
	created, err := c.Create(req)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, created.Id)

	response, err := c.GetList()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, 1, response.Count)
}
