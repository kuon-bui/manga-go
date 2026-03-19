package common

type Requester interface {
	GetUserId() uint
	GetEmail() string
}

type GinCtxKey string

const (
	CurrentUser GinCtxKey = "current_user"
	TokenId     GinCtxKey = "token_id"
)
