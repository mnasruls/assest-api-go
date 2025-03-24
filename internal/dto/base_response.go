package dto

type BaseResponse struct {
	Data             interface{} `json:"data,omitempty"`
	Message          string      `json:"message,omitempty"`
	Error            string      `json:"error,omitempty"`
	ErrorDescription string      `json:"error_description,omitempty"`
}
