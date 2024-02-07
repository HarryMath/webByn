package repository

import (
	"errors"
	"sync"
	"webByn/src/util"
)

// Repository represents a generic repository.
type Repository[T Entity] struct {
	data        []T
	transaction *sync.Mutex
}

// NewRepository creates a new instance of Repository.
func NewRepository[T Entity]() *Repository[T] {
	return &Repository[T]{
		data: make([]T, 0),
	}
}

// Add adds an item to the repository.
func (r *Repository[T]) Add(item T) error {
	var uid = item.GetUid()
	var duplicate, _ = r.GetById(uid, true)
	if duplicate == nil {
		r.data = append(r.data, item)
		return nil
	}
	return errors.New("Conflict creating entity")
}

func (r *Repository[T]) GetById(uid interface{}, silentMode bool) (*T, error) {
	criteria := func(t T) bool { return t.GetUid() == uid }
	return r.FindOneBy(criteria, silentMode)
}

// FindOneBy retrieves the first element by criteria.
func (r *Repository[T]) FindOneBy(criteria func(T) bool, silentMode bool) (*T, error) {
	r.transaction.Lock()
	for _, item := range r.data {
		// Assume items have an "ID" field
		if criteria(item) {
			r.transaction.Unlock()
			return &item, nil
		}
	}
	r.transaction.Unlock()

	return nil, util.Ternary(silentMode, nil, errors.New("Entity not found"))
}
