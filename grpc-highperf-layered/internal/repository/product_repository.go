package repository

import (
	"codegen/internal/domain"
	"codegen/internal/infrastructure/database"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

// ProductRepository ...
type ProductRepository interface {
	GetAll(ctx context.Context) ([]*domain.Product, error)
	GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Insert(ctx context.Context, e *domain.Product) (int64, error)
	Update(ctx context.Context, e *domain.Product) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) (int64, error)
	DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error)
	Count(ctx context.Context) (int64, error)

	Upsert(ctx context.Context, e *domain.Product) error
	BulkInsertOneShot(ctx context.Context, list []*domain.Product) (int64, error)
	BulkInsert(ctx context.Context, list []*domain.Product, batchSize int) (int64, error)
	BulkUpdate(ctx context.Context, list []*domain.Product, batchSize int) (int64, error)
	Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error)
}

type productRepository struct {
	db *database.DB
}

// NewProductRepository ...
func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{db: db}
}

func scanProduct(row pgx.Row) (*domain.Product, error) {
	e := &domain.Product{}

	err := row.Scan(&e.Id, &e.BrandId, &e.Name, &e.Sku, &e.Summary, &e.Storyline, &e.StockQuantity, &e.Price, &e.Deleted, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt, &e.DeletedBy, &e.DeletedAt)
	if err == pgx.ErrNoRows { // sql: no rows in result set
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithMessage(err, rowScanError)
	}
	return e, nil
}

// GetAll ...
func (repo *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, product)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetByIds ...
func (repo *productRepository) GetByIds(ctx context.Context, ids []uuid.UUID) ([]*domain.Product, error) {
	if len(ids) == 0 {
		return []*domain.Product{}, nil
	}

	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false AND id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Product
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, product)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetById ...
func (repo *productRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	const stmt string = "SELECT id, brand_id, name, sku, summary, storyline, stock_quantity, price::numeric, deleted, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM catalog.products WHERE deleted=false AND id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := scanProduct(row)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *productRepository) Insert(ctx context.Context, e *domain.Product) (int64, error) {
	const command string = `
		INSERT INTO catalog.products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.Deleted,
		e.CreatedBy,
		e.CreatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToInsert)
	}

	return result.RowsAffected(), nil
}

// Update ...
func (repo *productRepository) Update(ctx context.Context, e *domain.Product) (int64, error) {
	const command string = `
		UPDATE catalog.products 
		SET brand_id=$2, name=$3, sku=$4, summary=$5, storyline=$6, stock_quantity=$7, price=$8, updated_by=$9, updated_at=$10 
		WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToUpdate)
	}

	return result.RowsAffected(), nil
}

// Delete ...
func (repo *productRepository) Delete(ctx context.Context, id uuid.UUID) (int64, error) {
	const command string = `DELETE FROM catalog.products WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDelete)
	}

	return result.RowsAffected(), nil
}

// DeleteByIds ...
func (repo *productRepository) DeleteByIds(ctx context.Context, ids []uuid.UUID) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	const command string = `DELETE FROM catalog.products WHERE id = ANY($1)`

	result, err := repo.db.Pool().Exec(ctx, command, ids)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDeletes)
	}

	return result.RowsAffected(), nil
}

// SoftDelete ...
func (repo *productRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) (int64, error) {
	const command string = `UPDATE catalog.products SET deleted=true, deleted_by=$2, deleted_at=$3 WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id, deletedBy, time.Now().UTC())
	if err != nil {
		return -1, errors.WithMessage(err, failedToSoftDelete)
	}

	return result.RowsAffected(), nil
}

// Count ...
func (repo *productRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM catalog.products WHERE deleted=false")
}

// Upsert ...
func (repo *productRepository) Upsert(ctx context.Context, e *domain.Product) error {
	if e.Id == uuid.Nil {
		var err error
		if e.Id, err = uuid.NewV7(); err != nil {
			return errors.WithMessage(err, "failed to generate id")
		}
	}

	const command string = `
		INSERT INTO catalog.products (id, brand_id, name, sku, summary, storyline, stock_quantity, price, deleted, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE 
		SET brand_id=EXCLUDED.brand_id,
			name=EXCLUDED.name,
			sku=EXCLUDED.sku,
			summary=EXCLUDED.summary,
			storyline=EXCLUDED.storyline,
			stock_quantity=EXCLUDED.stock_quantity,
			price=EXCLUDED.price,
			deleted=EXCLUDED.deleted,
			updated_by=$12,
			updated_at=$13`

	_, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.BrandId,
		e.Name,
		e.Sku,
		e.Summary,
		e.Storyline,
		e.StockQuantity,
		e.Price,
		e.Deleted,
		e.CreatedBy,
		e.CreatedAt,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return errors.WithMessage(err, failedToUpsert)
	}

	return nil
}

func (repo *productRepository) BulkInsertOneShot(ctx context.Context, list []*domain.Product) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	// 1. Tüm batch için TEK BİR zaman damgası oluştur (Consistency)
	now := time.Now().UTC()

	// 2. pgx'in beklediği veri yapısına (slice of slices) dönüştür
	// Kapasiteyi önceden belirlemek (make) memory allocation maliyetini düşürür.
	rows := make([][]any, len(list))
	for i, e := range list {
		if e.Id == uuid.Nil {
			newId, err := uuid.NewV7()
			if err != nil {
				return 0, errors.WithMessage(err, "failed to generate id")
			}
			e.Id = newId
		}
		e.CreatedAt = now
		rows[i] = []any{e.Id, e.BrandId, e.Name, e.Sku, e.Summary, e.Storyline, e.StockQuantity, e.Price, e.Deleted, e.CreatedBy, e.CreatedAt}
	}

	// 3. PostgreSQL COPY Protokolü ile Veriyi Akıt
	// Bu işlem standart INSERT'ten yaklaşık 5-10 kat daha hızlıdır.
	count, err := repo.db.Pool().CopyFrom(
		ctx,
		pgx.Identifier{"catalog", "products"},
		[]string{"id", "brand_id", "name", "sku", "summary", "storyline", "stock_quantity", "price", "deleted", "created_by", "created_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkInsert)
	}

	return count, nil
}

// batchSize parametresi ile chunk boyutunu belirleyebilirsiniz.
// batchSize = 0 ise varsayılan BatchSize kullanılır.
func (repo *productRepository) BulkInsert(ctx context.Context, products []*domain.Product, batchSize int) (int64, error) {
	// 1. Veri yoksa işlem yapma
	if len(products) == 0 {
		return 0, nil
	}

	if batchSize == 0 {
		batchSize = defaultBatchSize
	}

	var totalCount int64 = 0

	// 2. Transaction Başlat
	// Bu çok kritiktir. 5. parçada hata alırsak, önceki 4 parçayı da geri almalıyız.
	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("transaction başlatılamadı: %w", err)
	}
	// Fonksiyondan çıkarken hata varsa Rollback yapar, Commit edildiyse etkisizdir.
	defer tx.Rollback(ctx)

	// 3. Ortak Zaman Damgası
	// Tüm batch'in aynı saniyede oluştuğunu görmek loglama açısından temizdir.
	now := time.Now().UTC()

	// 4. Chunking (Parçalama) Döngüsü
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		chunk := products[i:end]
		rows := make([][]any, len(chunk))

		for i, e := range chunk {
			e.CreatedAt = now
			rows[i] = []any{e.Id, e.BrandId, e.Name, e.Sku, e.Summary, e.Storyline, e.StockQuantity, e.Price, e.Deleted, e.CreatedBy, e.CreatedAt}
		}

		// 5. Copy Protokolü ile Yazma (High Performance Insert)
		count, err := repo.db.Pool().CopyFrom(
			ctx,
			pgx.Identifier{"catalog", "products"},
			[]string{"id", "brand_id", "name", "sku", "summary", "storyline", "stock_quantity", "price", "deleted", "created_by", "created_at"},
			pgx.CopyFromRows(rows),
		)

		if err != nil {
			// Hata detayını yakala (Örn: Unique Constraint ihlali varsa)
			if pgErr, ok := err.(*pgconn.PgError); ok {
				return 0, fmt.Errorf("bulk insert hatası (SQL State: %s): %s", pgErr.Code, pgErr.Message)
			}
			return 0, fmt.Errorf("bulk insert bilinmeyen hata: %w", err)
		}

		totalCount += count
	}

	// 6. Transaction Commit
	// Buraya kadar geldiyse tüm parçalar hatasız yazılmıştır.
	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("transaction commit hatası: %w", err)
	}

	return totalCount, nil
}

// BulkUpdate, büyük veri setlerini güvenli parçalar halinde günceller.
func (repo *productRepository) BulkUpdate(ctx context.Context, list []*domain.Product, batchSize int) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	if batchSize == 0 {
		batchSize = defaultBatchSize
	}

	// 1. ÖN HAZIRLIK: DEADLOCK KORUMASI
	// Veriyi işlemeye başlamadan önce ID'ye göre sıralıyoruz.
	// Bu, tüm batch'lerin her zaman aynı sırada kilit almasını garanti eder.
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})

	// 2. TRANSACTION BAŞLATMA
	// Tüm parçaların (chunks) ya hep ya hiç (all-or-nothing) mantığıyla
	// işlenmesi için tek bir transaction başlatıyoruz.
	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to begin transaction")
	}
	// Hata durumunda rollback, başarı durumunda commit öncesi no-op olur.
	defer tx.Rollback(ctx)

	var totalAffected int64 = 0

	// 3. CHUNKING DÖNGÜSÜ
	for i := 0; i < len(list); i += batchSize {
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		batch := list[i:end]

		// Batch için slice'ları hazırla (Memory Allocation optimization)
		// Her turda yeniden oluşturuyoruz ki GC eski batch'i temizleyebilsin.
		count := len(batch)
		ids := make([]uuid.UUID, count)
		brandIds := make([]int32, count)
		names := make([]string, count)
		skus := make([]*string, count)
		summaries := make([]*string, count)
		storylines := make([]*string, count)
		stockQuantities := make([]int32, count)
		prices := make([]float64, count)
		updatedBys := make([]*uuid.UUID, count)
		updatedAts := make([]*time.Time, count)

		for k, e := range batch {
			ids[k] = e.Id
			brandIds[k] = int32(e.BrandId)
			names[k] = e.Name
			skus[k] = e.Sku
			summaries[k] = e.Summary
			storylines[k] = e.Storyline
			stockQuantities[k] = int32(e.StockQuantity)
			prices[k] = e.Price
			updatedBys[k] = e.UpdatedBy
			updatedAts[k] = e.UpdatedAt
		}

		// UNNEST Sorgusu
		// $1...$10 parametreleri o anki batch'in arrayleridir.
		const query = `
			UPDATE catalog.products p
			SET 
				brand_id = data.brand_id,
				name = data.name,
				sku = data.sku,
				summary = data.summary,
				storyline = data.storyline,
				stock_quantity = data.stock_quantity,
				price = data.price,
				updated_by = data.updated_by,
				updated_at = data.updated_at
			FROM (
				SELECT * FROM UNNEST(
					$1::uuid[], 
					$2::int[], 
					$3::text[], 
					$4::text[], 
					$5::text[], 
					$6::text[], 
					$7::int[], 
					$8::numeric[], 
					$9::uuid[], 
					$10::timestamp[]
				) AS t(id, brand_id, name, sku, summary, storyline, stock_quantity, price, updated_by, updated_at)
			) AS data
			WHERE p.id = data.id
		`

		// Batch'i Transaction context'i ile çalıştır
		cmdTag, err := tx.Exec(ctx, query, ids, brandIds, names, skus, summaries, storylines, stockQuantities, prices, updatedBys, updatedAts)
		if err != nil {
			// Hata detayını hangi batch'te olduğunu belirterek dönüyoruz
			return 0, errors.WithMessagef(err, "failed to update batch starting at index %d", i)
		}

		totalAffected += cmdTag.RowsAffected()
	}

	// 4. COMMIT
	if err := tx.Commit(ctx); err != nil {
		return 0, errors.WithMessage(err, "failed to commit bulk update transaction")
	}

	return totalAffected, nil
}

/*
func (repo *productRepository) runInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := repo.db.Pool().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Context'e tx ekle
	txCtx := context.WithValue(ctx, txContextKey, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
*/

func (repo *productRepository) Search(ctx context.Context, filter *domain.ProductSearchFilter) (*domain.ProductSearchResult, error) {
	// 1. Base Query
	baseQuery := " FROM catalog.products p JOIN catalog.brands b ON p.brand_id = b.id WHERE p.deleted=false "
	var args []any
	argId := 1

	// 2. Dynamic Filters
	if filter.BrandId > 0 {
		baseQuery += fmt.Sprintf(" AND p.brand_id = $%d", argId)
		args = append(args, filter.BrandId)
		argId++
	}

	if filter.Name != "" {
		baseQuery += fmt.Sprintf(" AND p.name ILIKE '%%' || $%d || '%%'", argId)
		args = append(args, filter.Name)
		argId++
	}

	if filter.LastSeenId != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND p.id < $%d", argId)
		args = append(args, filter.LastSeenId)
		argId++
	}

	// 3. Count Query
	// We rebuild count query to exclude cursor/pagination filters if needed,
	// but here we only have LastSeenId as pagination filter.
	// As discussed, we count based on business filters (Brand, Name).

	countBaseQuery := "FROM catalog.products p WHERE p.deleted=false"
	var countArgs []any
	countArgId := 1

	if filter.BrandId > 0 {
		countBaseQuery += fmt.Sprintf(" AND p.brand_id = $%d", countArgId)
		countArgs = append(countArgs, filter.BrandId)
		countArgId++
	}
	if filter.Name != "" {
		countBaseQuery += fmt.Sprintf(" AND p.name ILIKE '%%' || $%d || '%%'", countArgId)
		countArgs = append(countArgs, filter.Name)
		countArgId++
	}

	finalCountQuery := "SELECT COUNT(*) " + countBaseQuery

	var total int64
	err := repo.db.Pool().QueryRow(ctx, finalCountQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, errors.WithMessage(err, failedToCount)
	}

	if total == 0 {
		return &domain.ProductSearchResult{Total: 0, Items: []domain.ProductSummary{}}, nil
	}

	// 4. Data Query
	// Limit default
	limit := 10
	if filter.Limit > 0 {
		limit = filter.Limit
	}

	// Select columns
	selectQuery := "SELECT p.id, p.name, p.price::numeric, p.brand_id, b.name AS brand_name " + baseQuery + " ORDER BY p.id DESC LIMIT " + fmt.Sprintf("$%d", argId)
	args = append(args, limit)

	rows, err := repo.db.Pool().Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var items []domain.ProductSummary
	for rows.Next() {
		e := domain.ProductSummary{}
		// Scan matches the SELECT columns order
		err := rows.Scan(&e.Id, &e.Name, &e.Price, &e.BrandId, &e.BrandName)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		items = append(items, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return &domain.ProductSearchResult{
		Total: total,
		Items: items,
	}, nil
}
