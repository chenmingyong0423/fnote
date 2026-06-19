package domain

type FileUsageEvent struct {
	EntityId       string   `json:"entity_id"`
	AddedFileIds   []string `json:"added_file_ids,omitempty"`
	DeletedFileIds []string `json:"deleted_file_ids,omitempty"`
}
