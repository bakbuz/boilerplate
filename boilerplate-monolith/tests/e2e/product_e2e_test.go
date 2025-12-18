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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// 1. Sanal Ağ Dinleyicisi (Bufconn Listener)
const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

// Test interceptor to set user ID
func testInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Set a dummy user ID if not set
	if ctx.Value(interceptor.UserIdKey) == nil {
		ctx = context.WithValue(ctx, interceptor.UserIdKey, uuid.New().String())
	}
	return handler(ctx, req)
}

// 2. Test Ortamını Ayağa Kaldıran Yardımcı Fonksiyon
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type InMemoryProductRepo struct {
	items map[string]*entity.Product
	mu    sync.RWMutex
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		items: make(map[string]*entity.Product),
	}
}

func (r *InMemoryProductRepo) GetAll(ctx context.Context) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Product
	for _, p := range r.items {
		if !p.Deleted {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemoryProductRepo) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Product
	for _, id := range ids {
		if p, exists := r.items[id.String()]; exists && !p.Deleted {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemoryProductRepo) GetById(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if p, exists := r.items[id.String()]; exists && !p.Deleted {
		return p, nil
	}
	return nil, nil
}

func (r *InMemoryProductRepo) Insert(ctx context.Context, e *entity.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[e.Id.String()] = e
	return 1, nil
}

func (r *InMemoryProductRepo) Update(ctx context.Context, e *entity.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[e.Id.String()]; exists {
		r.items[e.Id.String()] = e
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryProductRepo) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[id.String()]; exists {
		delete(r.items, id.String())
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryProductRepo) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := int64(0)
	for _, id := range ids {
		if _, exists := r.items[id.String()]; exists {
			delete(r.items, id.String())
			count++
		}
	}
	return count, nil
}

func (r *InMemoryProductRepo) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p, exists := r.items[id.String()]; exists {
		now := time.Now().UTC()
		p.Deleted = true
		p.DeletedBy = &deletedBy
		p.DeletedAt = &now
		return 1, nil
	}
	return 0, nil
}

func (r *InMemoryProductRepo) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := int64(0)
	for _, p := range r.items {
		if !p.Deleted {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryProductRepo) Upsert(ctx context.Context, e *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[e.Id.String()] = e
	return nil
}

func (r *InMemoryProductRepo) BulkInsert(ctx context.Context, list []*entity.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range list {
		r.items[e.Id.String()] = e
	}
	return int64(len(list)), nil
}

func (r *InMemoryProductRepo) BulkUpdate(ctx context.Context, list []*entity.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := int64(0)
	for _, e := range list {
		if _, exists := r.items[e.Id.String()]; exists {
			r.items[e.Id.String()] = e
			count++
		}
	}
	return count, nil
}

func (r *InMemoryProductRepo) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return nil
}

func (r *InMemoryProductRepo) Search(ctx context.Context, filter *entity.ProductSearchFilter) (*entity.ProductSearchResult, error) {
	// Not implemented for test
	return &entity.ProductSearchResult{}, nil
}

func setupTestServer(t *testing.T) (catalogv1.ProductServiceClient, func()) {
	// gRPC sunucusunu oluştur
	s := grpc.NewServer(
		grpc.UnaryInterceptor(testInterceptor),
	)

	// 1. Önce sahte veritabanını oluştur (RAM'de çalışır, çok hızlıdır)
	mockRepo := NewInMemoryProductRepo()

	// 2. Servisi bu sahte veritabanıyla başlat. Servis bunu gerçek veritabanı sanacak çünkü interface'e uyuyor.
	mockSvc := service.NewProductService(mockRepo)

	// 3. Handler oluştur
	mockHandler := handler.NewProductHandler(mockSvc)

	catalogv1.RegisterProductServiceServer(s, mockHandler)

	// Sunucuyu sanal listener üzerinde başlat (Goroutine içinde)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	// Client bağlantısını oluştur
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := catalogv1.NewProductServiceClient(conn)

	// Teardown (Temizlik) fonksiyonu döndür
	return client, func() {
		err := conn.Close()
		require.NoError(t, err)
		s.Stop()
	}
}

// 3. E2E Test Senaryosu
func TestProductService_E2E(t *testing.T) {
	client, teardown := setupTestServer(t)
	defer teardown()

	userId := uuid.New().String()
	ctx := context.WithValue(context.Background(), interceptor.UserIdKey, userId)
	var createdId string

	// Adım 1: Create Product
	t.Run("Create Product", func(t *testing.T) {
		req := &catalogv1.CreateProductRequest{
			Name:      "MSI Raider GE78",
			BrandId:   1,
			Storyline: ptr("High performance gaming laptop"),
			Price:     3500.00,
		}

		resp, err := client.Create(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Id)
		require.Equal(t, req.Name, resp.Name)
		require.Equal(t, req.BrandId, resp.BrandId)
		require.Equal(t, *req.Storyline, *resp.Storyline)
		require.Equal(t, float32(req.Price), resp.Price)
		createdId = resp.Id // Sonraki testler için Id'yi sakla
	})

	// Adım 2: Get Product
	t.Run("Get Product", func(t *testing.T) {
		req := &catalogv1.GetProductRequest{Id: createdId}
		resp, err := client.Get(ctx, req)

		require.NoError(t, err)
		require.Equal(t, createdId, resp.Id)
		require.Equal(t, "MSI Raider GE78", resp.Name)
		require.Equal(t, int32(1), resp.BrandId)
		require.Equal(t, "High performance gaming laptop", *resp.Storyline)
		require.Equal(t, float32(3500.00), resp.Price)
	})

	// Adım 3: Update Product
	t.Run("Update Product", func(t *testing.T) {
		// Update name and price
		req := &catalogv1.UpdateProductRequest{
			Id:      createdId,
			Name:    "MSI Raider GE78 HX", // Name changed
			BrandId: 1,
			Price:   4000.00, // Price changed
		}

		resp, err := client.Update(ctx, req)
		require.NoError(t, err)
		require.Equal(t, "MSI Raider GE78 HX", resp.Name)
		require.Equal(t, float32(4000.00), resp.Price) // proto'da float32
		require.Equal(t, req.BrandId, resp.BrandId)
		require.Nil(t, resp.Storyline) // Not set in update, assuming cleared or unchanged per service logic
	})

	// Adım 4: List Products
	t.Run("List Products", func(t *testing.T) {
		req := &catalogv1.ListProductsRequest{Limit: 10}
		resp, err := client.List(ctx, req)

		require.NoError(t, err)
		require.GreaterOrEqual(t, len(resp.Items), 1)

		// Listede bizim eleman var mı kontrol et
		found := false
		for _, item := range resp.Items {
			if item.Id == createdId {
				found = true
				break
			}
		}
		require.True(t, found, "Created product should be in the list")
	})

	// Adım 5: Delete Product
	t.Run("Delete Product", func(t *testing.T) {
		req := &catalogv1.DeleteProductRequest{Id: createdId}
		_, err := client.Delete(ctx, req)
		require.NoError(t, err)
	})

	// Adım 6: Verify Deletion (Get should fail)
	t.Run("Verify Deletion", func(t *testing.T) {
		req := &catalogv1.GetProductRequest{Id: createdId}
		_, err := client.Get(ctx, req)

		require.Error(t, err)
		// Hata kodunun NOT_FOUND olduğunu doğrula
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	// Error case: Get non-existing product
	t.Run("Get Non-Existing Product", func(t *testing.T) {
		req := &catalogv1.GetProductRequest{Id: uuid.New().String()}
		_, err := client.Get(ctx, req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	// Error case: Update non-existing product
	t.Run("Update Non-Existing Product", func(t *testing.T) {
		req := &catalogv1.UpdateProductRequest{
			Id:      uuid.New().String(),
			Name:    "Non-Existing",
			BrandId: 1,
			Price:   100.00,
		}
		_, err := client.Update(ctx, req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})

	// Error case: Delete non-existing product
	t.Run("Delete Non-Existing Product", func(t *testing.T) {
		req := &catalogv1.DeleteProductRequest{Id: uuid.New().String()}
		_, err := client.Delete(ctx, req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, codes.NotFound, st.Code())
	})
}
