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
	require.NoError(t, err) // Should be no error if we return nil, nil
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

	// 1. Upsert (Insert Scenario - Auto ID)
	t.Run("InsertNewBrand_AutoID", func(t *testing.T) {
		brand := &entity.Brand{
			Name:      "New Brand AutoID " + uuid.New().String(),
			Slug:      "new-brand-auto-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
		// Should handle ID=0 by Inserting
		err := repo.Upsert(ctx, brand)
		require.NoError(t, err)
		assert.NotZero(t, brand.Id, "ID should be generated")
	})

	// 2. Upsert (Insert Scenario - Explicit ID)
	t.Run("InsertNewBrand_ExplicitID", func(t *testing.T) {
		upsertID := int32(99999)
		// Clean up potential leftover
		_, _ = repo.Delete(ctx, upsertID)

		upsertBrand := &entity.Brand{
			Id:        upsertID,
			Name:      "Upsert Brand Explicit",
			Slug:      "upsert-explicit-" + uuid.New().String(),
			Logo:      "upsert-logo.png",
			CreatedBy: uuid.New(),
			CreatedAt: time.Now(),
		}

		err := repo.Upsert(ctx, upsertBrand)
		require.NoError(t, err)

		fetchedUpsert, err := repo.GetById(ctx, upsertID)
		require.NoError(t, err)
		assert.NotNil(t, fetchedUpsert)
		assert.Equal(t, upsertBrand.Name, fetchedUpsert.Name)
	})

	// 3. Upsert (Update Scenario)
	t.Run("UpdateExistingBrand", func(t *testing.T) {
		// First Insert
		brand := &entity.Brand{
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

func TestBrandRepository_GetById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)

	fetched, err := repo.GetById(context.Background(), 2147483647) // Max int32
	// Expecting nil, nil based on user preference
	require.NoError(t, err)
	assert.Nil(t, fetched, "Should return nil for non-existent record")
}

func TestBrandRepository_BulkInsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	count := 10
	list := make([]*entity.Brand, count)
	for i := 0; i < count; i++ {
		list[i] = &entity.Brand{
			Name:      "CopyFrom Brand " + uuid.New().String(),
			Slug:      "copyfrom-brand-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
	}

	insertedCount, err := repo.BulkInsert(ctx, list)
	require.NoError(t, err)
	assert.Equal(t, int64(count), insertedCount)
}

func TestBrandRepository_BulkUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Setup
	brand1 := &entity.Brand{
		Name:      "Pre-Update-Copy 1",
		Slug:      "pre-upd-copy-1-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	brand2 := &entity.Brand{
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
	count, err := repo.BulkUpdate(ctx, []*entity.Brand{brand1, brand2})
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

func TestBrandRepository_BulkInsertTran(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	count := 10
	list := make([]*entity.Brand, count)
	for i := 0; i < count; i++ {
		list[i] = &entity.Brand{
			Name:      "Bulk Brand " + uuid.New().String(),
			Slug:      "bulk-brand-" + uuid.New().String(),
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
	}

	err := repo.BulkInsertTran(ctx, list)
	require.NoError(t, err)

	// Verify count roughly (or exact if we clear DB first, but we didn't clear here)
	// Just ensure no error
}

func TestBrandRepository_BulkUpdateTran(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// 1. Prepare data (Insert two brands)
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
	err = repo.BulkUpdateTran(ctx, []*entity.Brand{brand1, brand2})
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
