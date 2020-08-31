package sync

import (
	"sync"
)

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]struct{})
	},
}
var l2MapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]map[string]struct{})
	},
}

func getMap() map[string]struct{} {
	return mapPool.Get().(map[string]struct{})
}

func putMap(m map[string]struct{}) {
	for k := range m {
		delete(m, k)
	}
	mapPool.Put(m)
}

func getL2Map() map[string]map[string]struct{} {
	return l2MapPool.Get().(map[string]map[string]struct{})
}

func putL2Map(m map[string]map[string]struct{}) {
	for k, v := range m {
		delete(m, k)
		putMap(v)
	}
	l2MapPool.Put(m)
}
