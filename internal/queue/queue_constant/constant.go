package queueconstant

// queue
const (
	MAIL_DELIVER_QUEUE = "mail_queue"
	NOTIFICATION_QUEUE = "notification_queue"
)

// task
const (
	MAIL_DELIVER_TASK        = "mail"
	MULTI_MAIL_DELIVER_TASK  = "multi_mail"
	NOTIFICATION_FANOUT_TASK = "notification_fanout"
)
