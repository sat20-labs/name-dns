package ns

import (
	"github.com/gin-gonic/gin"
	common "github.com/sat20-labs/name-ns/common"
	serverCommon "github.com/sat20-labs/name-ns/server/define"
	"go.etcd.io/bbolt"
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
	err = common.InitBucket(s.DB, BUCKET_NAME)
	if err != nil {
		return
	}
	s.initRouter(r)
	return
}

func (s *Service) initRouter(r *gin.Engine) {
	r.GET("/", s.getContent)
	r.GET("/namecount/:name", s.getNameCount)
	r.GET("/proxyreq", s.proxyReq)
}
