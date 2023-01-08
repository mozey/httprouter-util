package share

type Response struct {
	Message string `json:"message"`
}

type ErrResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type JSONRaw string

type ClientVersion struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}
