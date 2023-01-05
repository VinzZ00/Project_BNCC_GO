package model

type Picture struct {
	BaseModel
	Data     string `json:"data"`
	MemoryID uint   `json:"memory_id"`
}
