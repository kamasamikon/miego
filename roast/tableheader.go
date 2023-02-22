package roast

import (
	"strconv"
	"time"
)

const (
	// 没有被删除
	RemWhy_NotRem = 0

	// 同一个UUID的添加了新记录，老记录自动标记为删除
	RemWhy_Update = 1

	// 显式删除
	RemWhy_Delete = 2
)

// NewAt, RemAt, CrtAt: NNNNYYRRSSFFMM
type TableHeader struct {
	// 记录的序列号而已
	ID uint `gorm:"Column:ID;primary_key"`

	// NewAt: 记录的添加时间
	NewAt uint64 `gorm:"Column:NewAt"`
	NewBy string `gorm:"Column:NewBy;size:40"`

	// See RemWhy_Delete etc.
	RemAt  uint64 `gorm:"Column:RemAt"`
	RemBy  string `gorm:"Column:RemBy;size:40"`
	RemWhy uint8  `gorm:"Column:RemWhy"`

	// CrtAt: UUID对应的项目的日期
	// UUID: 真正的记录的ID
	CrtAt uint64 `gorm:"Column:CrtAt"`
	UUID  string `gorm:"Column:UUID;size:40"`
}

func Setup(h *TableHeader, NewBy string) {
	now, _ := strconv.ParseUint(time.Now().Format("20060102150405"), 0, 64)

	h.NewAt = now
	h.NewBy = NewBy
	if h.CrtAt == 0 {
		h.CrtAt = now
	}
}
