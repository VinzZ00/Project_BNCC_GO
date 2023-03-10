package model

type User struct {
	BaseModel
	Username string `json:"username"`
	// by default, `string` uses `longtext` type but MySQL doesn't support "foreign key" or "index"
	// on that data type. source: https://stackoverflow.com/a/69201429
	Email    string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string   `json:"-"`
	Memories []Memory `json:"memories"`
}
