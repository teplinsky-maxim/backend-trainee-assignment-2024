package cache

import (
	"avito-backend-2024-trainee/internal/entity"
	"errors"
	"sync"
	"time"
)

var (
	ElementDoesNotExistError = errors.New("element does not exists")
)

type compositeKey struct {
	tagId     uint
	featureId uint
}

func newCompositeKey(tagId, featureId uint) compositeKey {
	return compositeKey{
		tagId:     tagId,
		featureId: featureId,
	}
}

type timedItem[T any] struct {
	timeCreated time.Time
	item        T
}

type InMemoryCache struct {
	m      sync.RWMutex
	values map[compositeKey]timedItem[entity.ProductionBanner]
	ttl    time.Duration
}

func NewInMemoryCache(ttl time.Duration) InMemoryCache {
	return InMemoryCache{
		values: make(map[compositeKey]timedItem[entity.ProductionBanner], 50),
		ttl:    ttl,
	}
}

func (r *InMemoryCache) Get(tagId, featureId uint) (entity.ProductionBanner, error) {
	key := newCompositeKey(tagId, featureId)

	r.m.RLock()
	item, ok := r.values[key]
	r.m.RUnlock()

	if !ok {
		return entity.ProductionBanner{}, ElementDoesNotExistError
	}
	now := time.Now().Add(r.ttl)
	if item.timeCreated.Before(now) {
		return item.item, nil
	}

	r.m.Lock()
	delete(r.values, key)
	r.m.Unlock()

	return entity.ProductionBanner{}, ElementDoesNotExistError
}

func (r *InMemoryCache) Set(featureId, tagId uint, banner entity.ProductionBanner) {
	key := newCompositeKey(tagId, featureId)

	r.m.Lock()
	defer r.m.Unlock()

	r.values[key] = timedItem[entity.ProductionBanner]{
		timeCreated: time.Now(),
		item:        banner,
	}
}
