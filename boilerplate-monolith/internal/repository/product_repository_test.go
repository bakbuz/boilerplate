package repository_test

import (
	"codegen/internal/entity"
	"codegen/internal/repository"
	"codegen/internal/repository/dto"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string {
	return &s
}

func TestProductRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	ctx := context.Background()

	// Cleanup
	_, err := db.Pool().Exec(ctx, "DELETE FROM catalog.products")
	require.NoError(t, err)

	brandRepo := repository.NewBrandRepository(db)
	brand := &entity.Brand{
		Name:      "Product Test Brand " + uuid.New().String(),
		Slug:      "prod-test-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	newProduct := &entity.Product{
		Id:            uuid.New(),
		BrandId:       int(brand.Id),
		Name:          "Test Product",
		Sku:           strPtr("SKU-" + uuid.New().String()),
		Summary:       strPtr("Summary"),
		Storyline:     strPtr("Storyline"),
		StockQuantity: 100,
		Price:         99.99,
		CreatedBy:     uuid.New(),
		CreatedAt:     time.Now(),
	}

	// 1. Insert
	affected, err := repo.Insert(ctx, newProduct)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	// 2. GetById
	fetched, err := repo.GetById(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, newProduct.Name, fetched.Name)
	assert.Equal(t, newProduct.BrandId, fetched.BrandId)

	// 3. Update
	newProduct.Name = "Updated Product"
	updatedBy := uuid.New()
	updatedAt := time.Now()
	newProduct.UpdatedBy = &updatedBy
	newProduct.UpdatedAt = &updatedAt

	affected, err = repo.Update(ctx, newProduct)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	fetchedUpdated, err := repo.GetById(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.Equal(t, "Updated Product", fetchedUpdated.Name)

	// 4. Search
	t.Run("Search", func(t *testing.T) {
		res, err := repo.Search(ctx, &dto.ProductSearchFilter{
			Name: "Updated",
			Take: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, res.Total)
		assert.NotEmpty(t, res.Items)
		assert.Equal(t, newProduct.Id, res.Items[0].Id)
	})

	// 5. Delete
	affected, err = repo.Delete(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	fetchedDeleted, err := repo.GetById(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.Nil(t, fetchedDeleted)
}

func TestProductRepository_BulkInsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Insert Brand
	brand := &entity.Brand{
		Name:      "Bulk Product Brand",
		Slug:      "bulk-prod-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	count := 5
	list := make([]*entity.Product, count)
	for i := 0; i < count; i++ {
		list[i] = &entity.Product{
			Id:        uuid.New(),
			BrandId:   int(brand.Id),
			Name:      "Bulk Product",
			Sku:       strPtr("BULK-" + uuid.New().String()),
			Price:     10.0,
			CreatedAt: time.Now(),
			CreatedBy: uuid.New(),
		}
	}

	err := repo.BulkInsert(ctx, list)
	require.NoError(t, err)

	// Verify count
	c, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, c, int64(count))
}
