package model

type Tag struct {
	BaseModel
	Name         string
	MemoriesTags []MemoryTag
}
