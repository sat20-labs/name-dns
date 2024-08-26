package ns

import (
	"github.com/gin-gonic/gin"
	serverCommon "github.com/sat20-labs/name-ns/server/define"
)

type Service struct {
	rpcConfig     *serverCommon.Rpc
	ordxRpcConfig *serverCommon.OrdxRpc
}

func New(rpcConfig *serverCommon.Rpc, ordxRpcConfig *serverCommon.OrdxRpc) *Service {
	return &Service{
		rpcConfig:     rpcConfig,
		ordxRpcConfig: ordxRpcConfig,
	}
}

func (s *Service) InitRouter(r *gin.Engine) {
	r.GET("/", s.getContent)
}
