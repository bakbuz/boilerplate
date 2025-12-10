package domain

type Product struct {
	ID    int32
	Name  string
	Price float64
	Stock int32
}
type Order struct {
	ID        int32
	ProductID int32
	Quantity  int32
	CreatedAt string // Basitlik için string, gerçek uygulamada time.Time kullanılır
}
