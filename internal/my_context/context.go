package mycontext

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	TrackKey  contextKey = "track"
)
