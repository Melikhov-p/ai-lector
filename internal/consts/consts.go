package consts

type ContextKey string

type URLParam string

const (
	UserIDContextKey ContextKey = "UserID"
)

const (
	UserIDURLParam     URLParam = "userID"
	InterestIDURLParam URLParam = "interestID"
)

func (r URLParam) String() string {
	return string(r)
}
