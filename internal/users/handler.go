package users

// Handler stores methods for handling users http-requests.
type Handler struct{}

// NewHandler creates and returns a new users [Handler] instance.
func NewHandler() *Handler {
	return &Handler{}
}
