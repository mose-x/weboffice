package weboffice

import (
	"io"
)

var Referer = "https://solution.wps.cn"

// Reply 返回参数结构体
type Reply struct {
	Code    Code   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

// Empty 空结构体
type Empty struct {
}

// GetFileReply 获取文件返回信息
type GetFileReply struct {
	CreateTime int64  `json:"create_time"`
	CreatorId  string `json:"creator_id"`
	ID         string `json:"id"`
	ModifierId string `json:"modifier_id"`
	ModifyTime int64  `json:"modify_time"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Version    int32  `json:"version"`
}

// GetFileDownloadReply 获取文件下载地址返回信息
type GetFileDownloadReply struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// GetFilePermissionReply 获取文件权限返回信息
type GetFilePermissionReply struct {
	Comment  int    `json:"comment"`
	Copy     int    `json:"copy"`
	Download int    `json:"download"`
	History  int    `json:"history"`
	Print    int    `json:"print"`
	Read     int    `json:"read"`
	Rename   int    `json:"rename"`
	SaveAs   int    `json:"saveas"`
	Update   int    `json:"update"`
	UserId   string `json:"user_id"`
}

// UserReply 用户信息
type UserReply struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// GetWatermarkReply 获取水印返回信息
type GetWatermarkReply struct {
	Type       int     `json:"type"`
	Value      string  `json:"value"`
	FillStyle  string  `json:"fill_style"`
	Font       string  `json:"font"`
	Rotate     float64 `json:"rotate"`
	Horizontal int     `json:"horizontal"`
	Vertical   int     `json:"vertical"`
}

// UpdateFilePhaseArgs 更新文件阶段参数
type UpdateFilePhaseArgs struct {
	Name     string
	Size     int64
	SHA1     string
	IsManual bool
	Content  io.Reader
}

// RenameFileArgs 重命名文件参数
type RenameFileArgs struct {
	Name string `json:"name"`
}

// NotifyArgs 通知参数
type NotifyArgs struct {
	FileId  string        `json:"file_id,omitempty"`
	Type    string        `json:"type,omitempty"`
	Content NotifyContent `json:"content,omitempty"`
}

// NotifyContent 通知内容
type NotifyContent struct {
	SessionId       string `json:"session_id,omitempty"`
	InitVersion     int    `json:"init_version,omitempty"`
	Readonly        bool   `json:"readonly,omitempty"`
	UploadedVersion int    `json:"uploaded_version,omitempty"`
	LastModifierId  string `json:"last_modifier_id,omitempty"`
	ConnectionId    string `json:"connection_id,omitempty"`
	UserId          string `json:"user_id,omitempty"`
	Permission      string `json:"permission,omitempty"`
	Print           bool   `json:"print,omitempty"`
	Format          string `json:"format,omitempty"`
}

// Config 配置
type Config struct {
	// 注册预览服务
	PreviewProvider

	// 注册用户服务
	UserProvider

	// 注册水印服务
	WatermarkProvider

	// 注册编辑服务
	EditProvider

	// 注册版本服务
	VersionProvider

	// 路由前缀（上传用）
	Prefix string

	// 注册日志服务
	Logger

	// 注册通知服务
	NotifyProvider
}
