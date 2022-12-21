package model

import "time"

type Memory struct {
	Memoryid     uint      `gorm:"column:memoryid;primarykey;auto_increment"`
	DateAdded    time.Time `gorm:"column:dateAdded;"`
	DateModified time.Time `gorm:"column:dateModified;"`
	Desc         string    `gorm:"column:MemoryDesc;"`
	UserId       uint      `gorm:"column:userId"`
	Picture      []Picture `gorm:"foreignkey:MemoryId;references:Memoryid"`
	Tag          []Tag     `gorm:"foreignKey:MemoryId;referencecs:Memoryid"`
}
