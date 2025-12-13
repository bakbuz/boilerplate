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

	// 6. DeleteByIds
	brandToDelete1 := &entity.Brand{
		Name:      "DeleteByIds 1",
		Slug:      "del-1-" + uuid.New().String(),
		Logo:      "logo1.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}
	brandToDelete2 := &entity.Brand{
		Name:      "DeleteByIds 2",
		Slug:      "del-2-" + uuid.New().String(),
		Logo:      "logo2.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	err = repo.Insert(ctx, brandToDelete1)
	require.NoError(t, err)
	err = repo.Insert(ctx, brandToDelete2)
	require.NoError(t, err)

	idsToDelete := []int32{brandToDelete1.Id, brandToDelete2.Id}
	deletedCount, err := repo.DeleteByIds(ctx, idsToDelete)
	require.NoError(t, err)
	assert.Equal(t, int64(2), deletedCount)

	// Verify they are gone
	fetched1, err := repo.GetById(ctx, brandToDelete1.Id)
	require.NoError(t, err)
	assert.Nil(t, fetched1)

	fetched2, err := repo.GetById(ctx, brandToDelete2.Id)
	require.NoError(t, err)
	assert.Nil(t, fetched2)

}

func TestBrandRepository_Upsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// 1. Upsert (Insert Scenario)
	upsertID := int32(99999) // Intentionally explicit ID
	upsertBrand := &entity.Brand{
		Id:        upsertID,
		Name:      "Upsert Brand",
		Slug:      "upsert-brand-" + uuid.New().String(),
		Logo:      "upsert-logo.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	// Clean up potential leftover
	_, _ = repo.Delete(ctx, upsertID)

	err := repo.Upsert(ctx, upsertBrand)
	require.NoError(t, err)

	fetchedUpsert, err := repo.GetById(ctx, upsertID)
	require.NoError(t, err)
	assert.NotNil(t, fetchedUpsert)
	assert.Equal(t, upsertBrand.Name, fetchedUpsert.Name)

	// 2. Upsert (Update Scenario)
	upsertBrand.Name = "Upsert Brand Updated"
	updatedByUpsert := uuid.New()
	updatedAtUpsert := time.Now()
	upsertBrand.UpdatedBy = &updatedByUpsert
	upsertBrand.UpdatedAt = &updatedAtUpsert

	err = repo.Upsert(ctx, upsertBrand)
	require.NoError(t, err)

	fetchedUpsertUpdated, err := repo.GetById(ctx, upsertID)
	require.NoError(t, err)
	assert.Equal(t, "Upsert Brand Updated", fetchedUpsertUpdated.Name)
	assert.Equal(t, updatedByUpsert, *fetchedUpsertUpdated.UpdatedBy)

	// Cleanup Upsert
	_, _ = repo.Delete(ctx, upsertID)
}

func TestBrandRepository_BulkInsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// BulkInsert
	bulkList := []*entity.Brand{
		{
			Name:      "Bulk 1",
			Slug:      "bulk-1-" + uuid.New().String(),
			Logo:      "logo1.png",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		},
		{
			Name:      "Bulk 2",
			Slug:      "bulk-2-" + uuid.New().String(),
			Logo:      "logo2.png",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		},
	}

	err := repo.BulkInsert(ctx, bulkList)
	require.NoError(t, err)

	// Verify BulkInsert count
	finalCount, err := repo.Count(ctx)
	require.NoError(t, err)
	// We expect at least 2 records
	assert.GreaterOrEqual(t, finalCount, int64(2))
}

func TestBrandRepository_BulkUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// 1. Prepare data (Insert first)
	brand1 := &entity.Brand{
		Name:      "To Update 1",
		Slug:      "upd-1-" + uuid.New().String(),
		Logo:      "logo1.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}
	brand2 := &entity.Brand{
		Name:      "To Update 2",
		Slug:      "upd-2-" + uuid.New().String(),
		Logo:      "logo2.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}

	// 1. Prepare data (Insert one by one to get IDs)
	err := repo.Insert(ctx, brand1)
	require.NoError(t, err)
	err = repo.Insert(ctx, brand2)
	require.NoError(t, err)

	// 2. Modify objects
	updatedBy := uuid.New()
	updatedAt := time.Now()

	brand1.Name = "Updated Name 1"
	brand1.UpdatedBy = &updatedBy
	brand1.UpdatedAt = &updatedAt

	brand2.Name = "Updated Name 2"
	brand2.UpdatedBy = &updatedBy
	brand2.UpdatedAt = &updatedAt

	// 3. Bulk Update
	err = repo.BulkUpdate(ctx, []*entity.Brand{brand1, brand2})
	require.NoError(t, err)

	// 4. Verify
	fetchedList, err := repo.GetByIds(ctx, []int32{brand1.Id, brand2.Id})
	require.NoError(t, err)
	assert.Len(t, fetchedList, 2)

	for _, b := range fetchedList {
		if b.Id == brand1.Id {
			assert.Equal(t, "Updated Name 1", b.Name)
			assert.Equal(t, updatedBy, *b.UpdatedBy)
		} else if b.Id == brand2.Id {
			assert.Equal(t, "Updated Name 2", b.Name)
			assert.Equal(t, updatedBy, *b.UpdatedBy)
		}
	}
}
