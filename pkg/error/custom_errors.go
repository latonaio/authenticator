package custmerr

type CustomErrMessage struct {
	message string
}

func (c CustomErrMessage) Error() string {
	return c.message
}

var (
	Unknown         = CustomErrMessage{message: "UNKNOWN_ERROR"}
	ErrInternal     = CustomErrMessage{message: "INTERNAL_ERROR"}
	ErrNotFound     = CustomErrMessage{message: "NOT_FOUND_RESOURCE"}
	ErrBadRequest   = CustomErrMessage{message: "BAD_REQUEST"}
	ErrUnauthorized = CustomErrMessage{message: "UNAUTHORIZED"}
	ErrConflict     = CustomErrMessage{message: "CONFLICT"}
)
