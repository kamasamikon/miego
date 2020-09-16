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
	ID uint `gorm:"Column:ID;primary_key"`

	// NewAt: 这个时间应该是什么？UTC？国外访问？
	NewAt uint64 `gorm:"Column:NewAt"`
	NewBy string `gorm:"Column:NewBy"`

	// See RemWhy_Delete etc.
	RemAt  uint64 `gorm:"Column:RemAt"`
	RemBy  string `gorm:"Column:RemBy"`
	RemWhy int    `gorm:"Column:RemWhy"`
}
