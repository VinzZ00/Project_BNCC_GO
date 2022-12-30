package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetDB() (*gorm.DB, error) {
	dsn := "root@tcp(127.0.0.1:3306)/bncc_go?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	}})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Memory{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Picture{})
	db.AutoMigrate(&MemoryTag{})
	return db, err
}
