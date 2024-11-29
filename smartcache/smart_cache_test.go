package smartcache

import (
	"fmt"
	"testing"
	"time"
)

func TestSmartCache(t *testing.T) {
	// 创建一个新的SmartCache实例
	cache := NewSmartCache(16, 100, 5*time.Minute, 10*time.Minute, "LRU")

	// 设置缓存项
	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 2*time.Minute)

	// 获取缓存项
	if value, found := cache.Get("key1"); found {
		fmt.Printf("Found key1: %v\n", value)
	}

	// 删除缓存项
	cache.Delete("key2")

	// 检查删除的项
	if _, found := cache.Get("key2"); !found {
		fmt.Println("key2 was successfully deleted")
	}

	// 等待过期
	time.Sleep(2 * time.Minute)

	// 检查过期的项
	if _, found := cache.Get("key1"); !found {
		fmt.Println("key1 has expired")
	}
}
