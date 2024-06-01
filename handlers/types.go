package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Comment string `json:"comment,omitempty"`
}
