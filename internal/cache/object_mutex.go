package cache

import (
	"context"
	"sync"
	"time"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/common"
)

const collectIntervalSec = 1

type DataItem struct {
	locked bool
	mutex  sync.Mutex
}

type ObjectMutex[K comparable] struct {
	ctx     context.Context
	mutexes sync.Map
	mut     sync.RWMutex
	logger  common.Logger
}

func NewObjectMutex[K comparable](ctx context.Context, logger common.Logger) *ObjectMutex[K] {
	obj := &ObjectMutex[K]{
		ctx:    ctx,
		logger: logger,
	}
	obj.startGC()
	return obj
}

func (m *ObjectMutex[K]) Lock(key K) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	data, _ := m.mutexes.LoadOrStore(key, &DataItem{})

	syncObj, ok := data.(*DataItem)

	if !ok {
		panic("lock: unknown type of object")
	}

	syncObj.mutex.Lock()
	syncObj.locked = true
}

func (m *ObjectMutex[K]) Unlock(key K) {
	data, ok := m.mutexes.Load(key)
	if !ok {
		panic("unlock: key not found")
	}

	syncObj, ok := data.(*DataItem)
	if !ok {
		panic("unlock: unknown type of object")
	}

	if !syncObj.locked {
		panic("unlock: key not locked: wrong lock/unlock sequence")
	}

	syncObj.locked = false
	syncObj.mutex.Unlock()
}

func (m *ObjectMutex[K]) startGC() {
	go func() {
		ticker := time.NewTicker(collectIntervalSec * time.Second)

		for {
			select {
			case <-ticker.C:
				m.collectKeys()
			case <-m.ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *ObjectMutex[K]) collectKeys() {
	var lockedCnt, totalCnt int
	keysToDelete := make([]K, 0)
	m.mut.Lock()
	defer m.mut.Unlock()

	m.mutexes.Range(func(k any, v any) bool {
		totalCnt++
		key, ok := k.(K)
		if !ok {
			panic("collectKeys: unknown type of key")
		}
		syncObj, ok := v.(*DataItem)
		if !ok {
			panic("collectKeys: unknown type of value")
		}

		if !syncObj.locked {
			lockedCnt++
			keysToDelete = append(keysToDelete, key)
		}

		return true
	})

	for _, key := range keysToDelete {
		m.mutexes.Delete(key)
	}
	// m.logger.Debug(fmt.Sprintf("collected %v unused keys of %v", lockedCnt, totalCnt))
}
