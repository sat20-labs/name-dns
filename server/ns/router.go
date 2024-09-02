package ns

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/sat20-labs/name-dns/common"
	serverCommon "github.com/sat20-labs/name-dns/server/define"
	"go.etcd.io/bbolt"
)

//go:embed static/*
var staticFiles embed.FS

type Service struct {
	RpcConfig     *serverCommon.Rpc
	OrdxRpcConfig *serverCommon.OrdxRpc
	DB            *bbolt.DB
}

func New(
	rpcConfig *serverCommon.Rpc,
	ordxRpcConfig *serverCommon.OrdxRpc,
	db *bbolt.DB) *Service {

	return &Service{
		RpcConfig:     rpcConfig,
		OrdxRpcConfig: ordxRpcConfig,
		DB:            db,
	}
}

func (s *Service) Init(r *gin.Engine) (err error) {
	err = common.InitBucket(s.DB, BUCKET_NAME_COUNT)
	if err != nil {
		return
	}
	s.initRouter(r)
	return
}

func (s *Service) initRouter(r *gin.Engine) {
	staticServer := http.FS(staticFiles)
	r.StaticFS("/static", staticServer)
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.FileFromFS("static/favicon.ico", staticServer)
	})
	r.GET("/sitemap.xml", s.siteMap)
	r.GET("/robots.txt", func(c *gin.Context) {
		c.FileFromFS("static/robots.txt", staticServer)
	})

	r.GET("/", s.content)
	r.GET("/summary/name-count", s.countHtml)
}
