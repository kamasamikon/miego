package xinit

import (
	"container/list"
	"fmt"
	"sync"
)

type CBInfo struct {
	bitmapIndex int // 位图位置 0-255
	cb          func() bool
}

var (
	cbList *list.List
	mu     sync.Mutex

	// 位图相关 - 4个uint64支持256个CB
	bitmap    [4]uint64          // bitmap[0]:0-63, bitmap[1]:64-127, bitmap[2]:128-191, bitmap[3]:192-255
	stateMap  map[[4]uint64]bool // 记录出现过的状态
	nextIndex int                // 下一个可用的位图索引
)

func init() {
	cbList = list.New()
	stateMap = make(map[[4]uint64]bool)
	nextIndex = 0
}

func Add(cb func() bool) {
	if cb == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if nextIndex >= 256 {
		panic("最多支持256个并发CB")
	}

	// 分配位图索引
	index := nextIndex
	nextIndex++

	// 在位图中标记
	setBitmapBit(index)

	// 添加到队列
	cbInfo := &CBInfo{
		bitmapIndex: index,
		cb:          cb,
	}
	cbList.PushBack(cbInfo)
}

func Done() {
	mu.Lock()
	defer mu.Unlock()

	for cbList.Len() > 0 {
		// 在处理前检查死循环
		currentState := bitmap
		if stateMap[currentState] {
			panic(fmt.Sprintf("检测到死循环！状态 %v 重复出现", currentState))
		}
		stateMap[currentState] = true

		elem := cbList.Front()
		if elem == nil {
			break
		}

		cbInfo := elem.Value.(*CBInfo)
		cbList.Remove(elem)

		// 执行回调
		ok := cbInfo.cb()

		if !ok {
			// 失败，重新入队（位图保持不变）
			if cbList.Len() > 0 {
				cbList.PushBack(cbInfo)
			}
		} else {
			// 成功，清除位图标记
			clearBitmapBit(cbInfo.bitmapIndex)
		}
	}
}

func setBitmapBit(index int) {
	slot := index / 64
	offset := uint(index % 64)
	bitmap[slot] |= 1 << offset
}

func clearBitmapBit(index int) {
	slot := index / 64
	offset := uint(index % 64)
	bitmap[slot] &^= 1 << offset
}
