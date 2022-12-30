package model

import "time"

type BaseModel struct {
	Id         uint `gorm:"primarykey;auto_increment"`
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}
