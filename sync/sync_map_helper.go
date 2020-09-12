package sync

import (
	gosync "sync"
)

func EraseSyncMap(m *gosync.Map) {
	m.Range(func(key interface{}, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
