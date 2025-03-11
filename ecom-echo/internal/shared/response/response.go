// internal/shared/response/response.go
package response

type Response struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(message string, data interface{}) *Response {
    return &Response{
        Status:  true,
        Message: message,
        Data:    data,
    }
}

func Error(message string, err interface{}) *Response {
    return &Response{
        Status:  false,
        Message: message,
        Error:   err,
    }
}