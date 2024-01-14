package infra

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	AccountUUIDKey Key = "AccountUUID"
	TokenKey       Key = "user-token"
	SessionKey     Key = "Session"
)
