package data

import (
	"container/list"
	"sync"

	"github.com/cmd4coder/cmd4coder/internal/model"
)

// Cache 缓存管理器
type Cache struct {
	maxSize int
	mu      sync.RWMutex
	items   map[string]*cacheItem
	lru     *list.List
}

// cacheItem 缓存项
type cacheItem struct {
	key     string
	value   interface{}
	element *list.Element
}

// NewCache 创建缓存管理器
func NewCache(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
		items:   make(map[string]*cacheItem),
		lru:     list.New(),
	}
}

// Get 获取缓存项
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// 移到列表前端（最近使用）
	c.lru.MoveToFront(item.element)
	return item.value, true
}

// Set 设置缓存项
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果已存在，更新并移到前端
	if item, ok := c.items[key]; ok {
		item.value = value
		c.lru.MoveToFront(item.element)
		return
	}

	// 如果缓存已满，删除最久未使用的项
	if c.lru.Len() >= c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			c.lru.Remove(oldest)
			delete(c.items, oldest.Value.(*cacheItem).key)
		}
	}

	// 添加新项
	item := &cacheItem{
		key:   key,
		value: value,
	}
	item.element = c.lru.PushFront(item)
	c.items[key] = item
}

// Clear 清空缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	c.lru = list.New()
}

// Size 获取缓存大小
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lru.Len()
}

// SearchCache 搜索结果缓存
type SearchCache struct {
	cache *Cache
}

// NewSearchCache 创建搜索缓存
func NewSearchCache(maxSize int) *SearchCache {
	return &SearchCache{
		cache: NewCache(maxSize),
	}
}

// GetSearchResult 获取搜索结果
func (sc *SearchCache) GetSearchResult(query string) ([]*model.Command, bool) {
	result, ok := sc.cache.Get(query)
	if !ok {
		return nil, false
	}
	return result.([]*model.Command), true
}

// SetSearchResult 设置搜索结果
func (sc *SearchCache) SetSearchResult(query string, commands []*model.Command) {
	sc.cache.Set(query, commands)
}

// Clear 清空搜索缓存
func (sc *SearchCache) Clear() {
	sc.cache.Clear()
}
