# [WPS WebOffice 开放平台](https://solution.wps.cn) Go SDK V3

## 依赖

- go 1.16+
- gin framework

## How to use
```shell
go get github.com/mose-x/weboffice
```

```go
package main

func InitRouter() {
	// 初始化路由
	e := gin.Default()

	// 注册wps web office服务
	provider := &Provider{}
	weboffice.NewServer(weboffice.Config{
		PreviewProvider:   provider,
		UserProvider:      provider,
		WatermarkProvider: provider,
		EditProvider:      provider,
		VersionProvider:   provider,
		Logger:            weboffice.DefaultLogger(),
		NotifyProvider:    provider,
	}, e)
	
	// 启动服务
	_ = e.Run(":8080")
}
```
### 实现接口
```go
package main

type Provider struct {
}

func (*Provider) GetFileWatermark(_ weboffice.Context, _ string) (*weboffice.GetWatermarkReply, error) {
	return &weboffice.GetWatermarkReply{
		Type:       1,
		Value:      "mose",
		FillStyle:  "rgba(192,192,192,0.6)",
		Font:       "bold 20px Serif",
		Rotate:     0.5,
		Horizontal: 50,
		Vertical:   50,
	}, nil
}
```

### 实际效果

--- --
[docx 在线预览/编辑](https://qnfile.ljserver.cn/weboffice/docx.html)
-- -------------------------------------------
[pptx 在线预览/编辑](https://qnfile.ljserver.cn/weboffice/pptx.html)
-- -------------------------------------------
[xlsx 在线预览/编辑](https://qnfile.ljserver.cn/weboffice/xlsx.html)
-- -------------------------------------------
[pdf 在线预览/编辑](https://qnfile.ljserver.cn/weboffice/pdf.html)
--- --


## 更多

关于接口的更多说明，请参考[WebOffice开放平台-WebOffice回调配置](https://solution.wps.cn/docs/callback/summary.html)。