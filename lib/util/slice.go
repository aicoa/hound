package util

import (
	"sync"
)

type (
	// SafeSlice 为并发安全的切片结构体
	SafeSlice struct {
		sync.RWMutex         // 嵌入读写锁以保证并发安全
		items        []Items // 存储 Items 类型的切片
	}

	// Items 为切片存储的元素，包含值和额外的整型字段
	Items struct {
		num   int // 额外的整型字段
		value any // 存储的值，使用 any 以接受任何类型
	}
)

// Append 向 SafeSlice 添加新元素
func (ss *SafeSlice) Append(item any) {
	ss.Lock()
	defer ss.Unlock()

	ss.items = append(ss.items, Items{value: item})
}

// Len 返回 SafeSlice 中元素的数量
func (ss *SafeSlice) Len() int {
	ss.RLock()
	defer ss.RUnlock()

	return len(ss.items)
}

// Key 返回特定元素在 SafeSlice 中的索引，不存在则返回 -1
func (ss *SafeSlice) Key(item any) int {
	ss.RLock()
	defer ss.RUnlock()

	for i, v := range ss.items {
		if v.value == item {
			return i
		}
	}

	return -1
}

// Get 返回 SafeSlice 在特定索引上的元素
func (ss *SafeSlice) Get(index int) any {
	ss.RLock()
	defer ss.RUnlock()

	return ss.items[index].value
}

// Update 更新 SafeSlice 在特定索引上的元素
func (ss *SafeSlice) Update(index int, item any) {
	ss.Lock()
	defer ss.Unlock()

	ss.items[index].value = item
}

// List 返回 SafeSlice 中所有元素的副本
func (ss *SafeSlice) List() []any {
	ss.RLock()
	defer ss.RUnlock()

	var r []any
	for _, v := range ss.items {
		r = append(r, v.value)
	}

	return r
}

// Iter 返回一个通道，用于遍历 SafeSlice 中的所有元素
func (ss *SafeSlice) Iter() chan any {
	ss.Lock()

	out := make(chan any)

	go func() {
		defer close(out)
		defer ss.Unlock()

		for _, item := range ss.items {
			out <- item.value
		}
	}()

	return out
}

// Num 返回特定元素的整型字段值
func (ss *SafeSlice) Num(item any) int {
	ss.RLock()
	defer ss.RUnlock()

	for _, v := range ss.items {
		if v.value == item {
			return v.num
		}
	}

	return 0
}

// UpdateNum 更新特定元素的整型字段值
func (ss *SafeSlice) UpdateNum(item any, num int) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num += num
			return
		}
	}
}

// ResetNum 将特定元素的整型字段重置为0
func (ss *SafeSlice) ResetNum(item any) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num = 0
			return
		}
	}
}

// SetNum 设置特定元素的整型字段为给定值
func (ss *SafeSlice) SetNum(item any, num int) {
	ss.Lock()
	defer ss.Unlock()

	for i, v := range ss.items {
		if v.value == item {
			ss.items[i].num = num
			return
		}
	}
}
