package model

type MemoryTag struct {
	MemoryID uint `gorm:"primaryKey"`
	TagID    uint `gorm:"primaryKey"`
}
