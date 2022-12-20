package model

import "time"

type Memory struct {
	Memoryid     uint      `gorm:"column:memoryid;primarykey;auto_increment" json:"memoryid,omitempty"`
	DateAdded    time.Time `gorm:"column:dateAdded;" json:"dateAdded"`
	DateModified time.Time `gorm:"column:dateModified;" json:"dateModified"`
	Desc         string    `gorm:"column:MemoryDesc;" json:"MemoryDesc"`
	UserId       uint      `gorm:"column:userId" json:"UserId"`
	Picture      []Picture `gorm:"foreignkey:MemoryId;references:Memoryid" json:"pictures"`
	Tag          []Tag     `gorm:"foreignKey:MemoryId;referencecs:Memoryid" json:"tags"`
}
