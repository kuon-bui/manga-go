package queueconstant

// queue
const (
	MAIL_DELIVER_QUEUE  = "mail_queue"
	NOTIFICATION_QUEUE  = "notification_queue"
	IMAGE_PROCESS_QUEUE = "image_process_queue"
)

// task
const (
	MAIL_DELIVER_TASK          = "mail"
	MULTI_MAIL_DELIVER_TASK    = "multi_mail"
	NOTIFICATION_FANOUT_TASK   = "notification_fanout"
	IMAGE_PROCESS_TASK         = "image_process"
	IMAGE_PROCESS_CLEANUP_TASK = "image_process_cleanup"
)
