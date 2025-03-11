package apperrors

type ErrorCode string

const (
    ErrNotFound       ErrorCode = "NOT_FOUND"
    ErrInvalidInput   ErrorCode = "INVALID_INPUT"
    ErrInternal       ErrorCode = "INTERNAL_ERROR"
    ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
    ErrAlreadyExists  ErrorCode = "ALREADY_EXISTS"
)

type Error struct {
    Code       ErrorCode   `json:"code"`
    Message    string     `json:"message"`
    Details    any        `json:"details,omitempty"`
}

func (e *Error) Error() string {
    return e.Message
}

func NewError(code ErrorCode, message string, details ...any) *Error {
    var detail any
    if len(details) > 0 {
        detail = details[0]
    }
    
    return &Error{
        Code:    code,
        Message: message,
        Details: detail,
    }
}