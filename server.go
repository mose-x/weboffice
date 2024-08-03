package weboffice

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// UserContext 用户上下文
type UserContext struct {
	context.Context

	appID     string
	token     string
	query     url.Values
	requestID string
}

// Server 服务配置
type Server struct {
	config Config
	engine *gin.Engine
	root   gin.IRouter
}

func (uc *UserContext) AppID() string {
	return uc.appID
}
func (uc *UserContext) Token() string {
	return uc.token
}
func (uc *UserContext) Query() url.Values {
	return uc.query
}
func (uc *UserContext) RequestID() string {
	return uc.requestID
}

// ParseContext 解析上下文
func ParseContext(req *http.Request) Context {
	uc := &UserContext{
		Context:   req.Context(),
		appID:     req.Header.Get("X-App-ID"),
		token:     req.Header.Get("X-WebOffice-Token"),
		requestID: req.Header.Get("X-Request-ID"),
	}
	if uc.Token() == "" {
		log.Panic("token is empty")
	}
	if v, err := url.ParseQuery(req.Header.Get("X-User-Query")); err == nil {
		uc.query = v
	} else {
		uc.query = url.Values{}
	}
	return uc
}

// wrapHandlerFunc 包装handler
func (srv *Server) wrapHandlerFunc(f func(*gin.Context) (any, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		begin := time.Now()
		data, err := f(c)
		cost := time.Since(begin)

		if err != nil {
			var respErr *Error
			var e *Error
			if errors.As(err, &e) {
				respErr = e
			}

			c.JSON(respErr.StatusCode(), &Reply{Code: respErr.Code(), Message: respErr.Message()})
			srv.config.Logger.Error("%s %s code=%d message=%s cost=%s", c.Request.Method, c.Request.RequestURI, respErr.Code(), cost.String())
		} else {
			c.JSON(http.StatusOK, &Reply{Code: OK, Data: data})
			srv.config.Logger.Info("%s %s code=OK cost=%s", c.Request.Method, c.Request.RequestURI, cost.String())
		}
	}
}

// registerRoutes 注册路由
func (srv *Server) registerRoutes(router gin.IRouter) {
	router.GET("/v3/3rd/files/:file_id", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
		fileID := c.Param("file_id")
		ctx := ParseContext(c.Request)

		return srv.config.GetFile(ctx, fileID)
	}))
	router.GET("/v3/3rd/files/:file_id/download", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
		fileID := c.Param("file_id")
		ctx := ParseContext(c.Request)

		return srv.config.GetFileDownload(ctx, fileID)
	}))
	router.GET("/v3/3rd/files/:file_id/permission", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
		fileID := c.Param("file_id")
		ctx := ParseContext(c.Request)

		return srv.config.GetFilePermission(ctx, fileID)
	}))

	if srv.config.UserProvider != nil {
		router.GET("/v3/3rd/users", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			userIDs := c.QueryArray("user_ids")
			ctx := ParseContext(c.Request)

			return srv.config.GetUsers(ctx, userIDs)
		}))
	}
	if srv.config.WatermarkProvider != nil {
		router.GET("/v3/3rd/files/:file_id/watermark", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			fileID := c.Param("file_id")
			ctx := ParseContext(c.Request)

			return srv.config.GetFileWatermark(ctx, fileID)
		}))
	}

	if srv.config.EditProvider != nil {
		router.POST("/v3/3rd/files/:file_id/upload", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			fileID := c.Param("file_id")
			ctx := ParseContext(c.Request)

			fileHeader, err := c.FormFile("file")
			if err != nil {
				return nil, ErrInvalidArguments.WithMessage(err.Error())
			}
			f, err := fileHeader.Open()
			if err != nil {
				return nil, ErrInternalError.WithMessage(err.Error())
			}
			defer func(f multipart.File) {
				err := f.Close()
				if err != nil {
					log.Println(err.Error())
				}
			}(f)

			var args UpdateFilePhaseArgs
			args.Name = c.PostForm("name")
			args.SHA1 = c.PostForm("sha1")
			args.Size, _ = strconv.ParseInt(c.PostForm("size"), 10, 64)
			args.IsManual, _ = strconv.ParseBool(c.PostForm("is_manual"))
			args.Content = f

			return srv.config.UpdateFile(ctx, fileID, &args)
		}))

		router.PUT("/v3/3rd/files/:file_id/name", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			ctx := ParseContext(c.Request)
			fileID := c.Param("file_id")

			var args RenameFileArgs
			if err := c.BindJSON(&args); err != nil {
				return nil, ErrInvalidArguments.WithMessage(err.Error())
			}
			if err := srv.config.RenameFile(ctx, fileID, &args); err != nil {
				return nil, err
			} else {
				return &Empty{}, nil
			}
		}))
	}

	// 版本管理 api
	if srv.config.VersionProvider != nil {
		router.GET("/v3/3rd/files/:file_id/versions", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			ctx := ParseContext(c.Request)
			fileID := c.Param("file_id")
			offset, _ := strconv.Atoi(c.Query("offset"))
			limit, _ := strconv.Atoi(c.Query("limit"))

			return srv.config.VersionProvider.GetFileVersions(ctx, fileID, offset, limit)
		}))
		router.GET("/v3/3rd/files/:file_id/versions/:version", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			ctx := ParseContext(c.Request)
			fileID := c.Param("file_id")
			versionID, _ := strconv.Atoi(c.Param("version"))

			return srv.config.VersionProvider.GetFileVersion(ctx, fileID, int32(versionID))
		}))
		router.GET("/v3/3rd/files/:file_id/versions/:version/download", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			ctx := ParseContext(c.Request)
			fileID := c.Param("file_id")
			versionID, _ := strconv.Atoi(c.Param("version"))

			return srv.config.VersionProvider.GetFileVersionDownload(ctx, fileID, int32(versionID))
		}))
	}

	// notify 通知api
	if srv.config.NotifyProvider != nil {
		router.POST("/v3/3rd/notify", srv.wrapHandlerFunc(func(c *gin.Context) (any, error) {
			var args NotifyArgs
			if err := c.BindJSON(&args); err != nil {
				return nil, ErrInvalidArguments.WithMessage(err.Error())
			}
			err := srv.config.OnNotify(ParseContext(c.Request), &args)
			if err != nil {
				return nil, err
			} else {
				return &Empty{}, nil
			}
		}))
	}
}

// NewServer 创建服务
func NewServer(config Config, e *gin.Engine) {
	if config.PreviewProvider == nil {
		log.Panic("PreviewProvider must not nil")
	}
	if config.Logger == nil {
		config.Logger = &noopLogger{}
	}
	srv := &Server{
		engine: e,
		config: config,
	}
	if config.Prefix == "" {
		srv.root = srv.engine
	} else {
		srv.root = srv.engine.Group(config.Prefix)
	}

	// 注册路由到主engine
	srv.registerRoutes(srv.root)
}

// Run 启动服务
func (srv *Server) Run(addr string) error {
	return srv.engine.Run(addr)
}

// Router 获取路由
func (srv *Server) Router() gin.IRouter {
	return srv.root
}

// Handler 获取http handler
func (srv *Server) Handler() http.Handler {
	return srv.engine
}
