package model

type Picture struct {
	BaseModel
	Data      []byte    `gorm:"column:data;"`
	Memory_Id uint      `gorm:"column:memory_id;"`
}
