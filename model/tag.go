package model

type Tag struct {
	BaseModel
	Name      string
	MemoryTag []MemoryTag `gorm:"foreignKey:Tag_Id"`
}
