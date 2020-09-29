package roast

const (
	// 没有被删除
	RemWhy_NotRem = 0

	// 同一个UUID的添加了新记录，老记录自动标记为删除
	RemWhy_Update = 1

	// 显式删除
	RemWhy_Delete = 2
)

type TableHeader struct {
	// 记录的序列号而已
	ID uint `gorm:"Column:ID;primary_key"`

	// NewAt: 记录的添加时间
	NewAt uint64 `gorm:"Column:NewAt"`
	NewBy string `gorm:"Column:NewBy"`

	// See RemWhy_Delete etc.
	RemAt  uint64 `gorm:"Column:RemAt"`
	RemBy  string `gorm:"Column:RemBy"`
	RemWhy int    `gorm:"Column:RemWhy"`
}