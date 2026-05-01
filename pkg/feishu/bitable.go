package feishu

import (
	"encoding/json"
	"fmt"
	"moremei/ai-saas/pkg/httpc"
	"net/http"
	"strings"
)

// BatchCreateRecordsRequest represents the request body for batch creating records
type BatchCreateRecordsRequest struct {
	Records []Record `json:"records"`
}

// Record represents a single record in the batch create request
type Record struct {
	Fields map[string]interface{} `json:"fields"`
}

// BatchCreateRecordsResponse represents the response from batch creating records
type BatchCreateRecordsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Records []Record `json:"records"`
	} `json:"data"`
}

// CreateBitableRequest represents the request body for creating a Bitable
type CreateBitableRequest struct {
	Name        string `json:"name,omitempty"`
	FolderToken string `json:"folder_token,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
}

// CreateBitableResponse represents the response from creating a Bitable
type CreateBitableResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		App struct {
			AppToken       string `json:"app_token"`
			DefaultTableID string `json:"default_table_id"`
			FolderToken    string `json:"folder_token"`
			Name           string `json:"name"`
			URL            string `json:"url"`
		} `json:"app"`
	} `json:"data"`
}

// Field represents a field in a Bitable table
type Field struct {
	FieldID     string      `json:"field_id"`
	FieldName   string      `json:"field_name"`
	Type        int         `json:"type"`
	IsPrimary   bool        `json:"is_primary"`
	Property    interface{} `json:"property"`
	Description string      `json:"description"`
}

// ListFieldsResponse represents the response from listing fields
type ListFieldsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		HasMore   bool    `json:"has_more"`
		Items     []Field `json:"items"`
		PageToken string  `json:"page_token,omitempty"`
	} `json:"data"`
}

// UpdateFieldRequest represents the request body for updating a field
type UpdateFieldRequest struct {
	FieldName   string      `json:"field_name"`
	Type        int         `json:"type"`
	Property    interface{} `json:"property,omitempty"`
	Description string      `json:"description,omitempty"`
	UIType      string      `json:"ui_type,omitempty"`
}

// DeleteFieldResponse represents the response body for deleting a field
type DeleteFieldResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		FieldID string `json:"field_id"`
		Deleted bool   `json:"deleted"`
	} `json:"data"`
}

// UpdateFieldResponse represents the response body for updating a field
type UpdateFieldResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Field Field `json:"field"`
	} `json:"data"`
}

// CreateFieldRequest represents the request body for creating a new field
type CreateFieldRequest struct {
	FieldName   string      `json:"field_name"`
	Type        int         `json:"type"`
	Property    interface{} `json:"property,omitempty"`
	Description string      `json:"description,omitempty"`
	UIType      string      `json:"ui_type,omitempty"`
}

// CreateFieldResponse represents the response body for creating a new field
type CreateFieldResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Field Field `json:"field"`
	} `json:"data"`
}

// BatchGetIdRequest represents the request body for batch getting user IDs
type BatchGetIdRequest struct {
	Emails          []string `json:"emails,omitempty"`
	Mobiles         []string `json:"mobiles,omitempty"`
	IncludeResigned bool     `json:"include_resigned,omitempty"`
}

// BatchGetIdResponse represents the response for batch getting user IDs
type BatchGetIdResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserList []struct {
			Mobile  string `json:"mobile,omitempty"`
			Email   string `json:"email,omitempty"`
			UserId  string `json:"user_id,omitempty"`
			OpenId  string `json:"open_id,omitempty"`
			UnionId string `json:"union_id,omitempty"`
			Status  struct {
				IsActivated bool `json:"is_activated,omitempty"`
				IsExited    bool `json:"is_exited,omitempty"`
				IsFrozen    bool `json:"is_frozen,omitempty"`
				IsResigned  bool `json:"is_resigned,omitempty"`
				IsUnjoin    bool `json:"is_unjoin,omitempty"`
			} `json:"status"`
		} `json:"user_list"`
	} `json:"data"`
}

type BatchAddPermissionMemberRequest struct {
	Members []AddPermissionMemberRequest `json:"members"`
}

type MemberResp struct {
	MemberType string `json:"member_type"`
	MemberId   string `json:"member_id"`
	Perm       string `json:"perm"`
	PermType   string `json:"perm_type"`
	Type       string `json:"type"`
}

type BatchAddPermissionMemberResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Members []MemberResp `json:"members"`
	} `json:"data"`
}

// AddPermissionMemberRequest represents the request body for adding permission member
type AddPermissionMemberRequest struct {
	MemberType string `json:"member_type"`
	MemberId   string `json:"member_id"`
	Perm       string `json:"perm"`
	PermType   string `json:"perm_type,omitempty"`
	Type       string `json:"type,omitempty"`
}

// AddPermissionMemberResponse represents the response for adding permission member
type AddPermissionMemberResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Member MemberResp `json:"member"`
	} `json:"data"`
}

// QueryRecordsRequest 查询记录请求
type QueryRecordsRequest struct {
	ViewId          string      `json:"view_id,omitempty"`
	FieldNames      []string    `json:"field_names,omitempty"`
	Filter          *FilterInfo `json:"filter,omitempty"`
	Sort            []Sort      `json:"sort,omitempty"`
	AutomaticFields bool        `json:"automatic_fields,omitempty"`
}

// FilterInfo 筛选条件
type FilterInfo struct {
	Conditions  []FilterCondition `json:"conditions,omitempty"`
	Conjunction string            `json:"conjunction,omitempty"` // "and" or "or"
}

// FilterCondition 筛选条件
type FilterCondition struct {
	FieldName string      `json:"field_name"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

// Sort 排序条件
type Sort struct {
	FieldName string `json:"field_name"`
	Desc      bool   `json:"desc,omitempty"`
}

// QueryRecordsResponse 查询记录响应
type QueryRecordsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total     int                      `json:"total"`
		HasMore   bool                     `json:"has_more"`
		PageToken string                   `json:"page_token"`
		Items     []map[string]interface{} `json:"items"`
	} `json:"data"`
}

// BatchDeleteRecordsRequest 批量删除记录请求
type BatchDeleteRecordsRequest struct {
	Records []string `json:"records"`
}

// BatchDeleteRecordsResponse 批量删除记录响应
type BatchDeleteRecordsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Records []struct {
			Deleted  bool   `json:"deleted"`
			RecordId string `json:"record_id"`
		} `json:"records"`
	} `json:"data"`
}

// BatchCreateRecords creates multiple records in a Bitable table
// appToken: the unique identifier of the Bitable app
// tableId: the unique identifier of the table
// records: slice of records to create, each record is a map of field name to value
func (c *Client) BatchCreateRecords(appToken, tableId string, records []map[string]interface{}) (*BatchCreateRecordsResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	// Prepare request body
	reqRecords := make([]Record, len(records))
	for i, r := range records {
		reqRecords[i] = Record{Fields: r}
	}

	reqBody := BatchCreateRecordsRequest{
		Records: reqRecords,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make API request
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_create", appToken, tableId)
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp BatchCreateRecordsResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

// CreateBitable creates a new Bitable (multi-dimensional table) in the specified folder
func (c *Client) CreateBitable(name, folderToken string) (*CreateBitableResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBody := CreateBitableRequest{
		Name:        name,
		FolderToken: folderToken,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make API request
	url := "https://open.feishu.cn/open-apis/bitable/v1/apps"
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp CreateBitableResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

// ListFields lists all fields in a Bitable table
func (c *Client) ListFields(appToken, tableID string, pageSize int, pageToken string) (*ListFieldsResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	// Build URL with query parameters
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/fields", appToken, tableID)
	if pageSize > 0 {
		url += fmt.Sprintf("?page_size=%d", pageSize)
		if pageToken != "" {
			url += fmt.Sprintf("&page_token=%s", pageToken)
		}
	} else if pageToken != "" {
		url += fmt.Sprintf("?page_token=%s", pageToken)
	}

	// Make API request
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodGet, url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp ListFieldsResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

// DeleteField deletes a field from a Bitable table
func (c *Client) DeleteField(appToken, tableID, fieldID string) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/fields/%s", appToken, tableID, fieldID)
	token, err := c.getTenantAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get tenant access token: %w", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request("DELETE", url, headers, nil)
	if err != nil {
		return fmt.Errorf("failed to make DELETE request: %w", err)
	}

	var resp DeleteFieldResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("API error: %s", resp.Msg)
	}
	return nil
}

// UpdateField updates a field in a Bitable table
func (c *Client) UpdateField(appToken, tableID, fieldID string, req UpdateFieldRequest) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/fields/%s", appToken, tableID, fieldID)
	token, err := c.getTenantAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get tenant access token: %w", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request("PUT", url, headers, reqBytes)
	if err != nil {
		return fmt.Errorf("failed to make PUT request: %w", err)
	}

	var resp UpdateFieldResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("API error: %s", resp.Msg)
	}
	return nil
}

// CreateField adds a new field to a Bitable table
func (c *Client) CreateField(appToken, tableID string, req CreateFieldRequest) error {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/fields", appToken, tableID)
	token, err := c.getTenantAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get tenant access token: %w", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request("POST", url, headers, reqBytes)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}

	var resp CreateFieldResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("API error: %s", resp.Msg)
	}
	return nil
}

func (c *Client) BatchGetId(req BatchGetIdRequest) (*BatchGetIdResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make API request
	url := "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id"
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	var resp BatchGetIdResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

func (c *Client) BatchAddPermissionMember(docToken, reqType string, req BatchAddPermissionMemberRequest) (*BatchAddPermissionMemberResponse, error) {
	accessToken, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make API request
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/drive/v1/permissions/%s/members/batch_create?need_notification=false&type=%s", docToken, reqType)
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	fmt.Println(respBody)
	var resp BatchAddPermissionMemberResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

func (c *Client) AddPermissionMember(docToken, reqType string, req AddPermissionMemberRequest) (*AddPermissionMemberResponse, error) {
	accessToken, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Make API request
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/drive/v1/permissions/%s/members?need_notification=false&type=%s", docToken, reqType)
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp AddPermissionMemberResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

// QueryRecords 查询记录
func (c *Client) QueryRecords(appToken, tableId string, req QueryRecordsRequest, pageToken string, pageSize int) (*QueryRecordsResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 构建URL和查询参数
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/search", appToken, tableId)
	if pageToken != "" {
		url = fmt.Sprintf("%s?page_token=%s", url, pageToken)
	}
	if pageSize > 0 {
		if strings.Contains(url, "?") {
			url = fmt.Sprintf("%s&page_size=%d", url, pageSize)
		} else {
			url = fmt.Sprintf("%s?page_size=%d", url, pageSize)
		}
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp QueryRecordsResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

// BatchDeleteRecords 批量删除记录
func (c *Client) BatchDeleteRecords(appToken, tableId string, req BatchDeleteRecordsRequest) (*BatchDeleteRecordsResponse, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_delete", appToken, tableId)
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp BatchDeleteRecordsResponse
	if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("API error: %s", resp.Msg)
	}

	return &resp, nil
}

type UpdatePermissionReq struct {
	ExternalAccessEntity     string `json:"external_access_entity,omitempty"`     // 允许内容被分享到组织外
	SecurityEntity           string `json:"security_entity,omitempty"`            // 谁可以创建副本、打印、下载
	CommentEntity            string `json:"comment_entity,omitempty"`             // 谁可以评论
	ShareEntity              string `json:"share_entity,omitempty"`               // 谁可以添加和管理协作者-组织维度
	ManageCollaboratorEntity string `json:"manage_collaborator_entity,omitempty"` // 谁可以添加和管理协作者-协作者维度
	LinkShareEntity          string `json:"link_share_entity,omitempty"`          // 链接分享设置
	CopyEntity               string `json:"copy_entity,omitempty"`                // 谁可以复制内容
}

type UpdatePermissionResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type OpenPasswordResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Password string `json:"password"`
	} `json:"data"`
}

func (c *Client) OpenPassword(appToken, fileType string) (*OpenPasswordResp, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/drive/v1/permissions/%s/public/password?type=%s", appToken, fileType)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	respBody, err := httpc.Request(http.MethodPost, url, headers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp OpenPasswordResp
	err = json.Unmarshal([]byte(respBody), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &resp, nil
}

func (c *Client) UpdatePermission(appToken, fileType string, req *UpdatePermissionReq) (*UpdatePermissionResp, error) {
	token, err := c.getTenantAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant access token: %v", err)
	}

	url := fmt.Sprintf("https://open.feishu.cn/open-apis/drive/v2/permissions/%s/public?type=%s", appToken, fileType)

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json; charset=utf-8",
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	respBody, err := httpc.Request(http.MethodPatch, url, headers, reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}

	var resp UpdatePermissionResp
	err = json.Unmarshal([]byte(respBody), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &resp, nil
}
