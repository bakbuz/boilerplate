package e2e_test

import (
	catalogv1 "codegen/api/gen/catalog/v1"
	"codegen/internal/service"
	"context"
	"net"
	"sync"
	"testing"

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

// 2. Test Ortamını Ayağa Kaldıran Yardımcı Fonksiyon
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type InMemoryProductRepo struct {
	items map[string]*catalogv1.Product
	mu    sync.RWMutex
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		items: make(map[string]*catalogv1.Product),
	}
}

func setupTestServer(t *testing.T) (catalogv1.ProductServiceClient, func()) {
	// gRPC sunucusunu oluştur
	s := grpc.NewServer()

	// 1. Önce sahte veritabanını oluştur (RAM'de çalışır, çok hızlıdır)
	mockRepo := NewInMemoryProductRepo()

	// 2. Servisi bu sahte veritabanıyla başlat. Servis bunu gerçek veritabanı sanacak çünkü interface'e uyuyor.
	mockSvc := service.NewProductService(mockRepo)

	catalogv1.RegisterProductServiceServer(s, mockSvc)

	// Sunucuyu sanal listener üzerinde başlat (Goroutine içinde)
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	// Client bağlantısını oluştur
	conn, err := grpc.NewClient("bufnet",
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

	ctx := context.Background()
	var createdId string

	// Adım 1: Create Product
	t.Run("Create Product", func(t *testing.T) {
		req := &catalogv1.CreateProductRequest{
			Name:      "MSI Raider GE78",
			Storyline: strPtr("High performance gaming laptop"),
			Price:     3500.00,
		}

		resp, err := client.Create(ctx, req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Id)
		require.Equal(t, req.Name, resp.Name)

		createdId = resp.Id // Sonraki testler için Id'yi sakla
	})

	// Adım 2: Get Product
	t.Run("Get Product", func(t *testing.T) {
		req := &catalogv1.GetProductRequest{Id: createdId}
		resp, err := client.Get(ctx, req)

		require.NoError(t, err)
		require.Equal(t, createdId, resp.Id)
		require.Equal(t, "MSI Raider GE78", resp.Name)
	})

	// Adım 3: Update Product
	t.Run("Update Product", func(t *testing.T) {
		// Sadece fiyatı güncelleyelim
		req := &catalogv1.UpdateProductRequest{
			Id:    createdId,
			Name:  "MSI Raider GE78 HX", // İsim değişti
			Price: 4000.00,              // Fiyat değişti
		}

		resp, err := client.Update(ctx, req)
		require.NoError(t, err)
		require.Equal(t, "MSI Raider GE78 HX", resp.Name)
		require.Equal(t, 4000.00, resp.Price) // proto'da float ise
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
}
