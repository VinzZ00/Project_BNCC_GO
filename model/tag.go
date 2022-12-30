package model

type Tag struct {
	BaseModel BaseModel `gorm:"embedded"`
	Name      string
	MemoryTag []MemoryTag `gorm:"foreignKey:Tag_Id"`
}
