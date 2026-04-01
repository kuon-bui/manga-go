package readinghistoryrequest

import (
	"time"
)

type UpdateReadingHistoryRequest struct {
	LastReadAt *time.Time `json:"lastReadAt"`
}
