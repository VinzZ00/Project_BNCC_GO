package model

type Memory struct {
	BaseModel
	Description  string      `json:"description"`
	UserID       uint        `json:"user_id"`
	Pictures     []Picture   `json:"pictures"`
	MemoriesTags []MemoryTag `json:"memories_tags"`
}
