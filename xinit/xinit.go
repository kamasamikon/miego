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
	// bitmap[0]:0-63, bitmap[1]:64-127, bitmap[2]:128-191, bitmap[3]:192-255
	bitmap    [4]uint64
	nextIndex int // 下一个可用的位图索引
)

func init() {
	cbList = list.New()
	nextIndex = 0
}

func Add(cb func() bool) {
	if cb == nil {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if nextIndex >= 256 {
		panic("最多支持256个CB")
	}

	index := nextIndex
	nextIndex++

	setBitmapBit(index)

	cbInfo := &CBInfo{
		bitmapIndex: index,
		cb:          cb,
	}
	cbList.PushBack(cbInfo)
}

func Done() error {
	mu.Lock()
	defer mu.Unlock()

	if cbList.Len() == 0 {
		return nil
	}

	// 一轮完整扫描的长度：连续这么多次没有成功，即认为无法推进
	maxNoProgress := cbList.Len()
	noProgress := 0

	for cbList.Len() > 0 {
		elem := cbList.Front()
		if elem == nil {
			break
		}

		cbInfo := elem.Value.(*CBInfo)
		cbList.Remove(elem)

		ok := cbInfo.cb()

		if ok {
			clearBitmapBit(cbInfo.bitmapIndex)
			noProgress = 0
		} else {
			// 失败：无条件重新入队
			cbList.PushBack(cbInfo)
			noProgress++
			if noProgress >= maxNoProgress {
				return fmt.Errorf("初始化无法完成：连续 %d 次没有回调成功，剩余位图=%v", noProgress, bitmap,)
			}
		}
	}

	return nil
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
