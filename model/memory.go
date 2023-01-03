package model

type Memory struct {
	BaseModel
	Description  string
	UserID       uint
	Pictures     []Picture
	MemoriesTags []MemoryTag
}
