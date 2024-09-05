package ns

import (
	"embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	common "github.com/sat20-labs/name-dns/common"
	serverCommon "github.com/sat20-labs/name-dns/server/define"
	"go.etcd.io/bbolt"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed assets/*
var assetsFiles embed.FS

var siteMapIndex = &SiteMapIndex{
	XMLNS:           "http://www.sitemaps.org/schemas/sitemap/0.9",
	SiteMapItemList: []*SiteMapIndexItem{},
}

const (
	SITE_MAP_ITEM_COUNT = 500
	GEN_SITE_MAP_TIME   = 10 * time.Minute
)

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
	err = common.InitBucket(s.DB, BUCKET_COMMON_SUMMARY)
	if err != nil {
		return
	}

	s.initRouter(r)
	return
}

func (s *Service) initRouter(r *gin.Engine) {
	staticServer := http.FS(staticFiles)
	staticFileServer := http.FileServer(staticServer)
	assetsServer := http.FS(assetsFiles)
	assetsFileServer := http.FileServer(assetsServer)
	r.GET("/static/*filepath", func(c *gin.Context) {
		staticFileServer.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/assets/*filepath", func(c *gin.Context) {
		assetsFileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/robots.txt", s.robots)
	r.GET("/sitemap/sitemap_index.xml", s.siteMapIndex)
	r.GET("/sitemap/:index.xml", s.siteMapItem)
	r.GET("/", s.content)

	r.GET("/summary/name-count", s.countHtml)
}
