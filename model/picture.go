package model

type Picture struct {
	BaseModel BaseModel `gorm:"embedded"`
	Data      []byte    `gorm:"column:data;"`
	Memory_Id uint      `gorm:"column:memory_id;"`
}
