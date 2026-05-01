package dify

const (
	ChatMessagesUrl            = "/v1/chat-messages"
	MessagesUrl                = "/v1/messages"
	MessagesFeedbacksUrl       = "/v1/messages/{message_id}/feedbacks"
	MessagesSuggestedUrl       = "/v1/messages/{message_id}/suggested"
	ConversationsUrl           = "/v1/conversations"
	ConversationsAddMessageUrl = "/v1/conversations/{conversation_id}/message"
	ConversationsRenameUrl     = "/v1/conversations/{conversation_id}/name"
	ParametersUrl              = "/v1/parameters"
	WorkflowsRunUrl            = "/v1/workflows/run"
	TextToAudioUrl             = "/v1/text-to-audio"

	FeedbackLike    = "like"
	FeedbackDislike = "dislike"
	FeedBackCancel  = "null"

	ResponseModeStreaming = "streaming"
	ResponseModeBlocking  = "blocking"

	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
	StatusStopped   = "stopped"

	ChatFileTypeImage            = "image"      // 图片类型
	ChatFileTypeDocument         = "document"   // 文档类型
	ChatFileTransferMethodRemote = "remote_url" // 远程文件, 使用url
	ChatFileTransferMethodLocal  = "local_file" // 本地文件，使用dify的文件ID
)

// 事件类型常量
const (
	EventWorkflowStarted  = "workflow_started"
	EventNodeStarted      = "node_started"
	EventNodeFinished     = "node_finished"
	EventWorkflowFinished = "workflow_finished"
	EventTTSMessage       = "tts_message"
	EventTTSMessageEnd    = "tts_message_end"
	EventTextChunk        = "text_chunk"
)
