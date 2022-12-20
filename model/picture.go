package model

type Picture struct {
	PictureId uint   `gorm:"column:pictureid;PrimaryKey;auto_increment" json:"pictureId,omitempty"`
	Picture   []byte `gorm:"column:picture;" json:"picture"`
	MemoryId  uint   `gorm:"column:memoryid;"`
}
