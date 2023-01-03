package model

type User struct {
	BaseModel
	Username string
	Email     string    `gorm:"column:email;uniqueIndex;"`
	Password string
	Memories []Memory
}
