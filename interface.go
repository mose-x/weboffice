package weboffice

import (
	"context"
	"net/url"
)

// Context 定义上下文接口
type Context interface {
	context.Context

	// AppID 获取请求ID
	AppID() string

	// Token 获取请求Token
	Token() string

	// Query 获取请求参数
	Query() url.Values

	// RequestID 获取请求ID
	RequestID() string
}

// PreviewProvider 定义预览接口
type PreviewProvider interface {
	// GetFile 获取文件信息
	GetFile(ctx Context, fileID string) (*GetFileReply, error)

	// GetFileDownload 获取文件下载地址
	GetFileDownload(ctx Context, fileID string) (*GetFileDownloadReply, error)

	// GetFilePermission 获取文件权限信息
	GetFilePermission(ctx Context, fileID string) (*GetFilePermissionReply, error)
}

// UserProvider 定义用户接口
type UserProvider interface {
	// GetUsers 获取用户信息
	GetUsers(ctx Context, userIDs []string) ([]*UserReply, error)
}

// WatermarkProvider 定义水印接口
type WatermarkProvider interface {
	// GetFileWatermark 获取文件水印信息
	GetFileWatermark(ctx Context, fileID string) (*GetWatermarkReply, error)
}

// EditProvider 定义编辑接口
type EditProvider interface {
	// UpdateFile 更新文件
	UpdateFile(ctx Context, fileID string, args *UpdateFilePhaseArgs) (*GetFileReply, error)

	// RenameFile 重命名文件
	RenameFile(ctx Context, fileID string, args *RenameFileArgs) error
}

// VersionProvider 定义版本接口
type VersionProvider interface {
	// GetFileVersions 获取文件版本列表
	GetFileVersions(ctx Context, fileID string, offset, limit int) ([]*GetFileReply, error)

	// GetFileVersion 获取文件指定版本信息
	GetFileVersion(ctx Context, fileID string, version int32) (*GetFileReply, error)

	// GetFileVersionDownload 获取文件指定版本下载地址
	GetFileVersionDownload(ctx Context, fileID string, version int32) (*GetFileDownloadReply, error)
}

// NotifyProvider 定义通知接口
type NotifyProvider interface {
	// OnNotify 触发通知
	OnNotify(ctx Context, args *NotifyArgs) error
}
