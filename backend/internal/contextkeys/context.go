package contextkeys

type key string

const (
	UserKey         key = "user"
	SessionTokenKey key = "session_token"
)

func (k key) String() string {
	return string(k)
}
