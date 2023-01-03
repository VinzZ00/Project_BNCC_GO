package model

type Picture struct {
	BaseModel
	Data     []byte
	MemoryID uint
}
