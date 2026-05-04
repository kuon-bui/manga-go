package constant

type CommentReportReason string

const (
	CommentReportReasonSpam         CommentReportReason = "SPAM"
	CommentReportReasonOffensive    CommentReportReason = "OFFENSIVE"
	CommentReportReasonHarassment   CommentReportReason = "HARASSMENT"
	CommentReportReasonAdultContent CommentReportReason = "ADULT_CONTENT"
)

func GetAllowedCommentReportReasons() map[CommentReportReason]bool {
	return map[CommentReportReason]bool{
		CommentReportReasonSpam:         true,
		CommentReportReasonOffensive:    true,
		CommentReportReasonHarassment:   true,
		CommentReportReasonAdultContent: true,
	}
}
