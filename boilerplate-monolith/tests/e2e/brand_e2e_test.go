package e2e_test

import (
	catalogv1 "codegen/api/gen/catalog/v1"
	"codegen/internal/entity"
	"codegen/internal/service"
	"codegen/internal/transport/handler"
	"codegen/internal/transport/interceptor"
	"context"
	"net"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// 1. Sanal Ağ Dinleyicisi (Bufconn Listener)
const brandBufSize = 1024 * 1024

var brandLis *bufconn.Listener

func init() {
	brandLis = bufconn.Listen(brandBufSize)
}

// Test interceptor to set user ID (duplicated for standalone run)
func brandTestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Set a dummy user ID if not set
	if ctx.Value(interceptor.UserIdKey) == nil {
		ctx = context.WithValue(ctx, interceptor.UserIdKey, uuid.New().String())
	}
	return handler(ctx, req)
}

// 2. Test Ortamını Ayağa Kaldıran Yardımcı Fonksiyon
func brandBufDialer(context.Context, string) (net.Conn, error) {
	return brandLis.Dial()
}

// InMemoryBrandRepo implements repository.BrandRepository
type InMemoryBrandRepo struct {
	items map[int32]*entity.Brand
	mu    sync.RWMutex
	// Auto-increment counter
	nextId int32
}

func NewInMemoryBrandRepo() *InMemoryBrandRepo {
	return &InMemoryBrandRepo{
		items:  make(map[int32]*entity.Brand),
		nextId: 1,
	}
}

func (r *InMemoryBrandRepo) GetAll(ctx context.Context) ([]*entity.Brand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Brand
	for _, b := range r.items {
		result = append(result, b)
	}
	return result, nil
}

func (r *InMemoryBrandRepo) GetByIds(ctx context.Context, ids []int32) ([]*entity.Brand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Brand
	for _, id := range ids {
		if b, exists := r.items[id]; exists {
			result = append(result, b)
		}
	}
	return result, nil
}

func (r *InMemoryBrandRepo) GetById(ctx context.Context, id int32) (*entity.Brand, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if b, exists := r.items[id]; exists {
		return b, nil
	}
	return nil, nil // Return nil, nil for not found as per repository pattern seen in product
}

func (r *InMemoryBrandRepo) Insert(ctx context.Context, e *entity.Brand) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.Id == 0 {
		e.Id = r.nextId
		r.nextId++
	}
	r.items[e.Id] = e
	return nil
}

func (r *InMemoryBrandRepo) Update(ctx context.Context, e *entity.Brand) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[e.Id]; exists {
		r.items[e.Id] = e
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryBrandRepo) Delete(ctx context.Context, id int32) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[id]; exists {
		delete(r.items, id)
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryBrandRepo) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := int64(0)
	for _, id := range ids {
		if _, exists := r.items[id]; exists {
			delete(r.items, id)
			count++
		}
	}
	return count, nil
}

func (r *InMemoryBrandRepo) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.items)), nil
}

func (r *InMemoryBrandRepo) Upsert(ctx context.Context, e *entity.Brand) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.Id == 0 {
		e.Id = r.nextId
		r.nextId++
	}
	r.items[e.Id] = e
	return nil
}

func (r *InMemoryBrandRepo) BulkInsert(ctx context.Context, list []*entity.Brand) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range list {
		if e.Id == 0 {
			e.Id = r.nextId
			r.nextId++
		}
		r.items[e.Id] = e
	}
	return int64(len(list)), nil
}

func (r *InMemoryBrandRepo) BulkUpdate(ctx context.Context, list []*entity.Brand) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := int64(0)
	for _, e := range list {
		if _, exists := r.items[e.Id]; exists {
			r.items[e.Id] = e
			count++
		}
	}
	return count, nil
}

func (r *InMemoryBrandRepo) BulkInsertTran(ctx context.Context, list []*entity.Brand) error {
	_, err := r.BulkInsert(ctx, list)
	return err
}

func (r *InMemoryBrandRepo) BulkUpdateTran(ctx context.Context, list []*entity.Brand) error {
	_, err := r.BulkUpdate(ctx, list)
	return err
}

func setupBrandTestServer(t *testing.T) (catalogv1.BrandServiceClient, func()) {
	// gRPC Server with Interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(brandTestInterceptor),
	)

	// Mock Repo
	mockRepo := NewInMemoryBrandRepo()

	// Service
	// Note: NewBrandService expects repository.BrandRepository
	mockSvc := service.NewBrandService(mockRepo)

	// Handler
	mockHandler := handler.NewBrandHandler(mockSvc)

	catalogv1.RegisterBrandServiceServer(s, mockHandler)

	// Launch server on dedicated listener
	go func() {
		if err := s.Serve(brandLis); err != nil {
			// We might expect error on Close
		}
	}()

	// Client
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(brandBufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := catalogv1.NewBrandServiceClient(conn)

	return client, func() {
		err := conn.Close()
		require.NoError(t, err)
		s.Stop()
	}
}

func TestBrandService_E2E(t *testing.T) {
	client, teardown := setupBrandTestServer(t)
	defer teardown()

	userId := uuid.New().String()
	ctx := context.WithValue(context.Background(), interceptor.UserIdKey, userId)
	var createdId int32

	// 1. Create Brand
	t.Run("Create Brand", func(t *testing.T) {
		req := &catalogv1.CreateBrandRequest{
			Name: "Razer",
			Logo: "https://example.com/razer.png",
		}

		resp, err := client.Create(ctx, req)
		require.NoError(t, err)
		require.NotZero(t, resp.Id)
		require.Equal(t, req.Name, resp.Name)
		// Request does not have URL, but Response does.
		// Since we don't send slug, check if service generates one or it remains empty if not handled.
		// For now let's just assert ID and Name.
		require.Equal(t, req.Logo, resp.Logo)
		createdId = resp.Id
	})

	// 2. Get Brand
	t.Run("Get Brand", func(t *testing.T) {
		req := &catalogv1.GetBrandRequest{Id: createdId}
		resp, err := client.Get(ctx, req)

		require.NoError(t, err)
		require.Equal(t, createdId, resp.Id)
		require.Equal(t, "Razer", resp.Name)
	})

	// 3. Update Brand
	t.Run("Update Brand", func(t *testing.T) {
		req := &catalogv1.UpdateBrandRequest{
			Id:   createdId,
			Name: "Razer Inc.",
			Logo: "https://example.com/razer-new.png",
		}

		resp, err := client.Update(ctx, req)
		require.NoError(t, err)
		require.Equal(t, "Razer Inc.", resp.Name)
		// Updated fields
	})

	// 4. List Brands
	t.Run("List Brands", func(t *testing.T) {
		req := &catalogv1.ListBrandsRequest{}
		resp, err := client.List(ctx, req)

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Items), 1)

		found := false
		for _, item := range resp.Items {
			if item.Id == createdId {
				found = true
				require.Equal(t, "Razer Inc.", item.Name)
				break
			}
		}
		require.True(t, found, "Created brand should be in the list")
	})

	// 5. Delete Brand
	t.Run("Delete Brand", func(t *testing.T) {
		req := &catalogv1.DeleteBrandRequest{Id: createdId}
		_, err := client.Delete(ctx, req)
		require.NoError(t, err)
	})

	// 6. Verify Deletion
	t.Run("Verify Deletion", func(t *testing.T) {
		req := &catalogv1.GetBrandRequest{Id: createdId}
		_, err := client.Get(ctx, req)

		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	// Error Case: Get Non-Existing
	t.Run("Get Non-Existing Brand", func(t *testing.T) {
		// arbitrary large ID
		req := &catalogv1.GetBrandRequest{Id: 99999}
		_, err := client.Get(ctx, req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	// Error Case: Delete Non-Existing
	t.Run("Delete Non-Existing Brand", func(t *testing.T) {
		req := &catalogv1.DeleteBrandRequest{Id: 99999}
		_, err := client.Delete(ctx, req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})
}
