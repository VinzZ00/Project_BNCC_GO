package model

type User struct {
	Userid   uint   `gorm:"column:userId;primaryKey;auto_increment" json:"userId,omitempty"`
	Username string `gorm:"column:userName" json:"userName"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:Password" json:"Password"`
	Memory   Memory `gorm:"foreignKey:UserId;references:Userid"`
}
