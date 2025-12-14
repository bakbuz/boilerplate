package utils

import (
	"sync"
)

// ObjectPool for zero-allocation object reuse
type ObjectPool[T any] struct {
	pool sync.Pool
}

func NewObjectPool[T any](newFunc func() T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: sync.Pool{
			New: func() any {
				return newFunc()
			},
		},
	}
}

func (p *ObjectPool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *ObjectPool[T]) Put(obj T) {
	// Reset object before putting back to pool
	if resettable, ok := any(obj).(interface{ Reset() }); ok {
		resettable.Reset()
	}
	p.pool.Put(obj)
}

// ByteSlicePool for reducing allocations
var ByteSlicePool = sync.Pool{
	New: func() any {
		return make([]byte, 0, 1024)
	},
}

func GetByteSlice() []byte {
	return ByteSlicePool.Get().([]byte)
}

func PutByteSlice(b []byte) {
	b = b[:0]
	ByteSlicePool.Put(b)
}
