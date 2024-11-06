package server

import (
	"BaselineCheck/server/compliance"
	"BaselineCheck/server/config"
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

type BaselineServer struct {
	gin     *gin.Engine
	conf    *config.Config
	handler *compliance.Handler
}

func NewServer(conf *config.Config, handler *compliance.Handler) *BaselineServer {
	return &BaselineServer{
		conf:    conf,
		handler: handler,
	}
}

//go:embed templates/*
var templates embed.FS

func (b *BaselineServer) Start() {
	b.gin = gin.Default()

	// 使用 embed 中的模板文件
	b.gin.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))
	b.gin.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})
	b.gin.POST("/check", b.handler.RegisterComplianceResult)
	b.gin.GET("/hostlist", b.handler.ShowComplianceHostList)
	b.gin.GET("/tasks", b.handler.ShowComplianceResultList)
	b.gin.GET("/check/:result_id/:id", b.handler.ShowComplianceDetails)
	addr := fmt.Sprintf("%s:%d", b.conf.Host, b.conf.Port)
	b.gin.Run(addr)
}
