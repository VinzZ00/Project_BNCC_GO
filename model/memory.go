package model

type Memory struct {
	BaseModel BaseModel   `gorm:"embedded"`
	Desc      string      `gorm:"column:MemoryDesc;"`
	Userid    uint        `gorm:"column:UserId"`
	Picture   []Picture   `gorm:"foreignkey:Memory_Id"`
	Tags      []MemoryTag `gorm:"foreignKey:Memory_Id"`
}
