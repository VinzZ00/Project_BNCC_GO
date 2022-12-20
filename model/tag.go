package model

type Tag struct {
	TagId    uint   `gorm:"column:tagId;primaryKey;auto_increment" json:"id,omitempty"`
	TagName  string `gorm:"column:tagName" json:"tag"`
	MemoryId uint   `gorm:"column:memoryId"`
}
