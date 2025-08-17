package clients

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int            `json:"update_id"`
	Message *IncomeMessage `json:"message"`
}

type IncomeMessage struct {
	Text string   `json:"text"`
	From From     `json:"from"`
	Chat Chat     `json:"chat"`
	File Document `json:"file"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type OutputMessage struct {
	Text string    `json:"text"`
	Chat Chat      `json:"chat"`
	File *Document `json:"file"`
}
