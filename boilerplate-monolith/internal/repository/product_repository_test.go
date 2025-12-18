package repository_test

import (
	"codegen/internal/entity"
	"codegen/internal/repository"
	"context"
	"fmt"
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

	// 6. DeleteByIds
	p1 := &entity.Product{
		Id:            uuid.New(),
		BrandId:       int(brand.Id),
		Name:          "DeleteByIds 1",
		Sku:           strPtr("Del-1-" + uuid.New().String()),
		StockQuantity: 10,
		Price:         10.0,
		CreatedBy:     uuid.New(),
		CreatedAt:     time.Now(),
	}
	p2 := &entity.Product{
		Id:            uuid.New(),
		BrandId:       int(brand.Id),
		Name:          "DeleteByIds 2",
		Sku:           strPtr("Del-2-" + uuid.New().String()),
		StockQuantity: 10,
		Price:         10.0,
		CreatedBy:     uuid.New(),
		CreatedAt:     time.Now(),
	}

	_, err = repo.Insert(ctx, p1)
	require.NoError(t, err)

	_, err = repo.Insert(ctx, p2)
	require.NoError(t, err)

	idsToDelete := []uuid.UUID{p1.Id, p2.Id}
	deletedCount, err := repo.DeleteByIds(ctx, idsToDelete)
	require.NoError(t, err)
	assert.Equal(t, int64(2), deletedCount)

	// Verify deletion
	f1, err := repo.GetById(ctx, p1.Id)
	require.NoError(t, err)
	assert.Nil(t, f1)

	f2, err := repo.GetById(ctx, p2.Id)
	require.NoError(t, err)
	assert.Nil(t, f2)
}

func TestProductRepository_GetById_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)

	fetched, err := repo.GetById(context.Background(), uuid.New()) // Max int32
	// Expecting nil, nil based on user preference
	require.NoError(t, err)
	assert.Nil(t, fetched, "Should return nil for non-existent record")
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

func TestProductRepository_BulkUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Cleanup
	_, err := db.Pool().Exec(ctx, "DELETE FROM catalog.products")
	require.NoError(t, err)

	// Insert Brand
	brand := &entity.Brand{
		Name:      "Bulk Update Brand",
		Slug:      "bulk-update-brand-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	// 1. Prepare products
	count := 5
	list := make([]*entity.Product, count)
	for i := 0; i < count; i++ {
		list[i] = &entity.Product{
			Id:            uuid.New(),
			BrandId:       int(brand.Id),
			Name:          fmt.Sprintf("Original Product %d", i),
			Sku:           strPtr(fmt.Sprintf("UPDATE-SKU-%d-%s", i, uuid.New().String())),
			StockQuantity: 10,
			Price:         50.0,
			CreatedBy:     uuid.New(),
			CreatedAt:     time.Now(),
		}
	}

	_, err = repo.BulkInsert(ctx, list)
	require.NoError(t, err)

	// 2. Modify products
	updatedBy := uuid.New()
	updatedAt := time.Now().UTC().Truncate(time.Millisecond) // Truncate to align with DB precision if needed, though usually fine.

	for i, p := range list {
		p.Name = fmt.Sprintf("Updated Product %d", i)
		p.Price = 100.0 + float64(i)
		p.StockQuantity = 20 + i
		p.UpdatedBy = &updatedBy
		p.UpdatedAt = &updatedAt
	}

	// 3. BulkUpdate
	affected, err := repo.BulkUpdate(ctx, list)
	require.NoError(t, err)
	assert.Equal(t, int64(count), affected)

	// 4. Verify updates
	for _, p := range list {
		fetched, err := repo.GetById(ctx, p.Id)
		require.NoError(t, err)
		assert.Equal(t, p.Name, fetched.Name)
		assert.Equal(t, p.Price, fetched.Price)
		assert.Equal(t, p.StockQuantity, fetched.StockQuantity)
		assert.Equal(t, updatedBy, *fetched.UpdatedBy)
	}
}

func TestProductRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Setup data: 1 Brand, 3 Products
	brand := &entity.Brand{
		Name:      "Search Brand " + uuid.New().String(),
		Slug:      "search-brand-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	p1 := &entity.Product{
		BrandId: int(brand.Id),
		Name:    "Alpha Product",
		Sku:     strPtr("S-1"),
		Price:   10,
	}
	p2 := &entity.Product{
		BrandId: int(brand.Id),
		Name:    "Alpha Beta Product",
		Sku:     strPtr("S-2"),
		Price:   20,
	}
	p3 := &entity.Product{
		BrandId: int(brand.Id),
		Name:    "Gamma Product",
		Sku:     strPtr("S-3"),
		Price:   30,
	}
	_, err := repo.BulkInsert(ctx, []*entity.Product{p1, p2, p3})
	require.NoError(t, err)

	// Test 1: Search by Name (Partial)
	res, err := repo.Search(ctx, &entity.ProductSearchFilter{
		Name:  "Alpha",
		Limit: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.Total)
	assert.Len(t, res.Items, 2)

	// Test 2: Search with Pagination
	res, err = repo.Search(ctx, &entity.ProductSearchFilter{
		Name:   "Alpha",
		Limit:  1,
		Offset: 0,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.Total) // Total should still be 2
	assert.Len(t, res.Items, 1)          // But we only took 1
}

func TestProductRepository_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewProductRepository(db)
	brandRepo := repository.NewBrandRepository(db)
	ctx := context.Background()

	// Need a brand
	brand := &entity.Brand{
		Name:      "Delete Brand " + uuid.New().String(),
		Slug:      "del-brand-" + uuid.New().String(),
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	require.NoError(t, brandRepo.Insert(ctx, brand))

	// Insert
	p := &entity.Product{
		Id:        uuid.New(),
		BrandId:   int(brand.Id),
		Name:      "Soft Delete Me",
		CreatedAt: time.Now(),
		CreatedBy: uuid.New(),
	}
	_, err := repo.Insert(ctx, p)
	require.NoError(t, err)

	// Soft Delete
	deleter := uuid.New()
	affected, err := repo.SoftDelete(ctx, p.Id, deleter)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)

	// Verify not found via GetById
	fetched, err := repo.GetById(ctx, p.Id)
	require.NoError(t, err)
	assert.Nil(t, fetched, "Soft deleted product should not be returned by GetById")
}
