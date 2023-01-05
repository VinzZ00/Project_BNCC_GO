package model

type Tag struct {
	BaseModel
	Name         string      `json:"name"`
	MemoriesTags []MemoryTag `json:"memories_tags"`
}
