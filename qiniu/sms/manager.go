/*
 * Copyright (c) 2018. Brickman.
 */

package sms

import (
	"fmt"
	"github.com/qiniu/api.v7/v7/sms"
	"net/http"
	"net/url"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/sms/client"
	"github.com/qiniu/api.v7/v7/sms/rpc"
)

var (
	// Host 为 QiNiu SMS Server API 服务域名
	Host = "https://sms.qiniuapi.com"
)

// Manager 提供了 Qiniu SMS Server API 相关功能
type Manager struct {
	mac    *auth.Credentials
	client rpc.Client
}

// NewManager 用来构建一个新的 Manager
func NewManager(mac *auth.Credentials) (manager *Manager) {

	manager = &Manager{}

	mac1 := &client.Mac{
		AccessKey: mac.AccessKey,
		SecretKey: mac.SecretKey,
	}

	transport := client.NewTransport(mac1, nil)
	manager.client = rpc.Client{Client: &http.Client{Transport: transport}}

	return
}

// QuerySignature 查询签名
func (m *Manager) QuerySignature(args QuerySignatureRequest) (pagination sms.SignaturePagination, err error) {
	values := url.Values{}

	if args.AuditStatus.IsValid() {
		values.Set("audit_status", args.AuditStatus.String())
	}

	if args.Page > 0 {
		values.Set("page", fmt.Sprintf("%d", args.Page))
	}

	if args.PageSize > 0 {
		values.Set("page_size", fmt.Sprintf("%d", args.PageSize))
	}

	reqURL := fmt.Sprintf("%s%s?%s", Host, "/v1/signature", values.Encode())
	err = m.client.GetCall(&pagination, reqURL)
	return
}

// QueryTemplate 查询模板
func (m *Manager) QueryTemplate(args sms.QueryTemplateRequest) (pagination TemplatePagination, err error) {
	values := url.Values{}

	if args.AuditStatus.IsValid() {
		values.Set("audit_status", args.AuditStatus.String())
	}

	if args.Page > 0 {
		values.Set("page", fmt.Sprintf("%d", args.Page))
	}

	if args.PageSize > 0 {
		values.Set("page_size", fmt.Sprintf("%d", args.PageSize))
	}

	reqURL := fmt.Sprintf("%s%s?%s", Host, "/v1/template", values.Encode())
	err = m.client.GetCall(&pagination, reqURL)
	return
}

// QueryTemplate 查询模板
func (m *Manager) QueryTemplateByID(id string) (smsTemplate Template, err error) {
	reqURL := fmt.Sprintf("%s%s/%s", Host, "/v1/template", id)
	err = m.client.GetCall(&smsTemplate, reqURL)
	return
}

// CreateTemplate 创建模板
func (m *Manager) CreateTemplate(args TemplateRequest) (ret sms.TemplateResponse, err error) {
	reqURL := fmt.Sprintf("%s%s", Host, "/v1/template")
	err = m.client.CallWithJSON(&ret, reqURL, args)
	return
}

// SendMessage 发送短信
func (m *Manager) SendMessage(args MessagesRequest) (ret MessagesResponse, err error) {
	reqURL := fmt.Sprintf("%s%s", Host, "/v1/message")
	err = m.client.CallWithJSON(&ret, reqURL, args)
	return
}

// QueryTemplateRequest 查询短信参数
type QueryTemplateRequest struct {
	JobID     string `json:"job_id"`
	MessageID string `json:"message_id"`
	Mobile    string `json:"mobile"`
}

// Message 模板
type Message struct {
	Content   string `json:"content" xml:"content"`
	Type      string `json:"type" xml:"type"`
	Error     string `json:"error" xml:"error"`
	Count     int    `json:"count" xml:"count"`
	JobID     string `json:"job_id" xml:"job_id"`
	Mobile    string `json:"mobile" xml:"mobile"`
	CreatedAt uint64 `json:"created_at" xml:"created_at"`
	DelivrdAt uint64 `json:"delivrd_at" xml:"delivrd_at"`
	MessageID string `json:"message_id" xml:"message_id"`
	Status    string `json:"status" xml:"status"`
}

type QueryMessageRequest struct {
	JobID     string `json:"job_id"`     // job_id
	MessageID string `json:"message_id"` // message_id
	Mobile    string `json:"mobile"`     // mobile
	Page      int    `json:"page"`       // 页码，默认为 1
	PageSize  int    `json:"page_size"`  // 分页大小，默认为 20
}

// MessagePagination 短信分页
type MessagePagination struct {
	Page     int       `json:"page"`      // 页码，默认为 1
	PageSize int       `json:"page_size"` // 分页大小，默认为 20
	Total    int       `json:"total"`     // 总记录条数
	Items    []Message `json:"items"`     // 模板
}

// QueryMessage 查询短信发送状态
func (m *Manager) QueryMessage(args QueryMessageRequest) (pagination MessagePagination, err error) {
	values := url.Values{}

	if args.JobID != "" {
		values.Set("job_id", args.JobID)
	}

	if args.MessageID != "" {
		values.Set("message_id", args.MessageID)
	}

	if args.Mobile != "" {
		values.Set("mobile", args.Mobile)
	}

	if args.Page > 0 {
		values.Set("page", fmt.Sprintf("%d", args.Page))
	}

	if args.PageSize > 0 {
		values.Set("page_size", fmt.Sprintf("%d", args.PageSize))
	}

	reqURL := fmt.Sprintf("%s%s?%s", Host, "/v1/messages", values.Encode())
	err = m.client.GetCall(&pagination, reqURL)
	return
}

type TemplateRequest struct {
	UID         uint32           `json:"uid"`
	Name        string           `json:"name"`
	Type        sms.TemplateType `json:"type"`
	Template    string           `json:"template"`
	Description string           `json:"description"`
	SignatureID string           `json:"signature_id"`
}

// Template 模板
type Template struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Type          sms.TemplateType `json:"type"`
	Template      string           `json:"template"`
	Description   string           `json:"description"`
	AuditStatus   sms.AuditStatus  `json:"audit_status"`
	RejectReason  string           `json:"reject_reason"`
	SignatureID   string           `json:"signature_id"`   // 模版绑定的签名ID
	SignatureText string           `json:"signature_text"` // 模版绑定的签名内容

	UpdatedAt uint64 `json:"updated_at"`
	CreatedAt uint64 `json:"created_at"`
}

// TemplatePagination 模板分页
type TemplatePagination struct {
	Page     int        `json:"page"`      // 页码，默认为 1
	PageSize int        `json:"page_size"` // 分页大小，默认为 20
	Total    int        `json:"total"`     // 总记录条数
	Items    []Template `json:"items"`     // 模板
}

// QuerySignatureRequest 查询签名参数
type QuerySignatureRequest struct {
	AuditStatus sms.AuditStatus `json:"audit_status"` // 审核状态
	Page        int             `json:"page"`         // 页码，默认为 1
	PageSize    int             `json:"page_size"`    // 分页大小，默认为 20
}

// MessagesRequest 短信消息
type MessagesRequest struct {
	TemplateID string                 `json:"template_id"`
	Mobiles    []string               `json:"mobiles"`
	Parameters map[string]interface{} `json:"parameters"`
}

// MessagesResponse 发送短信响应
type MessagesResponse struct {
	JobID string `json:"job_id"`
}
