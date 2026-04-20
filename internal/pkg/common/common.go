package common

import (
	"fmt"
	"strings"
)

type Requester interface {
	GetUserId() uint
	GetEmail() string
}

type GinCtxKey string

const (
	CurrentUser GinCtxKey = "current_user"
	TokenId     GinCtxKey = "token_id"
)

func ShowDebugTrace(module string, trace []byte) {
	module = fmt.Sprintf("DEBUG TRACE: %s ", module)
	generateTitle := func(title string) string {
		maxLength := 150
		length := (maxLength - 2 - len(title)) / 2

		return fmt.Sprintf("%s %s %s", strings.Repeat("=", length), title, strings.Repeat("=", length))
	}

	fmt.Printf("\n%s%s\n%s\n%s%s\n\n",
		"\033[31m", generateTitle(module),
		trace,
		generateTitle("END "+module), "\033[0m",
	)
}

func AddFileContentPrefix(s string) string {
	if s == "" || strings.HasPrefix(s, "/") || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}
	v := "/files/content/" + s
	return v
}
