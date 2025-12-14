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

	// 4. GetAll
	list, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, list)
	found := false
	for _, b := range list {
		if b.Id == newProduct.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "Newly created product should be in GetAll list")

	// 5. Delete
	affected, err = repo.Delete(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	fetchedDeleted, err := repo.GetById(ctx, newProduct.Id)
	require.NoError(t, err)
	assert.Nil(t, fetchedDeleted)
}

func TestProductRepository_Upsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Setup Brand
	brand := &entity.Brand{
		Name:      "Upsert Brand " + uuid.New().String(),
		Slug:      "upsert-brand-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	// 1. Insert new product via Upsert
	product := &entity.Product{
		BrandId:       int(brand.Id),
		Name:          "Upsert Product",
		Sku:           strPtr("UPSERT-1"),
		StockQuantity: 50,
		Price:         19.99,
		CreatedAt:     time.Now(),
		CreatedBy:     uuid.New(),
	}

	err := repo.Upsert(ctx, product)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, product.Id)

	fetched, err := repo.GetById(ctx, product.Id)
	require.NoError(t, err)
	assert.Equal(t, product.Name, fetched.Name)

	// 2. Update existing product via Upsert
	product.Name = "Upsert Product Updated"
	product.Price = 29.99
	updatedBy := uuid.New()
	updatedAt := time.Now()
	product.UpdatedBy = &updatedBy
	product.UpdatedAt = &updatedAt

	err = repo.Upsert(ctx, product)
	require.NoError(t, err)

	fetchedUpdated, err := repo.GetById(ctx, product.Id)
	require.NoError(t, err)
	assert.Equal(t, "Upsert Product Updated", fetchedUpdated.Name)
	assert.Equal(t, 29.99, fetchedUpdated.Price)
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

	insertedCount, err := repo.BulkInsert(ctx, list)
	require.NoError(t, err)
	assert.Equal(t, int64(count), insertedCount)

	// Verify count
	c, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, c, int64(count))
}
