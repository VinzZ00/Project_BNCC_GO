package model

type User struct {
	BaseModel
	Username  string    `gorm:"column:userName"`
	Email     string    `gorm:"column:email;uniqueIndex;"`
	Password  string    `gorm:"column:Password"`
	Memory    Memory    `gorm:"foreignKey:Userid;"`
}
