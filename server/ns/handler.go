package ns

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sat20-labs/name-ns/common"
)

func (s *Service) getContent(c *gin.Context) {
	if !strings.Contains(c.Request.Host, s.rpcConfig.Domain) {
		c.String(http.StatusNotFound, "invalid host")
		return
	}
	name := strings.TrimSuffix(c.Request.Host, "."+s.rpcConfig.Domain)
	if name == "" {
		c.String(http.StatusNotFound, "invalid host")
		return
	}
	nsRoutingResp, _, err := common.RpcRequest(s.ordxRpcConfig.NsRouting, name, "GET")
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var nameRoutingResp NameRoutingResp
	err = json.Unmarshal(nsRoutingResp, &nameRoutingResp)
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if nameRoutingResp.Code != 0 {
		c.String(http.StatusNotFound, nameRoutingResp.Msg)
		return
	}

	inscriptionContent, header, err := common.RpcRequest(s.ordxRpcConfig.InscriptionContent, nameRoutingResp.Data.InscriptionId, "GET")
	if err != nil {
		common.Log.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for key, values := range header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	contentType := header.Get("Content-Type")
	c.Data(http.StatusOK, contentType, inscriptionContent)
}
