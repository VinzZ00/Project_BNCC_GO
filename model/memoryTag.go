package model

type MemoryTag struct {
	MemoryID uint `gorm:"primaryKey" json:"memory_id"`
	TagID    uint `gorm:"primaryKey" json:"tag_id"`
	Tag      Tag  `gorm:"-:migration"`
}
