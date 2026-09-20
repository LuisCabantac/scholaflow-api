package response

type APIResponse[T any] struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
	Data       T      `json:"data,omitempty"`
}
