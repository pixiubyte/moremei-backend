package dify

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type ChatFile struct {
	Type           string `json:"type"`                     // 文件类型 image | document
	TransferMethod string `json:"transfer_method"`          // 传递方式，remote_url 图片地址 / local_file 上传文件
	Url            string `json:"url,omitempty"`            // 图片地址（仅当传递方式为 remote_url 时）
	UploadFileId   string `json:"upload_file_id,omitempty"` // 上传文件 ID（仅当传递方式为 local_file 时）
}

type ChatMessageRequest struct {
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query"`
	ResponseMode   string                 `json:"response_mode"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	User           string                 `json:"user"`
	Files          []*ChatFile            `json:"files,omitempty"`
}

type ChatMessageResponse struct {
	ID             string `json:"id"`
	Answer         string `json:"answer"`
	ConversationID string `json:"conversation_id"`
	CreatedAt      int    `json:"created_at"`
}

type ChatMessageStreamResponse struct {
	Event          string `json:"event"`
	TaskID         string `json:"task_id"`
	ID             string `json:"id"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
	ConversationID string `json:"conversation_id"`
}

type ChatMessageStreamChannelResponse struct {
	ChatMessageStreamResponse
	Err error `json:"-"`
}

type MessagesFeedbacksRequest struct {
	MessageID string `json:"message_id,omitempty"`
	Rating    string `json:"rating,omitempty"`
	User      string `json:"user"`
}

type MessagesFeedbacksResponse struct {
	HasMore bool                            `json:"has_more"`
	Data    []MessagesFeedbacksDataResponse `json:"data"`
}

type MessagesFeedbacksDataResponse struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	PhoneNumber    string `json:"phone_number"`
	AvatarURL      string `json:"avatar_url"`
	DisplayName    string `json:"display_name"`
	ConversationID string `json:"conversation_id"`
	LastActiveAt   int64  `json:"last_active_at"`
	CreatedAt      int64  `json:"created_at"`
}

type MessagesRequest struct {
	ConversationID string `json:"conversation_id"`
	FirstID        string `json:"first_id,omitempty"`
	Limit          int    `json:"limit"`
	User           string `json:"user"`
}

type MessagesResponse struct {
	Limit   int                    `json:"limit"`
	HasMore bool                   `json:"has_more"`
	Data    []MessagesDataResponse `json:"data"`
}

type MessagesDataResponse struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Inputs         map[string]interface{} `json:"inputs"`
	Query          string                 `json:"query"`
	Answer         string                 `json:"answer"`
	Feedback       interface{}            `json:"feedback"`
	CreatedAt      int64                  `json:"created_at"`
}

type ConversationsRequest struct {
	LastID string `json:"last_id,omitempty"`
	Limit  int    `json:"limit"`
	User   string `json:"user"`
}

type ConversationsResponse struct {
	Limit   int                         `json:"limit"`
	HasMore bool                        `json:"has_more"`
	Data    []ConversationsDataResponse `json:"data"`
}

type ConversationsDataResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Inputs    map[string]string `json:"inputs"`
	Status    string            `json:"status"`
	CreatedAt int64             `json:"created_at"`
}

type ConversationsRenamingRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Name           string `json:"name"`
	User           string `json:"user"`
}

type ConversationsRenamingResponse struct {
	Result string `json:"result"`
}

type ParametersRequest struct {
	User string `json:"user"`
}

type ParametersResponse struct {
	OpeningStatement              string        `json:"opening_statement"`
	SuggestedQuestions            []interface{} `json:"suggested_questions"`
	SuggestedQuestionsAfterAnswer struct {
		Enabled bool `json:"enabled"`
	} `json:"suggested_questions_after_answer"`
	MoreLikeThis struct {
		Enabled bool `json:"enabled"`
	} `json:"more_like_this"`
	UserInputForm []map[string]interface{} `json:"user_input_form"`
}

type ConversationsAddMessageRequest struct {
	ConversationId string
	User           string `json:"user"`
	Answer         string `json:"answer"`
}

type ConversationsAddMessageResponse struct {
	Result string `json:"result"`
}

type WorkflowRunRequest struct {
	Inputs       map[string]interface{} `json:"inputs"`
	ResponseMode string                 `json:"response_mode"`
	User         string                 `json:"user"`
}

type WorkflowRunData struct {
	Id          string      `json:"id"`
	Text        string      `json:"text"`
	WorkflowId  string      `json:"workflow_id"`
	Status      string      `json:"status"`
	Outputs     interface{} `json:"outputs"`
	Error       string      `json:"error,omitempty"`
	ElapsedTime float64     `json:"elapsed_time,omitempty"`
	TotalTokens int         `json:"total_tokens"`
	TotalSteps  int         `json:"total_steps"`
	CreatedAt   int64       `json:"created_at"`
	FinishedAt  int64       `json:"finished_at"`
}

type WorkflowRunResponse struct {
	TaskID        string          `json:"task_id"`
	WorkflowRunId string          `json:"workflow_run_id"`
	Data          WorkflowRunData `json:"data"`
}

type WorkflowRunStreamResponse struct {
	Event string `json:"event"`
	WorkflowRunResponse
}

type WorkflowRunStreamChannelResponse struct {
	WorkflowRunStreamResponse
	Err error `json:"-"`
}

type MessagesSuggestedRequest struct {
	MessageID string `json:"message_id"`
	User      string `json:"user"`
}

type MessagesSuggestedResponse struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

func (c *Client) ChatMessages(ctx context.Context, req ChatMessageRequest) (resp *ChatMessageResponse, err error) {
	req.ResponseMode = ResponseModeBlocking
	if req.Inputs == nil {
		req.Inputs = map[string]interface{}{}
	}

	var r *http.Request
	r, err = c.createPostRequest(ctx, c.buildRequestApi(ChatMessagesUrl), req)
	if err != nil {
		return
	}

	var _resp ChatMessageResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)

	if !c.thinking {
		_resp.Answer = c.trimThinkingContent(_resp.Answer)
	}

	resp = &_resp
	return
}

func (c *Client) ChatMessagesStreamRaw(ctx context.Context, req ChatMessageRequest) (resp *http.Response, err error) {
	req.ResponseMode = ResponseModeStreaming
	if req.Inputs == nil {
		req.Inputs = map[string]interface{}{}
	}
	var r *http.Request
	r, err = c.createPostRequest(ctx, c.buildRequestApi(ChatMessagesUrl), req)
	if err != nil {
		return
	}

	return c.SetHttpRequest(r).SendRequestStream()
}

func (c *Client) ChatMessagesStream(ctx context.Context, req ChatMessageRequest) (streamChannel chan ChatMessageStreamChannelResponse, err error) {
	var resp *http.Response
	if resp, err = c.ChatMessagesStreamRaw(ctx, req); err != nil {
		return
	}

	streamChannel = make(chan ChatMessageStreamChannelResponse)
	go c.chatMessagesStreamHandle(ctx, resp, streamChannel)
	return
}

func (c *Client) chatMessagesStreamHandle(ctx context.Context, resp *http.Response, streamChannel chan ChatMessageStreamChannelResponse) {
	var (
		body   = resp.Body
		reader = bufio.NewReader(body)

		err  error
		line []byte
	)

	defer resp.Body.Close()
	defer close(streamChannel)

	thinkingFlag := false
	thinkingStarted := false
	thinkingStopped := false
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if line, err = reader.ReadBytes('\n'); err != nil {
				streamChannel <- ChatMessageStreamChannelResponse{
					Err: errors.New("Error reading line: " + err.Error()),
				}
				return
			}

			if !bytes.HasPrefix(line, []byte("data:")) {
				continue
			}

			line = bytes.TrimPrefix(line, []byte("data:"))
			line = bytes.TrimSpace(line)

			var resp ChatMessageStreamChannelResponse
			if err = json.Unmarshal(line, &resp); err != nil {
				streamChannel <- ChatMessageStreamChannelResponse{
					Err: errors.New("Error unmarshalling event: " + err.Error()),
				}
				return
			}

			if resp.Event == "message_end" || resp.Event == "error" {
				return
			}

			if resp.Answer == "" {
				continue
			}

			if resp.Event == "message" {
				if !c.thinking {
					if !thinkingStarted {
						thinkingFlag = regexp.MustCompile(`<details.*?>`).MatchString(resp.Answer) || regexp.MustCompile(`<think.*?>`).MatchString(resp.Answer)
						thinkingStarted = true
					}
					if thinkingStarted && thinkingFlag && !thinkingStopped {
						thinkingStopped = strings.Contains(resp.Answer, "</details>") || strings.Contains(resp.Answer, "</think>")
						continue
					}
				}
			}

			streamChannel <- resp
		}
	}
}

func (c *Client) Messages(ctx context.Context, req MessagesRequest) (resp *MessagesResponse, err error) {
	var u = url.Values{}
	u.Set("conversation_id", req.ConversationID)
	u.Set("user", req.User)

	if req.FirstID != "" {
		u.Set("first_id", req.FirstID)
	}
	if req.Limit > 0 {
		var l = int64(req.Limit)
		u.Set("limit", strconv.FormatInt(l, 10))
	}

	var r *http.Request
	if r, err = c.createGetRequest(ctx, c.buildRequestApi(MessagesUrl), u); err != nil {
		return
	}

	var _resp MessagesResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) MessagesFeedbacks(ctx context.Context, req MessagesFeedbacksRequest) (resp *MessagesFeedbacksResponse, err error) {
	if req.MessageID == "" {
		err = errors.New("MessagesFeedbacksRequest.MessageID Illegal")
		return
	}

	u := c.buildRequestApi(MessagesFeedbacksUrl)
	u = strings.ReplaceAll(u, "{message_id}", req.MessageID)

	req.MessageID = ""

	var r *http.Request
	if r, err = c.createPostRequest(ctx, u, req); err != nil {
		return
	}

	var _resp MessagesFeedbacksResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

// MessagesSuggested 消息建议, get方法
func (c *Client) MessagesSuggested(ctx context.Context, req MessagesSuggestedRequest) (resp *MessagesSuggestedResponse, err error) {
	if req.MessageID == "" {
		err = errors.New("MessagesSuggestedRequest.MessageID Illegal")
		return
	}

	u := c.buildRequestApi(MessagesSuggestedUrl)
	u = strings.ReplaceAll(u, "{message_id}", req.MessageID)

	var uv = url.Values{}
	uv.Set("user", req.User)

	var r *http.Request
	if r, err = c.createGetRequest(ctx, u, uv); err != nil {
		return
	}

	var _resp MessagesSuggestedResponse
	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) Conversations(ctx context.Context, req ConversationsRequest) (resp *ConversationsResponse, err error) {
	if req.User == "" {
		err = errors.New("ConversationsRequest.User Illegal")
		return
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	var u = url.Values{}
	u.Set("last_id", req.LastID)
	u.Set("user", req.User)

	var l = int64(req.Limit)
	u.Set("limit", strconv.FormatInt(l, 10))

	var r *http.Request
	if r, err = c.createGetRequest(ctx, c.buildRequestApi(ConversationsUrl), u); err != nil {
		return
	}

	var _resp ConversationsResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) ConversationsRenaming(ctx context.Context, req ConversationsRenamingRequest) (resp *ConversationsRenamingResponse, err error) {
	u := c.buildRequestApi(ConversationsRenameUrl)
	u = strings.ReplaceAll(u, "{conversation_id}", req.ConversationID)

	req.ConversationID = ""

	var r *http.Request
	if r, err = c.createPostRequest(ctx, u, req); err != nil {
		return
	}

	var _resp ConversationsRenamingResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) ConversationsAddMessage(ctx context.Context, req ConversationsAddMessageRequest) (resp *ConversationsAddMessageResponse, err error) {
	u := c.buildRequestApi(ConversationsAddMessageUrl)
	u = strings.ReplaceAll(u, "{conversation_id}", req.ConversationId)

	var r *http.Request
	if r, err = c.createPostRequest(ctx, u, req); err != nil {
		return
	}

	var _resp ConversationsAddMessageResponse
	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) Parameters(ctx context.Context, req ParametersRequest) (resp *ParametersResponse, err error) {
	if req.User == "" {
		err = errors.New("ParametersRequest.User Illegal")
		return
	}

	var u = url.Values{}
	u.Set("user", req.User)

	var r *http.Request
	if r, err = c.createGetRequest(ctx, c.buildRequestApi(ParametersUrl), u); err != nil {
		return
	}

	var _resp ParametersResponse

	err = c.SetHttpRequest(r).SendRequest(&_resp)
	resp = &_resp
	return
}

func (c *Client) WorkflowsRun(ctx context.Context, req WorkflowRunRequest) (resp *WorkflowRunResponse, err error) {
	req.ResponseMode = ResponseModeBlocking
	var r *http.Request
	r, err = c.createPostRequest(ctx, c.buildRequestApi(WorkflowsRunUrl), req)
	if err != nil {
		return
	}

	var _resp WorkflowRunResponse
	err = c.SetHttpRequest(r).SendRequest(&_resp)
	if err != nil {
		return
	}
	resp = &_resp
	return
}

func (c *Client) WorkflowsRunStream(ctx context.Context, req WorkflowRunRequest) (streamChannel chan WorkflowRunStreamChannelResponse, err error) {
	req.ResponseMode = ResponseModeStreaming
	streamChannel = make(chan WorkflowRunStreamChannelResponse)
	var r *http.Request
	r, err = c.createPostRequest(ctx, c.buildRequestApi(WorkflowsRunUrl), req)
	if err != nil {
		return
	}

	resp, err := c.SetHttpRequest(r).SendRequestStream()
	if err != nil {
		return nil, err
	}

	go c.workflowRunStreamHandle(ctx, resp, streamChannel)
	return
}

func (c *Client) workflowRunStreamHandle(ctx context.Context, resp *http.Response, streamChannel chan WorkflowRunStreamChannelResponse) {
	var (
		body   = resp.Body
		reader = bufio.NewReader(body)

		err  error
		line []byte
	)

	defer resp.Body.Close()
	defer close(streamChannel)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if line, err = reader.ReadBytes('\n'); err != nil {
				streamChannel <- WorkflowRunStreamChannelResponse{
					Err: errors.New("Error reading line: " + err.Error()),
				}
				return
			}

			if !bytes.HasPrefix(line, []byte("data:")) {
				continue
			}

			line = bytes.TrimPrefix(line, []byte("data:"))
			line = bytes.TrimSpace(line)

			var resp WorkflowRunStreamChannelResponse
			if err = json.Unmarshal(line, &resp); err != nil {
				streamChannel <- WorkflowRunStreamChannelResponse{
					Err: errors.New("Error unmarshalling event: " + err.Error()),
				}
				return
			}
			streamChannel <- resp
		}
	}
}
