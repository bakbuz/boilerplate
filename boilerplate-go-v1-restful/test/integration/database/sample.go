package database_test

// import (
// 	"context"
// 	"fmt"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"codegen/internal/database"
// 	"codegen/internal/entity"
// 	"codegen/internal/repository"
// 	"codegen/utils"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func Test_Product_50(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	for i := 1; i <= 50; i++ {
// 		var name = "name " + strconv.Itoa(i)
// 		product := &entity.Product{Id: newGuid(), Name: name, CreatedBy: 0, CreatedAt: time.Now()}

// 		err = r.Insert(context.Background(), product)
// 		require.NoError(t, err)
// 	}
// }

// func Test_Product_BulkInsert(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	product1 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now()}
// 	product2 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now()}
// 	products := []*entity.Product{product1, product2}

// 	err = r.BulkInsert(context.Background(), products)
// 	require.NoError(t, err)

// }

// func Test_Product_GetByIds(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	// GetByIds
// 	var product1Id uuid.UUID = uuid.MustParse("018E8049-5A01-70F3-ACAE-4B6D5B5D91F5")
// 	var product2Id uuid.UUID = uuid.MustParse("018E8049-5A25-7D3B-AD2B-677CC1AAF661")
// 	data, err := r.GetByIds(context.Background(), []uuid.UUID{product1Id, product2Id})

// 	require.NoError(t, err)
// 	assert.Len(t, data, 2)
// }

// func Test_Product_CRUD(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	// Create
// 	product1 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now().UTC()} //UTC olması gerekiyor
// 	err = r.Insert(context.Background(), product1)
// 	require.NoError(t, err)

// 	product2 := &entity.Product{Id: newGuid(), Name: random.Str(4), CreatedBy: 0, CreatedAt: time.Now().UTC()}
// 	err = r.Insert(context.Background(), product2)
// 	require.NoError(t, err)

// 	// Count
// 	count, err := r.Count(context.Background())
// 	require.NoError(t, err)
// 	assert.GreaterOrEqual(t, count, int64(2))

// 	// GetByIds
// 	data, err := r.GetByIds(context.Background(), []uuid.UUID{product1.Id, product2.Id})
// 	require.NoError(t, err)
// 	assert.Len(t, data, 2)

// 	pid1 := product1.Id.String()
// 	pid2 := product2.Id.String()

// 	aid1 := data[0].Id.String()
// 	aid2 := data[1].Id.String()

// 	fmt.Println(pid1, pid2)
// 	fmt.Println(aid1, aid2)

// 	actual1 := data[0]
// 	actual2 := data[1]

// 	if actual1.Id != product1.Id {
// 		actual1 = data[1]
// 		actual2 = data[0]
// 	}

// 	//assert.Equal(t, actual1.Id, product1.Id)
// 	assert.Equal(t, actual1.Name, product1.Name)
// 	assert.Equal(t, actual1.CreatedBy, product1.CreatedBy)
// 	assert.Equal(t, actual1.CreatedAt.Truncate(time.Second), product1.CreatedAt.Truncate(time.Second))

// 	//assert.Equal(t, actual2.Id, product2.Id)
// 	assert.Equal(t, actual2.Name, product2.Name)
// 	assert.Equal(t, actual2.CreatedBy, product2.CreatedBy)
// 	assert.Equal(t, actual2.CreatedAt.Truncate(time.Second), product2.CreatedAt.Truncate(time.Second))

// 	// GetById
// 	single, err := r.GetById(context.Background(), product1.Id)
// 	require.NoError(t, err)
// 	assert.NotNil(t, single)

// 	//assert.Equal(t, single.Id, product1.Id)
// 	assert.Equal(t, single.Name, product1.Name)
// 	assert.Equal(t, single.CreatedBy, product1.CreatedBy)
// 	assert.Equal(t, single.CreatedAt.Truncate(time.Second), product1.CreatedAt.Truncate(time.Second))

// 	// Update
// 	single.UpdatedBy = pointer(1)
// 	single.UpdatedAt = pointer(time.Now().UTC())

// 	rowsAffected, err := r.Update(context.Background(), single)
// 	require.NoError(t, err)
// 	assert.Equal(t, int64(1), rowsAffected)

// 	// Updated
// 	updated, err := r.GetById(context.Background(), product1.Id)
// 	require.NoError(t, err)
// 	assert.NotNil(t, updated)

// 	assert.Equal(t, single.UpdatedBy, updated.UpdatedBy)
// 	assert.Equal(t, single.UpdatedAt.Truncate(time.Second), updated.UpdatedAt.Truncate(time.Second))

// 	// SoftDelete
// 	rowsAffected2, err := r.SoftDelete(context.Background(), product1.Id, 99)
// 	require.NoError(t, err)
// 	assert.NotNil(t, updated)

// }

// func Test_Product_Create(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	e := &entity.Product{
// 		Name:      random.Str(4),
// 		CreatedBy: 0,
// 		CreatedAt: time.Now(),
// 	}
// 	err = r.Insert(context.Background(), e)
// 	require.NoError(t, err)

// 	assert.NotEqual(t, e.Id, uuid.UUID{})

// 	result, err := db.GetString(fmt.Sprintf("SELECT Name FROM Products WHERE Id = '%s';", e.Id.String()))
// 	require.NoError(t, err)

// 	assert.Equal(t, e.Name, result)
// }

// /*
// func Test_Product_NoRecords(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)
// 	result, err := r.GetProduct(context.Background(), -1)
// 	require.NoError(t, err)

// 	assert.Nil(t, result)
// }

// func Test_Product_GetOne(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	name := random.Str(16)
// 	listId, err := r.CreateProduct(context.Background(), name)
// 	require.NoError(t, err)

// 	list, err := r.GetProduct(context.Background(), listId)
// 	require.NoError(t, err)

// 	require.NotNil(t, list.CreatedAt)
// 	assert.Equal(t, listId, list.Id)
// 	assert.Equal(t, name, list.Name)
// }

// func Test_Task_CreateTaskWithProduct(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)
// 	require.NoError(t, err)

// 	listName := random.Str(16)
// 	listId, err := r.CreateProduct(context.Background(), listName)
// 	require.NoError(t, err)

// 	taskName := random.Str(16)
// 	taskId, err := r.CreateTask(context.Background(), listId, taskName)
// 	require.NoError(t, err)

// 	list, err := r.GetProduct(context.Background(), listId)
// 	require.NoError(t, err)

// 	require.NotNil(t, list.CreatedAt)
// 	assert.Equal(t, listId, list.Id)
// 	assert.Equal(t, listName, list.Name)

// 	tasks, err := r.GetTasks(context.Background(), listId)
// 	require.NoError(t, err)

// 	assert.Greater(t, len(tasks), 0)

// 	task := tasks[0]
// 	assert.Equal(t, taskId, task.Id)
// 	assert.Equal(t, taskName, task.Name)
// }

// func Test_Task_CreateTaskWithoutProduct(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)

// 	name := random.Str(8)
// 	id, err := r.CreateTask(context.Background(), -1, name)
// 	require.Error(t, err)
// 	assert.Equal(t, id, -1)
// }

// func Test_Task_DeleteTask(t *testing.T) {
// 	c := RequireConfig(t)
// 	db, err := database.New(c.DataSources.Default)
// 	require.NoError(t, err)
// 	defer db.Close()

// 	r := repository.NewProductRepository(db)
// 	require.NoError(t, err)

// 	listName := random.Str(16)
// 	listId, err := r.CreateProduct(context.Background(), listName)
// 	require.NoError(t, err)

// 	taskName := random.Str(16)
// 	taskId, err := r.CreateTask(context.Background(), listId, taskName)
// 	require.NoError(t, err)

// 	err = r.DeleteTask(context.Background(), taskId)
// 	require.NoError(t, err)

// 	count, err := db.Count(context.Background(), "SELECT COUNT(*) FROM todo_tasks WHERE id = "+strconv.Itoa(taskId))
// 	require.NoError(t, err)

// 	assert.Equal(t, int64(0), count)
// }
// */
