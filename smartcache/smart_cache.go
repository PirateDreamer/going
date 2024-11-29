package smartcache

// 实现go本地缓存，同时支持单个值过期和多个值同时过期，00gc，fe分片存储，支持缓存淘汰策略FIFO和LRU
import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// 缓存项结构
type cacheItem struct {
	key        interface{}
	value      interface{}
	expiration int64
}

// 缓存分片结构
type cacheShard struct {
	items             map[interface{}]*list.Element
	lruList           *list.List
	lock              sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	maxItems          int
}

// SmartCache 结构
type SmartCache struct {
	shards         []*cacheShard
	shardCount     int
	evictionPolicy string // "FIFO" 或 "LRU"
}

// 创建新的SmartCache
func NewSmartCache(shardCount, maxItemsPerShard int, defaultExpiration, cleanupInterval time.Duration, evictionPolicy string) *SmartCache {
	sc := &SmartCache{
		shards:         make([]*cacheShard, shardCount),
		shardCount:     shardCount,
		evictionPolicy: evictionPolicy,
	}
	for i := 0; i < shardCount; i++ {
		sc.shards[i] = &cacheShard{
			items:             make(map[interface{}]*list.Element),
			lruList:           list.New(),
			defaultExpiration: defaultExpiration,
			cleanupInterval:   cleanupInterval,
			maxItems:          maxItemsPerShard,
		}
		go sc.shards[i].startCleanupTimer()
	}
	return sc
}

// 获取分片索引
func (sc *SmartCache) getShard(key interface{}) *cacheShard {
	hash := fnv32(key)
	return sc.shards[hash%uint32(sc.shardCount)]
}

// 设置缓存项
func (sc *SmartCache) Set(key, value interface{}, expiration time.Duration) {
	shard := sc.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()

	var exp int64
	if expiration == 0 {
		exp = 0
	} else {
		exp = time.Now().Add(expiration).UnixNano()
	}

	item := &cacheItem{
		key:        key,
		value:      value,
		expiration: exp,
	}

	if element, exists := shard.items[key]; exists {
		if sc.evictionPolicy == "LRU" {
			shard.lruList.MoveToFront(element)
		}
		element.Value = item
	} else {
		if len(shard.items) >= shard.maxItems {
			sc.evict(shard)
		}
		var element *list.Element
		if sc.evictionPolicy == "FIFO" {
			element = shard.lruList.PushBack(item)
		} else {
			element = shard.lruList.PushFront(item)
		}
		shard.items[key] = element
	}
}

// 获取缓存项
func (sc *SmartCache) Get(key interface{}) (interface{}, bool) {
	shard := sc.getShard(key)
	shard.lock.RLock()
	defer shard.lock.RUnlock()

	element, exists := shard.items[key]
	if !exists {
		return nil, false
	}

	item := element.Value.(*cacheItem)
	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	if sc.evictionPolicy == "LRU" {
		shard.lruList.MoveToFront(element)
	}

	return item.value, true
}

// 删除缓存项
func (sc *SmartCache) Delete(key interface{}) {
	shard := sc.getShard(key)
	shard.lock.Lock()
	defer shard.lock.Unlock()

	if element, exists := shard.items[key]; exists {
		shard.lruList.Remove(element)
		delete(shard.items, key)
	}
}

// 清理过期项
func (cs *cacheShard) cleanupExpired() {
	now := time.Now().UnixNano()
	cs.lock.Lock()
	defer cs.lock.Unlock()

	for key, element := range cs.items {
		item := element.Value.(*cacheItem)
		if item.expiration > 0 && now > item.expiration {
			cs.lruList.Remove(element)
			delete(cs.items, key)
		}
	}
}

// 启动清理定时器
func (cs *cacheShard) startCleanupTimer() {
	ticker := time.NewTicker(cs.cleanupInterval)
	for range ticker.C {
		cs.cleanupExpired()
	}
}

// 淘汰缓存项
func (sc *SmartCache) evict(shard *cacheShard) {
	var element *list.Element
	if sc.evictionPolicy == "FIFO" {
		element = shard.lruList.Front()
	} else {
		element = shard.lruList.Back()
	}
	if element != nil {
		item := element.Value.(*cacheItem)
		delete(shard.items, item.key)
		shard.lruList.Remove(element)
	}
}

// FNV-1a 哈希函数
func fnv32(key interface{}) uint32 {
	h := uint32(2166136261)
	bytes := []byte(fmt.Sprintf("%v", key))
	for _, b := range bytes {
		h ^= uint32(b)
		h *= 16777619
	}
	return h
}
