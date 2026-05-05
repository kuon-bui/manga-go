package constant

import "strings"

type CommentReportReason string

const (
	CommentReportReasonSpam         CommentReportReason = "SPAM"
	CommentReportReasonOffensive    CommentReportReason = "OFFENSIVE"
	CommentReportReasonHarassment   CommentReportReason = "HARASSMENT"
	CommentReportReasonAdultContent CommentReportReason = "ADULT_CONTENT"
)

func GetAllCommentReportReasons() []CommentReportReason {
	return []CommentReportReason{
		CommentReportReasonSpam,
		CommentReportReasonOffensive,
		CommentReportReasonHarassment,
		CommentReportReasonAdultContent,
	}
}

func GetAllCommentReportReasonsStr() string {
	reasons := GetAllCommentReportReasons()
	var res strings.Builder
	res.WriteString(string(reasons[0]))
	for _, r := range reasons[1:] {
		res.WriteString(", " + string(r))
	}
	return res.String()
}

func GetAllowedCommentReportReasons() map[CommentReportReason]bool {
	return map[CommentReportReason]bool{
		CommentReportReasonSpam:         true,
		CommentReportReasonOffensive:    true,
		CommentReportReasonHarassment:   true,
		CommentReportReasonAdultContent: true,
	}
}
