package repository

import (
	"errors"
	"sync"
	"webByn/src/util"
)

// Repository represents a generic repository.
type Repository[T Entity[Json], Json any] struct {
	data        []T
	transaction sync.Mutex
}

// NewRepository creates a new instance of Repository.
func NewRepository[T Entity[Json], Json any]() *Repository[T, Json] {
	return &Repository[T, Json]{
		data: make([]T, 0),
	}
}

// Add adds an item to the repository.
func (r *Repository[T, Json]) Add(item T) error {
	var uid = item.GetUid()
	var duplicate, _ = r.GetById(uid, true)
	if duplicate == nil {
		r.data = append(r.data, item)
		return nil
	}
	return errors.New("conflict creating entity")
}

func (r *Repository[T, Json]) GetById(uid interface{}, silentMode bool) (*T, error) {
	criteria := func(t T) bool { return t.GetUid() == uid }
	return r.FindOneBy(criteria, silentMode)
}

// FindOneBy retrieves the first element by criteria.
func (r *Repository[T, Json]) FindOneBy(criteria func(T) bool, silentMode bool) (*T, error) {
	r.transaction.Lock()
	for _, item := range r.data {
		if criteria(item) {
			r.transaction.Unlock()
			return &item, nil
		}
	}
	r.transaction.Unlock()

	return nil, util.Ternary(silentMode, nil, errors.New("entity not found"))
}

// GelAll returns all accounts calling ToJson for each entity
func (r *Repository[T, Json]) GelAll() []Json {
	transform := func(entity T) Json {
		return entity.ToJson()
	}
	return util.Map(r.data, transform)
}
