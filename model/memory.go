package model

type Memory struct {
	BaseModel
	Desc      string      `gorm:"column:MemoryDesc;"`
	Userid    uint        `gorm:"column:UserId"`
	Picture   []Picture   `gorm:"foreignkey:Memory_Id"`
	Tags      []MemoryTag `gorm:"foreignKey:Memory_Id"`
}
