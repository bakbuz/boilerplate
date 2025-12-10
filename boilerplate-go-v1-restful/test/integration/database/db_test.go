package database_test

import (
	"testing"

	"codegen/internal/database"

	"github.com/stretchr/testify/require"
)

func TestDatabase_Ping(t *testing.T) {
	c := RequireConfig(t)
	db, err := database.New(c.DataSources.Default)
	require.NoError(t, err)
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err)
}

func TestDatabase_GetString(t *testing.T) {
	c := RequireConfig(t)
	db, err := database.New(c.DataSources.Default)
	require.NoError(t, err)
	defer db.Close()

	str, err := db.GetString("SELECT [name] FROM Products WHERE id = 1")
	require.NoError(t, err)

	require.NotEmpty(t, str)
}
