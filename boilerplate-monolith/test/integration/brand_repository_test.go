package repository_test

import (
	"codegen/internal/entity"
	"codegen/internal/repository"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrandRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Clean up before test (optional - depends on test strategy)
	_, err := db.Pool().Exec(ctx, "DELETE FROM catalog.brands")
	require.NoError(t, err)

	// 1. Insert
	newBrand := &entity.Brand{
		Name:      "Test Brand",
		Slug:      "test-brand-" + uuid.New().String(),
		Logo:      "logo.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	err = repo.Insert(ctx, newBrand)
	require.NoError(t, err)
	assert.NotZero(t, newBrand.Id)

	// 2. GetById
	fetched, err := repo.GetById(ctx, newBrand.Id)
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, newBrand.Name, fetched.Name)
	assert.Equal(t, newBrand.CreatedBy, fetched.CreatedBy)

	// 3. Update
	newBrand.Name = "Updated Brand"
	updatedBy := uuid.New()
	updatedAt := time.Now()
	newBrand.UpdatedBy = &updatedBy
	newBrand.UpdatedAt = &updatedAt

	affected, err := repo.Update(ctx, newBrand)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	fetchedUpdated, err := repo.GetById(ctx, newBrand.Id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Brand", fetchedUpdated.Name)
	assert.Equal(t, updatedBy, *fetchedUpdated.UpdatedBy)

	// 4. GetAll
	list, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, b := range list {
		if b.Id == newBrand.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "Newly created brand should be in GetAll list")

	// 5. Delete
	affected, err = repo.Delete(ctx, newBrand.Id)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	fetchedDeleted, err := repo.GetById(ctx, newBrand.Id)
	require.NoError(t, err)
	assert.Nil(t, fetchedDeleted)
}
