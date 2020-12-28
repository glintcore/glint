package api

type AccountPasswordRequest struct {
	Password string `json:"password"`
}

type PostRequest struct {
	Data string `json:"data"`
}

type PostResponse struct {
	Url string `json:"url"`
}

type MetadataRequest struct {
	Metadata string `json:"metadata"`
}

type LoginRequest struct {
}

type LoginResponse struct {
	SessionId string `json:"sessionId"`
}
