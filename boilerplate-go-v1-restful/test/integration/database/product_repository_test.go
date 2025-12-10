package database_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"codegen/internal/database"
	"codegen/internal/entity"
	"codegen/internal/repository"
	"codegen/utils/random"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Product_BulkInsert(t *testing.T) {
	c := RequireConfig(t)
	db, err := database.New(c.DataSources.Default)
	require.NoError(t, err)
	defer db.Close()

	r := repository.NewProductRepository(db)

	product1 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now()}
	product2 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now()}
	products := []*entity.Product{product1, product2}

	err = r.BulkInsert(context.Background(), products)
	require.NoError(t, err)

}

func Test_Product_CRUD(t *testing.T) {
	c := RequireConfig(t)
	db, err := database.New(c.DataSources.Default)
	require.NoError(t, err)
	defer db.Close()

	r := repository.NewProductRepository(db)

	// Create
	product1 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now().UTC()} //UTC olması gerekiyor
	rowsAffected1, err := r.Insert(context.Background(), product1)
	require.NoError(t, err)
	assert.Equal(t, 1, rowsAffected1)

	product2 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now().UTC()}
	rowsAffected2, err := r.Insert(context.Background(), product2)
	require.NoError(t, err)
	assert.Equal(t, 1, rowsAffected2)

	// Count
	count, err := r.Count(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(2))

	// GetByIds
	data, err := r.GetByIds(context.Background(), []uuid.UUID{product1.Id, product2.Id})
	require.NoError(t, err)
	assert.Len(t, data, 2)

	pid1 := product1.Id.String()
	pid2 := product2.Id.String()

	aid1 := data[0].Id.String()
	aid2 := data[1].Id.String()

	fmt.Println(pid1, pid2)
	fmt.Println(aid1, aid2)

	actual1 := data[0]
	actual2 := data[1]

	if actual1.Id != product1.Id {
		actual1 = data[1]
		actual2 = data[0]
	}

	//assert.Equal(t, actual1.Id, product1.Id)
	assert.Equal(t, actual1.Name, product1.Name)
	assert.Equal(t, actual1.CreatedBy, product1.CreatedBy)
	assert.Equal(t, actual1.CreatedAt.Truncate(time.Second), product1.CreatedAt.Truncate(time.Second))

	//assert.Equal(t, actual2.Id, product2.Id)
	assert.Equal(t, actual2.Name, product2.Name)
	assert.Equal(t, actual2.CreatedBy, product2.CreatedBy)
	assert.Equal(t, actual2.CreatedAt.Truncate(time.Second), product2.CreatedAt.Truncate(time.Second))

	// GetById
	single, err := r.GetById(context.Background(), product1.Id)
	require.NoError(t, err)
	assert.NotNil(t, single)

	//assert.Equal(t, single.Id, product1.Id)
	assert.Equal(t, single.Name, product1.Name)
	assert.Equal(t, single.CreatedBy, product1.CreatedBy)
	assert.Equal(t, single.CreatedAt.Truncate(time.Second), product1.CreatedAt.Truncate(time.Second))

	// Update
	single.UpdatedBy = pointer(1)
	single.UpdatedAt = pointer(time.Now().UTC())

	rowsAffected, err := r.Update(context.Background(), single)
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	// Updated
	updated, err := r.GetById(context.Background(), product1.Id)
	require.NoError(t, err)
	assert.NotNil(t, updated)

	assert.Equal(t, single.UpdatedBy, updated.UpdatedBy)
	assert.Equal(t, single.UpdatedAt.Truncate(time.Second), updated.UpdatedAt.Truncate(time.Second))

	// SoftDelete
	softDeletedAffected, err := r.SoftDelete(context.Background(), product1.Id, 99)
	require.NoError(t, err)
	assert.Equal(t, 1, softDeletedAffected)

}
