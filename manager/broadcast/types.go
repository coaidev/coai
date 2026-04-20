package broadcast

type Broadcast struct {
	Index   int    `json:"index"`
	Content string `json:"content"`
}

type Info struct {
	Index     int    `json:"index"`
	Content   string `json:"content"`
	Poster    string `json:"poster"`
	CreatedAt string `json:"created_at"`
}

type listResponse struct {
	Data []Info `json:"data"`
}

type createRequest struct {
	Content string `json:"content"`
}

type updateRequest struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type createResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
