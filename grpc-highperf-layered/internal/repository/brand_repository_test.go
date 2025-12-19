package repository_test

import (
	"codegen/internal/domain"
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
	newBrand := &domain.Brand{
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
	require.NoError(t, err) // Should be no error if we return nil, nil
	assert.Nil(t, fetchedDeleted)

	// 6. DeleteByIds
	brandToDelete1 := &domain.Brand{
		Name:      "DeleteByIds 1",
		Slug:      "del-1-" + uuid.New().String(),
		Logo:      "logo1.png",
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
	}
	brandToDelete2 := &domain.Brand{
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

func TestBrandRepository_GetById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)

	fetched, err := repo.GetById(context.Background(), 2147483647) // Max int32
	// Expecting nil, nil based on user preference
	require.NoError(t, err)
	assert.Nil(t, fetched, "Should return nil for non-existent record")
}

func TestBrandRepository_Upsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// 1. Upsert (Insert Scenario - Auto Id)
	t.Run("InsertNewBrand_AutoId", func(t *testing.T) {
		brand := &domain.Brand{
			Name:      "New Brand AutoId " + uuid.New().String(),
			Slug:      "new-brand-auto-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
		// Should handle Id=0 by Inserting
		err := repo.Upsert(ctx, brand)
		require.NoError(t, err)
		assert.NotZero(t, brand.Id, "Id should be generated")
	})

	// 2. Upsert (Insert Scenario - Explicit Id)
	t.Run("InsertNewBrand_ExplicitId", func(t *testing.T) {
		upsertId := int32(99999)
		// Clean up potential leftover
		_, _ = repo.Delete(ctx, upsertId)

		upsertBrand := &domain.Brand{
			Id:        upsertId,
			Name:      "Upsert Brand Explicit",
			Slug:      "upsert-explicit-" + uuid.New().String(),
			Logo:      "upsert-logo.png",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		err := repo.Upsert(ctx, upsertBrand)
		require.NoError(t, err)

		fetchedUpsert, err := repo.GetById(ctx, upsertId)
		require.NoError(t, err)
		assert.NotNil(t, fetchedUpsert)
		assert.Equal(t, upsertBrand.Name, fetchedUpsert.Name)
	})

	// 3. Upsert (Update Scenario)
	t.Run("UpdateExistingBrand", func(t *testing.T) {
		// First Insert
		brand := &domain.Brand{
			Name:      "Old Name " + uuid.New().String(),
			Slug:      "old-name-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
		err := repo.Insert(ctx, brand)
		require.NoError(t, err)

		// Modify
		brand.Name = "New Name " + uuid.New().String()
		updatedBy := uuid.New()
		timeRef := time.Now()
		brand.UpdatedBy = &updatedBy
		brand.UpdatedAt = &timeRef

		err = repo.Upsert(ctx, brand)
		require.NoError(t, err)

		// Verify
		fetched, err := repo.GetById(ctx, brand.Id)
		require.NoError(t, err)
		assert.Equal(t, brand.Name, fetched.Name)
		assert.Equal(t, updatedBy, *fetched.UpdatedBy)
	})
}

func TestBrandRepository_BulkInsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	count := 10
	list := make([]*domain.Brand, count)
	for i := 0; i < count; i++ {
		list[i] = &domain.Brand{
			Name:      "CopyFrom Brand " + uuid.New().String(),
			Slug:      "copyfrom-brand-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
	}

	insertedCount, err := repo.BulkInsertOneShot(ctx, list)
	require.NoError(t, err)
	assert.Equal(t, int64(count), insertedCount)
}

func TestBrandRepository_BulkUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Setup
	brand1 := &domain.Brand{
		Name:      "Pre-Update-Copy 1",
		Slug:      "pre-upd-copy-1-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	brand2 := &domain.Brand{
		Name:      "Pre-Update-Copy 2",
		Slug:      "pre-upd-copy-2-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, repo.Insert(ctx, brand1))
	require.NoError(t, repo.Insert(ctx, brand2))

	// Modify
	newUUID := uuid.New()
	newTime := time.Now()

	brand1.Name = "Post-Update-Copy 1"
	brand1.UpdatedBy = &newUUID
	brand1.UpdatedAt = &newTime

	brand2.Name = "Post-Update-Copy 2"
	brand2.UpdatedBy = &newUUID
	brand2.UpdatedAt = &newTime

	// Execute BulkUpdate
	count, err := repo.BulkUpdate(ctx, []*domain.Brand{brand1, brand2}, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)

	// Verify
	fetched, err := repo.GetByIds(ctx, []int32{brand1.Id, brand2.Id})
	require.NoError(t, err)
	assert.Len(t, fetched, 2)

	for _, b := range fetched {
		if b.Id == brand1.Id {
			assert.Equal(t, "Post-Update-Copy 1", b.Name)
			assert.Equal(t, newUUID, *b.UpdatedBy)
		} else if b.Id == brand2.Id {
			assert.Equal(t, "Post-Update-Copy 2", b.Name)
		}
	}
}

func TestBrandRepository_BulkInsert_Large(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()
	batchSize := 1000

	// 2005 items to test chunking (batchSize=1000, so 3 chunks: 1000, 1000, 5)
	count := 2005
	list := make([]*domain.Brand, count)
	for i := 0; i < count; i++ {
		list[i] = &domain.Brand{
			Name:      "Large Brand " + uuid.New().String(),
			Slug:      "large-brand-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
	}

	insertedCount, err := repo.BulkInsert(ctx, list, batchSize)
	require.NoError(t, err)
	assert.Equal(t, int64(count), insertedCount)

	// Verify count logic
	// We might have other records, so we'll just check if it's at least count
	dbCount, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, dbCount, int64(count))
}
