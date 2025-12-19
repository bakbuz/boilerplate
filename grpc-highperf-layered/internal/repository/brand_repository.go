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

// BrandRepository ...
type BrandRepository interface {
	GetAll(ctx context.Context) ([]*domain.Brand, error)
	GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error)
	GetById(ctx context.Context, id int32) (*domain.Brand, error)
	Insert(ctx context.Context, e *domain.Brand) error
	Update(ctx context.Context, e *domain.Brand) (int64, error)
	Delete(ctx context.Context, id int32) (int64, error)
	DeleteByIds(ctx context.Context, ids []int32) (int64, error)
	Count(ctx context.Context) (int64, error)

	Upsert(ctx context.Context, e *domain.Brand) error
	BulkInsertOneShot(ctx context.Context, list []*domain.Brand) (int64, error)
	BulkInsert(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error)
	BulkUpdate(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error)
}

type brandRepository struct {
	db *database.DB
}

// NewBrandRepository ...
func NewBrandRepository(db *database.DB) BrandRepository {
	return &brandRepository{db: db}
}

func scanBrand(row pgx.Row) (*domain.Brand, error) {
	e := &domain.Brand{}

	err := row.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Or return specific ErrNotFound if preferred
		}
		return nil, errors.WithMessage(err, rowScanError)
	}
	return e, nil
}

// GetAll ...
func (repo *brandRepository) GetAll(ctx context.Context) ([]*domain.Brand, error) {
	// WARNING: Unbounded query. Added safety limit.
	// TODO: Update interface to support pagination.
	const stmt string = "SELECT * FROM catalog.brands LIMIT 1000"

	rows, err := repo.db.Pool().Query(ctx, stmt)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Brand
	for rows.Next() {
		e := &domain.Brand{}
		err := rows.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetByIds ...
func (repo *brandRepository) GetByIds(ctx context.Context, ids []int32) ([]*domain.Brand, error) {
	if len(ids) == 0 {
		return []*domain.Brand{}, nil
	}

	const stmt string = "SELECT * FROM catalog.brands WHERE id = ANY($1)"

	rows, err := repo.db.Pool().Query(ctx, stmt, ids)
	if err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}
	defer rows.Close()

	var list []*domain.Brand
	for rows.Next() {
		e := &domain.Brand{}
		err := rows.Scan(&e.Id, &e.Name, &e.Slug, &e.Logo, &e.CreatedBy, &e.CreatedAt, &e.UpdatedBy, &e.UpdatedAt)
		if err != nil {
			return nil, errors.WithMessage(err, rowScanError)
		}
		list = append(list, e)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithMessage(err, listQueryRowError)
	}

	return list, nil
}

// GetById ...
func (repo *brandRepository) GetById(ctx context.Context, id int32) (*domain.Brand, error) {
	const stmt string = "SELECT * FROM catalog.brands WHERE id=$1"

	row := repo.db.Pool().QueryRow(ctx, stmt, id)
	item, err := scanBrand(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return item, nil
}

// Insert ...
func (repo *brandRepository) Insert(ctx context.Context, e *domain.Brand) error {
	const command string = `
		INSERT INTO catalog.brands (name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := repo.db.Pool().QueryRow(ctx, command,
		e.Name,
		e.Slug,
		e.Logo,
		e.CreatedBy,
		e.CreatedAt,
	).Scan(&e.Id)

	if err != nil {
		return errors.WithMessage(err, failedToInsert)
	}

	return nil
}

// Update ...
func (repo *brandRepository) Update(ctx context.Context, e *domain.Brand) (int64, error) {
	const command string = `
		UPDATE catalog.brands 
		SET name=$2, slug=$3, logo=$4, updated_by=$5, updated_at=$6 
		WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.Name,
		e.Slug,
		e.Logo,
		e.UpdatedBy,
		e.UpdatedAt,
	)

	if err != nil {
		return -1, errors.WithMessage(err, failedToUpdate)
	}

	return result.RowsAffected(), nil
}

// Delete ...
func (repo *brandRepository) Delete(ctx context.Context, id int32) (int64, error) {
	const command string = `DELETE FROM catalog.brands WHERE id=$1`

	result, err := repo.db.Pool().Exec(ctx, command, id)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDelete)
	}

	return result.RowsAffected(), nil
}

// DeleteByIds ...
func (repo *brandRepository) DeleteByIds(ctx context.Context, ids []int32) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	const command string = `DELETE FROM catalog.brands WHERE id = ANY($1)`

	result, err := repo.db.Pool().Exec(ctx, command, ids)
	if err != nil {
		return -1, errors.WithMessage(err, failedToDeletes)
	}

	return result.RowsAffected(), nil
}

// Count ...
func (repo *brandRepository) Count(ctx context.Context) (int64, error) {
	return repo.db.Count(ctx, "SELECT COUNT(*) FROM catalog.brands")
}

// Upsert ...
func (repo *brandRepository) Upsert(ctx context.Context, e *domain.Brand) error {
	// If Id is missing, treating as Insert creates a record with ID 0 on some systems
	// or fails validation. For auto-increment, we must delegate to Insert.
	if e.Id == 0 {
		return repo.Insert(ctx, e)
	}

	const command string = `
		INSERT INTO catalog.brands (id, name, slug, logo, created_by, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET name=EXCLUDED.name,
			slug=EXCLUDED.slug, 
			logo=EXCLUDED.logo, 
			updated_by=$7, 
			updated_at=$8`

	_, err := repo.db.Pool().Exec(ctx, command,
		e.Id,
		e.Name,
		e.Slug,
		e.Logo,
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

// BulkInsertOneShot ...
func (repo *brandRepository) BulkInsertOneShot(ctx context.Context, list []*domain.Brand) (int64, error) {
	if len(list) == 0 {
		return 0, nil
	}

	// 1. Tüm batch için TEK BİR zaman damgası oluştur (Consistency)
	now := time.Now().UTC()

	// 2. pgx'in beklediği veri yapısına (slice of slices) dönüştür
	// Kapasiteyi önceden belirlemek (make) memory allocation maliyetini düşürür.
	rows := make([][]any, len(list))
	for i, e := range list {
		e.CreatedAt = now
		rows[i] = []any{e.Name, e.Slug, e.Logo, e.CreatedBy, e.CreatedAt}
	}

	// 3. PostgreSQL COPY Protokolü ile Veriyi Akıt
	// Bu işlem standart INSERT'ten yaklaşık 5-10 kat daha hızlıdır.
	count, err := repo.db.Pool().CopyFrom(
		ctx,
		pgx.Identifier{"catalog", "brands"},
		[]string{"name", "slug", "logo", "created_by", "created_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return 0, errors.WithMessage(err, failedToBulkInsert)
	}

	return count, nil
}

// batchSize parametresi ile chunk boyutunu belirleyebilirsiniz.
// batchSize = 0 ise varsayılan BatchSize kullanılır.
func (repo *brandRepository) BulkInsert(ctx context.Context, brands []*domain.Brand, batchSize int) (int64, error) {
	// 1. Veri yoksa işlem yapma
	if len(brands) == 0 {
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
	for i := 0; i < len(brands); i += batchSize {
		end := i + batchSize
		if end > len(brands) {
			end = len(brands)
		}

		chunk := brands[i:end]
		rows := make([][]any, 0, len(chunk))

		for i, e := range chunk {
			e.CreatedAt = now
			rows[i] = []any{e.Name, e.Slug, e.Logo, e.CreatedBy, e.CreatedAt}
		}

		// 5. Copy Protokolü ile Yazma (High Performance Insert)
		count, err := repo.db.Pool().CopyFrom(
			ctx,
			pgx.Identifier{"catalog", "brands"},
			[]string{"name", "slug", "logo", "created_by", "created_at"},
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
func (repo *brandRepository) BulkUpdate(ctx context.Context, list []*domain.Brand, batchSize int) (int64, error) {
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
		return list[i].Id < list[j].Id
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
		ids := make([]int32, count)
		names := make([]string, count)
		slugs := make([]string, count)
		logos := make([]string, count)
		updatedBys := make([]*uuid.UUID, count)
		updatedAts := make([]*time.Time, count)

		for k, e := range batch {
			ids[k] = int32(e.Id)
			names[k] = e.Name
			slugs[k] = e.Slug
			logos[k] = e.Logo
			updatedBys[k] = e.UpdatedBy
			updatedAts[k] = e.UpdatedAt
		}

		// UNNEST Sorgusu
		// $1...$6 parametreleri o anki batch'in arrayleridir.
		const query = `
			UPDATE catalog.brands b
			SET 
				name = data.name,
				slug = data.slug,
				logo = data.logo,
				updated_by = data.updated_by,
				updated_at = data.updated_at
			FROM (
				SELECT * FROM UNNEST(
					$1::int[], 
					$2::text[], 
					$3::text[], 
					$4::text[], 
					$5::uuid[], 
					$6::timestamp[]
				) AS t(id, name, slug, logo, updated_by, updated_at)
			) AS data
			WHERE b.id = data.id
		`

		// Batch'i Transaction context'i ile çalıştır
		cmdTag, err := tx.Exec(ctx, query, ids, names, slugs, logos, updatedBys, updatedAts)
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
