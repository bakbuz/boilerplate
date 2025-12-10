package postgres

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/rs/zerolog"

    "github.com/yourusername/grpc-highperf-backend/internal/domain"
)

type ProductRepository struct {
    pool *pgxpool.Pool
    log  zerolog.Logger
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
    return &ProductRepository{
        pool: pool,
        log:  zerolog.Nop(),
    }
}

// Create product with minimal allocations
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
    query := `
        INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

    now := time.Now().UTC()
    product.ID = generateUUID()
    product.CreatedAt = now
    product.UpdatedAt = now

    err := r.pool.QueryRow(ctx, query,
        product.ID,
        product.Name,
        product.Description,
        product.Price,
        product.Stock,
        product.CreatedAt,
        product.UpdatedAt,
    ).Scan(&product.ID)

    if err != nil {
        return fmt.Errorf("failed to create product: %w", err)
    }

    return nil
}

// GetByID with optimized query
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
    query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products 
        WHERE id = $1 AND deleted_at IS NULL
        LIMIT 1
    `

    var product domain.Product
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &product.ID,
        &product.Name,
        &product.Description,
        &product.Price,
        &product.Stock,
        &product.CreatedAt,
        &product.UpdatedAt,
    )

    if err == pgx.ErrNoRows {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get product: %w", err)
    }

    return &product, nil
}

// List with pagination and filtering
func (r *ProductRepository) List(ctx context.Context, page, pageSize int, filter string) ([]*domain.Product, int, error) {
    // Count total
    countQuery := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
    if filter != "" {
        countQuery += ` AND (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')`
    }

    var total int
    var err error
    if filter != "" {
        err = r.pool.QueryRow(ctx, countQuery, filter).Scan(&total)
    } else {
        err = r.pool.QueryRow(ctx, countQuery).Scan(&total)
    }
    if err != nil {
        return nil, 0, err
    }

    // Get paginated results
    offset := (page - 1) * pageSize
    query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products 
        WHERE deleted_at IS NULL
    `

    args := []interface{}{}
    argPos := 1

    if filter != "" {
        query += ` AND (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')`
        args = append(args, filter)
        argPos++
    }

    query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprint(argPos) + ` OFFSET $` + fmt.Sprint(argPos+1)
    args = append(args, pageSize, offset)

    rows, err := r.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    products := make([]*domain.Product, 0, pageSize)
    for rows.Next() {
        var product domain.Product
        if err := rows.Scan(
            &product.ID,
            &product.Name,
            &product.Description,
            &product.Price,
            &product.Stock,
            &product.CreatedAt,
            &product.UpdatedAt,
        ); err != nil {
            return nil, 0, err
        }
        products = append(products, &product)
    }

    return products, total, nil
}